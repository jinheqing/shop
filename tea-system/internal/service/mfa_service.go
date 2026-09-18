package service

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/pquerna/otp/totp"

	"tea-system/internal/models"
	"tea-system/internal/repository"
)

// MFAService — TOTP 生成 + 验证
type MFAService struct {
	staffRepo *repository.StaffRepo
}

func NewMFAService(staffRepo *repository.StaffRepo) *MFAService {
	return &MFAService{staffRepo: staffRepo}
}

// GenerateSecret — 生成新的 TOTP secret
func (s *MFAService) GenerateSecret() (string, error) {
	// 用随机字节生成一个 20 字节的 base32 secret
	buf := make([]byte, 20)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(buf), nil
}

// GenerateTOTP — 生成当前时刻的 TOTP（用于用户手动测试）
func (s *MFAService) GenerateTOTP(secret string) (string, error) {
	return totp.GenerateCode(secret, time.Now())
}

// Verify — 验证用户输入的 TOTP code
func (s *MFAService) Verify(secret, code string) bool {
	if secret == "" {
		return false
	}
	// 允许 ±1 个时间窗口（30s × 1 = ±30s 容差）
	return totp.Validate(code, secret)
}

// GenerateQRURI — 生成 Google Authenticator 扫码 URI
// 格式: otpauth://totp/Tea%20System:admin@ourdomain.com?secret=XXX&issuer=Tea%20System
func (s *MFAService) GenerateQRURI(staff *models.Staff) string {
	return fmt.Sprintf(
		"otpauth://totp/Tea%%20System:%s?secret=%s&issuer=Tea%%20System",
		staff.Email,
		staff.MfaSecret,
	)
}

// Fingerprint — 给 secret 生成一个短 fingerprint 用于调试日志
func (s *MFAService) Fingerprint(secret string) string {
	h := sha1.Sum([]byte(secret))
	return hex.EncodeToString(h[:4])
}
