package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"

	"tea-system/internal/invoice"
	"tea-system/internal/models"
	"tea-system/internal/repository"
)

// InvoiceService — 发票生成服务
type InvoiceService struct {
	orderRepo     *repository.OrderRepo
	invoiceRepo   *repository.InvoiceRepo
	userRepo      *repository.UserRepo
	storageDir    string
	rabbitMQURL   string
	rmqConn       *amqp091.Connection
	asyncEnabled  bool
}

func NewInvoiceService(
	orderRepo *repository.OrderRepo,
	invoiceRepo *repository.InvoiceRepo,
	userRepo *repository.UserRepo,
	storageDir string,
	rabbitMQURL string,
) *InvoiceService {
	if storageDir == "" {
		storageDir = "./storage/invoices"
	}
	// 确保目录存在
	_ = os.MkdirAll(storageDir, 0o755)

	svc := &InvoiceService{
		orderRepo:   orderRepo,
		invoiceRepo: invoiceRepo,
		userRepo:    userRepo,
		storageDir:  storageDir,
		rabbitMQURL: rabbitMQURL,
	}
	svc.asyncEnabled = svc.tryConnectRabbitMQ()
	return svc
}

// tryConnectRabbitMQ — 尝试连接 RabbitMQ；失败则降级为同步模式
func (s *InvoiceService) tryConnectRabbitMQ() bool {
	if s.rabbitMQURL == "" {
		return false
	}
	conn, err := amqp091.Dial(s.rabbitMQURL)
	if err != nil {
		log.Warn().Err(err).Msg("invoice_service: rabbitmq unavailable, using sync mode")
		return false
	}
	s.rmqConn = conn

	ch, err := conn.Channel()
	if err != nil {
		log.Warn().Err(err).Msg("invoice_service: rabbitmq channel failed, using sync mode")
		_ = conn.Close()
		return false
	}
	defer ch.Close()

	_, err = ch.QueueDeclare("invoice_generation", true, false, false, false, nil)
	if err != nil {
		log.Warn().Err(err).Msg("invoice_service: rabbitmq queue declare failed, using sync mode")
		_ = conn.Close()
		return false
	}

	log.Info().Msg("invoice_service: rabbitmq connected, async mode enabled")
	return true
}

// GenerateInvoice — 同步生成发票 PDF + 存 DB + 返回 Invoice
// 如果 rabbitmq 可用也会 publish 一个 async 任务（但同步先做）
func (s *InvoiceService) GenerateInvoice(ctx context.Context, orderID uint64) (*models.Invoice, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	var user *models.User
	if order.UserID != nil {
		user, err = s.userRepo.GetByID(ctx, *order.UserID)
	}
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 组装 Input
	invoiceNo := generateInvoiceNo()
	issueDate := time.Now()

	customerName, customerEmail, customerAddr := "", "", ""
	if user != nil {
		customerName = user.Name
		customerEmail = user.Email
		customerAddr = formatAddress(user.BillingAddress)
	}
	if customerName == "" && order.BillingAddressSnapshot != nil {
		if v, ok := order.BillingAddressSnapshot["name"].(string); ok { customerName = v }
		if v, ok := order.BillingAddressSnapshot["email"].(string); ok { customerEmail = v }
	}

	// 从 Order.CustomProductSnapshot 拿产品信息
	productName := ""
	if order.CustomProductSnapshot != nil {
		if v, ok := order.CustomProductSnapshot["title"].(string); ok {
			productName = v
		}
	}
	if productName == "" {
		productName = fmt.Sprintf("Tea Product (Order %s)", order.OrderNo)
	}

	inBytes, err := invoice.Generate(invoice.Input{
		InvoiceNo:       invoiceNo,
		IssueDate:       issueDate,
		OrderNo:         order.OrderNo,
		Currency:        "GBP",
		HsCode:          order.HsCode,
		CountryOfOrigin: order.CountryOfOrigin,
		ShippingCost:    order.ShippingCost,
		CustomerName:    customerName,
		CustomerEmail:   customerEmail,
		CustomerAddress: customerAddr,
		Items: []invoice.InvoiceItem{
			{
				Name:      productName,
				Quantity:  order.Quantity,
				UnitPrice: order.UnitPrice,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("pdf generate failed: %w", err)
	}

	// 存 PDF 到本地文件系统
	fileName := fmt.Sprintf("%s-%s.pdf", order.OrderNo, invoiceNo)
	filePath := filepath.Join(s.storageDir, fileName)
	if err := os.WriteFile(filePath, inBytes, 0o644); err != nil {
		return nil, fmt.Errorf("write pdf failed: %w", err)
	}

	// 创建 Invoice 记录
	inv := &models.Invoice{
		OrderID:         order.ID,
		InvoiceNo:       invoiceNo,
		IssueDate:       issueDate,
		PdfURL:          filePath,
		HsCode:          order.HsCode,
		CountryOfOrigin: order.CountryOfOrigin,
		TotalAmount:     order.TotalAmount,
		Locked:          false,
	}

	if err := s.invoiceRepo.CreateInvoice(ctx, inv); err != nil {
		return nil, fmt.Errorf("invoice db create failed: %w", err)
	}

	// 如果启用了 RabbitMQ，发布一个异步事件（用于归档/邮件通知等下游消费）
	if s.asyncEnabled {
		s.publishGeneratedEvent(ctx, inv)
	}

	return inv, nil
}

// RegenerateInvoice — 重新生成（更新 PDF + 更新 DB）
func (s *InvoiceService) RegenerateInvoice(ctx context.Context, invID uint64) (*models.Invoice, error) {
	inv, err := s.invoiceRepo.GetInvoiceByID(ctx, invID)
	if err != nil {
		return nil, err
	}

	order, err := s.orderRepo.GetByID(ctx, inv.OrderID)
	if err != nil {
		return nil, err
	}
	var user *models.User
	if order.UserID != nil {
		user, err = s.userRepo.GetByID(ctx, *order.UserID)
	}
	if err != nil {
		return nil, err
	}

	customerName, customerEmail, customerAddr := "", "", ""
	if user != nil {
		customerName = user.Name
		customerEmail = user.Email
		customerAddr = formatAddress(user.BillingAddress)
	}
	if customerName == "" && order.BillingAddressSnapshot != nil {
		if v, ok := order.BillingAddressSnapshot["name"].(string); ok { customerName = v }
	}

	productName := ""
	if order.CustomProductSnapshot != nil {
		if v, ok := order.CustomProductSnapshot["title"].(string); ok {
			productName = v
		}
	}
	if productName == "" {
		productName = fmt.Sprintf("Tea Product (Order %s)", order.OrderNo)
	}

	inBytes, err := invoice.Generate(invoice.Input{
		InvoiceNo:       inv.InvoiceNo,
		IssueDate:       time.Now(),
		OrderNo:         order.OrderNo,
		Currency:        "GBP",
		HsCode:          inv.HsCode,
		CountryOfOrigin: inv.CountryOfOrigin,
		ShippingCost:    order.ShippingCost,
		CustomerName:    customerName,
		CustomerEmail:   customerEmail,
		CustomerAddress: customerAddr,
		Items: []invoice.InvoiceItem{
			{Name: productName, Quantity: order.Quantity, UnitPrice: order.UnitPrice},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("pdf regenerate failed: %w", err)
	}

	// 更新文件
	fileName := fmt.Sprintf("%s-%s.pdf", order.OrderNo, inv.InvoiceNo)
	filePath := filepath.Join(s.storageDir, fileName)
	if err := os.WriteFile(filePath, inBytes, 0o644); err != nil {
		return nil, fmt.Errorf("write pdf failed: %w", err)
	}

	if err := s.invoiceRepo.UpdateInvoicePDFURL(ctx, invID, filePath); err != nil {
		return nil, err
	}
	inv.PdfURL = filePath
	return inv, nil
}

// ReadPDF — 读取指定 invoice 的 PDF 字节（供 handler 下载）
func (s *InvoiceService) ReadPDF(pdfPath string) ([]byte, error) {
	return os.ReadFile(pdfPath)
}

// publishGeneratedEvent — 异步事件
func (s *InvoiceService) publishGeneratedEvent(ctx context.Context, inv *models.Invoice) {
	if s.rmqConn == nil {
		return
	}
	ch, err := s.rmqConn.Channel()
	if err != nil {
		log.Warn().Err(err).Msg("invoice_service: rmq channel failed on publish")
		return
	}
	defer ch.Close()

	body, _ := json.Marshal(map[string]interface{}{
		"event":      "invoice.generated",
		"invoice_id": inv.ID,
		"invoice_no": inv.InvoiceNo,
		"order_id":   inv.OrderID,
		"pdf_url":    inv.PdfURL,
		"generated_at": time.Now().UTC().Format(time.RFC3339),
	})

	_ = ch.PublishWithContext(ctx,
		"",
		"invoice_generation",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

// generateInvoiceNo — INV-YYYYMMDD-NNNN
func generateInvoiceNo() string {
	now := time.Now()
	randPart := fmt.Sprintf("%04d", rand.Int63n(10000))
	return fmt.Sprintf("INV-%s-%s", now.Format("20060102"), randPart)
}

// formatAddress — 把 JSONMap 地址格式化成多行字符串
func formatAddress(addr models.JSONMap) string {
	if addr == nil {
		return ""
	}
	var parts []string
	for _, key := range []string{"line1", "line2", "city", "county", "postcode", "country", "street", "state", "zip"} {
		if v, ok := addr[key]; ok {
			if s, ok := v.(string); ok && s != "" {
				parts = append(parts, s)
			}
		}
	}
	if len(parts) > 0 {
		// 尝试用 , 拼接
		out := ""
		for i, p := range parts {
			if i > 0 {
				out += ", "
			}
			out += p
		}
		return out
	}
	// fallback：序列化整个 map
	b, _ := json.Marshal(addr)
	return string(b)
}

