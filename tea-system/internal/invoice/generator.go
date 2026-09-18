// Package invoice — 发票生成器（PDF）
package invoice

import (
	"bytes"
	"fmt"
	"time"

	"github.com/go-pdf/fpdf"
)

// Company — 公司抬头（可以通过 config 注入）
type Company struct {
	Name      string
	Address1  string
	Address2  string
	TaxNumber string
	Email     string
	Website   string
}

// DefaultCompany — Our Tea Co., London, UK
var DefaultCompany = Company{
	Name:      "Our Tea Co.",
	Address1:  "123 Tea Lane",
	Address2:  "London, W1A 1AA, United Kingdom",
	TaxNumber: "GB 123 4567 89",
	Email:     "hello@ourtea.co",
	Website:   "www.ourtea.co",
}

// InvoiceItem — 发票行
type InvoiceItem struct {
	Name      string
	Quantity  int
	UnitPrice float64
}

// Input — 生成发票所需的完整数据
type Input struct {
	Company         Company
	InvoiceNo       string
	IssueDate       time.Time
	DueDate         time.Time
	OrderNo         string
	Currency        string
	HsCode          string
	CountryOfOrigin string
	ShippingCost    float64
	CustomerName    string
	CustomerEmail   string
	CustomerAddress string
	Items           []InvoiceItem
	Notes           string
}

// Generate — 生成 A4 PDF，返回 []byte
func Generate(in Input) ([]byte, error) {
	if in.Currency == "" {
		in.Currency = "GBP"
	}
	if in.Company.Name == "" {
		in.Company = DefaultCompany
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 15, 20)
	pdf.SetAutoPageBreak(true, 20)
	pdf.SetDisplayMode("fullpage", "TwoPageRight")

	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 10)

	_, pageH := pdf.GetPageSize()
	leftMargin, _, _, _ := pdf.GetMargins()
	usableW := pageH - leftMargin // 页面宽（mm），A4 P 模式 210mm

	startX := leftMargin

	// ============ 公司抬头 ============
	pdf.SetFont("Helvetica", "B", 20)
	pdf.SetXY(startX, 15)
	pdf.CellFormat(100, 8, in.Company.Name, "", 0, "L", false, 0, "")
	pdf.Ln(12)

	pdf.SetFont("Helvetica", "", 9)
	pdf.CellFormat(100, 4, in.Company.Address1, "", 0, "L", false, 0, "")
	pdf.Ln(4)
	pdf.CellFormat(100, 4, in.Company.Address2, "", 0, "L", false, 0, "")
	pdf.Ln(4)
	if in.Company.TaxNumber != "" {
		pdf.CellFormat(100, 4, "VAT No: "+in.Company.TaxNumber, "", 0, "L", false, 0, "")
		pdf.Ln(4)
	}
	if in.Company.Email != "" {
		pdf.CellFormat(100, 4, in.Company.Email, "", 0, "L", false, 0, "")
		pdf.Ln(4)
	}
	if in.Company.Website != "" {
		pdf.CellFormat(100, 4, in.Company.Website, "", 0, "L", false, 0, "")
		pdf.Ln(4)
	}

	// 右对齐 Invoice 标题
	pdf.SetXY(startX, pdf.GetY()-28)
	pdf.SetFont("Helvetica", "B", 16)
	pdf.CellFormat(0, 8, "COMMERCIAL INVOICE", "", 0, "R", false, 0, "")
	pdf.Ln(16)

	// ============ 发票号 + 日期 ============
	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(90, 6, fmt.Sprintf("Invoice No: %s", in.InvoiceNo), "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf("Date: %s", formatUKDate(in.IssueDate)), "", 1, "R", false, 0, "")
	if !in.DueDate.IsZero() {
		pdf.CellFormat(90, 6, "", "", 0, "L", false, 0, "")
		pdf.CellFormat(0, 6, fmt.Sprintf("Due: %s", formatUKDate(in.DueDate)), "", 1, "R", false, 0, "")
	}
	if in.OrderNo != "" {
		pdf.CellFormat(90, 6, fmt.Sprintf("Order No: %s", in.OrderNo), "", 0, "L", false, 0, "")
		pdf.Ln(8)
	}

	// ============ Customer 信息 ============
	pdf.SetFont("Helvetica", "B", 11)
	pdf.CellFormat(40, 6, "Bill To:", "", 0, "L", false, 0, "")
	pdf.Ln(6)
	pdf.SetFont("Helvetica", "", 10)

	// 客户名字（自动换行）
	customerBlock := in.CustomerName
	if in.CustomerEmail != "" {
		customerBlock += "\n" + in.CustomerEmail
	}
	if in.CustomerAddress != "" {
		customerBlock += "\n" + in.CustomerAddress
	}
	pdf.MultiCell(0, 5, customerBlock, "", "L", false)
	pdf.Ln(4)

	// ============ 商品表格 ============
	colName := 80.0
	colQty := 20.0
	colUnit := 35.0
	colTotal := 35.0
	rowH := 7.0

	// 表头
	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetFillColor(230, 230, 230)
	pdf.SetDrawColor(180, 180, 180)
	yHead := pdf.GetY()
	xHead := startX
	pdf.SetXY(xHead, yHead)
	pdf.CellFormat(colName, rowH, "Product", "1", 0, "L", true, 0, "")
	pdf.CellFormat(colQty, rowH, "Qty", "1", 0, "C", true, 0, "")
	pdf.CellFormat(colUnit, rowH, "Unit Price", "1", 0, "R", true, 0, "")
	pdf.CellFormat(colTotal, rowH, "Total", "1", 1, "R", true, 0, "")

	// 表体
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetFillColor(255, 255, 255)
	pdf.SetDrawColor(180, 180, 180)

	var subtotal float64
	for _, item := range in.Items {
		itemTotal := float64(item.Quantity) * item.UnitPrice
		subtotal += itemTotal

		yStart := pdf.GetY()
		xStart := leftMargin

		// 产品名：用 MultiCell 允许换行
		pdf.SetXY(xStart, yStart)
		pdf.MultiCell(colName, rowH, truncate(item.Name, 80), "1", "L", false)

		// 回到行起点，画其他列（覆盖 MultiCell 的推进）
		pdf.SetY(yStart)
		xAfterName := leftMargin + colName
		pdf.SetX(xAfterName)
		pdf.CellFormat(colQty, rowH, fmt.Sprintf("%d", item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colUnit, rowH, fmtPrice(item.UnitPrice, in.Currency), "1", 0, "R", false, 0, "")
		pdf.CellFormat(colTotal, rowH, fmtPrice(itemTotal, in.Currency), "1", 0, "R", false, 0, "")

		// 推到下一行（取 MultiCell 推进后的 Y 和 Cell 的 Y 里较大的）
		_ = xAfterName
	}

	// ============ 小计 + 运费 + 合计 ============
	pdf.Ln(4)
	pdf.SetFont("Helvetica", "", 10)
	totalW := colQty + colUnit + colTotal
	rightX := usableW - totalW - leftMargin + leftMargin // 右边块起始 X（留边距）

	pdf.SetX(rightX)
	pdf.CellFormat(totalW-colTotal, rowH, "Subtotal:", "", 0, "R", false, 0, "")
	pdf.CellFormat(colTotal, rowH, fmtPrice(subtotal, in.Currency), "", 1, "R", false, 0, "")

	if in.ShippingCost > 0 {
		pdf.SetX(rightX)
		pdf.CellFormat(totalW-colTotal, rowH, "Shipping:", "", 0, "R", false, 0, "")
		pdf.CellFormat(colTotal, rowH, fmtPrice(in.ShippingCost, in.Currency), "", 1, "R", false, 0, "")
	}

	grandTotal := subtotal + in.ShippingCost
	pdf.Ln(3)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetX(rightX)
	pdf.CellFormat(totalW-colTotal, rowH, "TOTAL:", "T", 0, "R", false, 0, "")
	pdf.CellFormat(colTotal, rowH, fmtPrice(grandTotal, in.Currency), "T", 1, "R", false, 0, "")

	// ============ 清关信息 ============
	pdf.Ln(6)
	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(60, 6, "Customs & Export Information", "", 1, "L", false, 0, "")
	pdf.Ln(1)
	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(60, 5, fmt.Sprintf("HS Code: %s", in.HsCode), "", 1, "L", false, 0, "")
	pdf.CellFormat(80, 5, fmt.Sprintf("Country of Origin: %s", in.CountryOfOrigin), "", 1, "L", false, 0, "")
	pdf.CellFormat(60, 5, fmt.Sprintf("Currency: %s", in.Currency), "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// ============ Notes ============
	if in.Notes != "" {
		pdf.SetFont("Helvetica", "B", 10)
		pdf.CellFormat(30, 6, "Notes:", "", 1, "L", false, 0, "")
		pdf.SetFont("Helvetica", "", 9)
		pdf.MultiCell(0, 5, in.Notes, "", "L", false)
	}

	// ============ 页脚 ============
	// 手动画页脚
	pdf.SetFont("Helvetica", "I", 8)
	pdf.SetTextColor(120, 120, 120)
	pdf.SetXY(leftMargin, pageH-12)
	pdf.CellFormat(0, 5, "Thank you for your business. This is a commercial invoice for customs purposes.", "", 0, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)

	// ============ 输出 ============
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// formatUKDate — 英式日期 DD/MM/YYYY
func formatUKDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return fmt.Sprintf("%02d/%02d/%04d", t.Day(), int(t.Month()), t.Year())
}

// fmtPrice — 格式化金额
func fmtPrice(v float64, currency string) string {
	sym := currencySymbol(currency)
	return fmt.Sprintf("%s %.2f", sym, v)
}

func currencySymbol(currency string) string {
	switch currency {
	case "GBP":
		return "\u00a3" // £
	case "USD":
		return "$"
	case "EUR":
		return "\u20ac" // €
	case "CNY", "RMB":
		return "\u00a5" // ¥
	default:
		return currency
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}
