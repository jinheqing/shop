package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// PaymentGateway — 支持的支付网关
type PaymentGateway string

const (
	PaymentGateway2Checkout PaymentGateway = "2checkout"
	PaymentGatewayPayPal    PaymentGateway = "paypal"
)

// PaymentService — 支付网关服务（mock 模式，暂不接真 SDK）
// 生产环境替换为真实 2Checkout / PayPal SDK 调用
type PaymentService struct {
	BaseURL   string
	APIKey    string
	APISecret string
}

func NewPaymentService(baseURL, apiKey, apiSecret string) *PaymentService {
	return &PaymentService{
		BaseURL:   baseURL,
		APIKey:    apiKey,
		APISecret: apiSecret,
	}
}

// InitPaymentRequest — 发起支付请求
type InitPaymentRequest struct {
	OrderNo  string
	Amount   float64
	Currency string
	ReturnURL string
}

// InitPaymentResult — 返回给前端
type InitPaymentResult struct {
	Gateway      PaymentGateway
	CheckoutURL  string
	MerchantRef  string
}

// InitPayment — 按网关类型初始化支付（统一入口）
func (s *PaymentService) InitPayment(gateway PaymentGateway, req InitPaymentRequest) (*InitPaymentResult, error) {
	switch gateway {
	case PaymentGateway2Checkout:
		return s.Init2CheckoutPayment(req)
	case PaymentGatewayPayPal:
		return s.InitPayPalPayment(req)
	default:
		return nil, fmt.Errorf("unsupported payment gateway: %s", gateway)
	}
}

// Init2CheckoutPayment — mock 2Checkout 支付初始化
func (s *PaymentService) Init2CheckoutPayment(req InitPaymentRequest) (*InitPaymentResult, error) {
	if req.OrderNo == "" || req.Amount <= 0 {
		return nil, errors.New("order_no and valid amount required")
	}
	if req.Currency == "" {
		req.Currency = "GBP"
	}

	mockTxID := generateMockTxID("2CO")
	checkoutURL := fmt.Sprintf(
		"%s/mock/2checkout/preview?order_no=%s&amount=%.2f&currency=%s&return=%s&tx=%s",
		s.BaseURL, req.OrderNo, req.Amount, req.Currency, req.ReturnURL, mockTxID,
	)

	return &InitPaymentResult{
		Gateway:     PaymentGateway2Checkout,
		CheckoutURL: checkoutURL,
		MerchantRef: req.OrderNo,
	}, nil
}

// InitPayPalPayment — mock PayPal 支付初始化
func (s *PaymentService) InitPayPalPayment(req InitPaymentRequest) (*InitPaymentResult, error) {
	if req.OrderNo == "" || req.Amount <= 0 {
		return nil, errors.New("order_no and valid amount required")
	}
	if req.Currency == "" {
		req.Currency = "GBP"
	}

	mockTxID := generateMockTxID("PP")
	checkoutURL := fmt.Sprintf(
		"%s/mock/paypal/checkout?order_no=%s&amount=%.2f&currency=%s&return=%s&tx=%s",
		s.BaseURL, req.OrderNo, req.Amount, req.Currency, req.ReturnURL, mockTxID,
	)

	return &InitPaymentResult{
		Gateway:     PaymentGatewayPayPal,
		CheckoutURL: checkoutURL,
		MerchantRef: req.OrderNo,
	}, nil
}

// ============================================================
// 回调处理（幂等入口：调用方先用 GatewayTransactionID 查 UNIQUE）
// ============================================================

// CallbackResult — 回调解析结果
type CallbackResult struct {
	GatewayTransactionID string
	OrderNo              string
	Amount               float64
	Currency             string
	Status               string // success / failed
	RawBody              map[string]interface{}
}

// HandleCallback — 按网关统一分发
func (s *PaymentService) HandleCallback(gateway PaymentGateway, body []byte) (*CallbackResult, error) {
	switch gateway {
	case PaymentGateway2Checkout:
		return s.Handle2CheckoutCallback(body)
	case PaymentGatewayPayPal:
		return s.HandlePayPalCallback(body)
	default:
		return nil, fmt.Errorf("unsupported payment gateway: %s", gateway)
	}
}

// Handle2CheckoutCallback — 处理 2Checkout INS 回调（mock 实现）
//
// 真实 2CO INS 是 x-www-form-urlencoded，mock 支持 JSON 或 form 两种：
//
//	{"order_number":"ORD-...","purchase_order_id":"MERCHANT_REF","invoice_list_amount":123.45,"is_valid":"1"}
func (s *PaymentService) Handle2CheckoutCallback(body []byte) (*CallbackResult, error) {
	raw := parseBody(body)

	// 真实 2CO 用 invoice_list_amount + order_number + is_valid
	var (
		orderNo string
		txID    string
		amount  float64
		currency = "GBP"
		status  = "pending"
	)

	if v, ok := raw["purchase_order_id"].(string); ok && v != "" {
		orderNo = v
	}
	if v, ok := raw["order_number"].(string); ok && v != "" {
		txID = v
	}
	if txID == "" {
		// mock 兼容：tx 字段
		if v, ok := raw["tx"].(string); ok {
			txID = v
		}
	}
	if v, ok := raw["invoice_list_amount"]; ok {
		amount = toFloat(v)
	} else if v, ok := raw["amount"]; ok {
		amount = toFloat(v)
	}
	if v, ok := raw["currency"].(string); ok && v != "" {
		currency = v
	}

	// 判定状态：真实 2CO 看 is_valid=="1" + fraud_status=="approved"
	if v, ok := raw["is_valid"].(string); ok && v == "1" {
		status = "success"
	} else if v, ok := raw["status"].(string); ok {
		if strings.EqualFold(v, "success") || strings.EqualFold(v, "approved") {
			status = "success"
		} else if strings.EqualFold(v, "failed") || strings.EqualFold(v, "declined") {
			status = "failed"
		}
	}

	if txID == "" {
		return nil, errors.New("2checkout callback missing gateway transaction id (order_number/tx)")
	}

	return &CallbackResult{
		GatewayTransactionID: txID,
		OrderNo:              orderNo,
		Amount:               amount,
		Currency:             currency,
		Status:               status,
		RawBody:              raw,
	}, nil
}

// HandlePayPalCallback — 处理 PayPal webhook 回调（mock 实现）
//
// 真实 PayPal webhook 是 JSON，主要字段：event_type / resource.sale / resource.amount.value
func (s *PaymentService) HandlePayPalCallback(body []byte) (*CallbackResult, error) {
	raw := parseBody(body)

	eventType, _ := raw["event_type"].(string)

	var (
		txID     string
		orderNo  string
		amount   float64
		currency = "GBP"
		status   = "pending"
	)

	// 从 resource 里取
	if res, ok := raw["resource"].(map[string]interface{}); ok {
		if v, ok := res["sale"].(map[string]interface{}); ok {
			if v2, ok := v["id"].(string); ok {
				txID = v2
			}
			if v2, ok := v["custom"].(string); ok && v2 != "" {
				orderNo = v2
			} else if v2, ok := v["invoice_number"].(string); ok {
				orderNo = v2
			}
		}
		if amt, ok := res["amount"].(map[string]interface{}); ok {
			amount = toFloat(amt["value"])
			if v, ok := amt["currency_code"].(string); ok && v != "" {
				currency = v
			}
		}
	}

	if txID == "" {
		// mock 兼容：顶层 tx 字段
		if v, ok := raw["tx"].(string); ok {
			txID = v
		}
	}
	if orderNo == "" {
		if v, ok := raw["custom"].(string); ok {
			orderNo = v
		}
	}
	if amount == 0 {
		if v, ok := raw["amount"]; ok {
			amount = toFloat(v)
		}
	}

	// 判定状态
	switch eventType {
	case "PAYMENT.SALE.COMPLETED", "CHECKOUT.ORDER.COMPLETED":
		status = "success"
	case "PAYMENT.SALE.DENIED", "PAYMENT.SALE.REFUNDED", "CHECKOUT.ORDER.APPROVAL.REVERSED":
		status = "failed"
	default:
		if v, ok := raw["status"].(string); ok {
			if strings.EqualFold(v, "success") || strings.EqualFold(v, "completed") || strings.EqualFold(v, "approved") {
				status = "success"
			} else if strings.EqualFold(v, "failed") || strings.EqualFold(v, "denied") || strings.EqualFold(v, "refunded") {
				status = "failed"
			}
		}
	}

	if txID == "" {
		return nil, errors.New("paypal callback missing gateway transaction id (resource.sale.id/tx)")
	}

	return &CallbackResult{
		GatewayTransactionID: txID,
		OrderNo:              orderNo,
		Amount:               amount,
		Currency:             currency,
		Status:               status,
		RawBody:              raw,
	}, nil
}

// ============================================================
// 内部工具
// ============================================================

func generateMockTxID(prefix string) string {
	randPart := fmt.Sprintf("%010d", rand.Int63n(1_000_000_0000))
	return fmt.Sprintf("%s-%s", prefix, randPart)
}

func parseBody(body []byte) map[string]interface{} {
	out := map[string]interface{}{}
	if len(body) == 0 {
		return out
	}
	// 先尝试 JSON
	if err := json.Unmarshal(body, &out); err == nil {
		return out
	}
	// 再尝试 form（简单 key=value 解析，不依赖标准库 url.Values 也可以）
	text := string(body)
	for _, pair := range strings.Split(text, "&") {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			out[kv[0]] = kv[1]
		}
	}
	return out
}

func toFloat(v interface{}) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case float32:
		return float64(t)
	case int:
		return float64(t)
	case int64:
		return float64(t)
	case string:
		f, _ := strconv.ParseFloat(t, 64)
		return f
	case json.Number:
		f, _ := t.Float64()
		return f
	default:
		return 0
	}
}
