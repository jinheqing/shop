package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// CustomProductService — 定制报价的工具方法（token / 版本 / 价格）
type CustomProductService struct{}

func NewCustomProductService() *CustomProductService {
	return &CustomProductService{}
}

// GenerateProductToken — 生成 64 字符 URL-safe 随机 token
// crypto/rand 生成 48 bytes → base64url 编码后刚好 64 字符（不补 =）
func (s *CustomProductService) GenerateProductToken() (string, error) {
	buf := make([]byte, 48)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate product token: %w", err)
	}
	token := base64.RawURLEncoding.EncodeToString(buf)
	if len(token) != 64 {
		return "", fmt.Errorf("unexpected token length: %d", len(token))
	}
	return token, nil
}

// BumpVersion — 版本号 +1
func (s *CustomProductService) BumpVersion(current int) int {
	if current < 0 {
		current = 0
	}
	return current + 1
}

// CalcTotalAmount — 价格计算: total = unit_price * quantity + shipping_cost
func (s *CustomProductService) CalcTotalAmount(unitPrice float64, quantity int, shippingCost float64) float64 {
	return unitPrice*float64(quantity) + shippingCost
}
