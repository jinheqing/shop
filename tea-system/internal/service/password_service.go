package service

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordService — bcrypt 哈希 + 验证
type PasswordService struct{}

func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// Hash — bcrypt 哈希（默认 cost 10）
func (s *PasswordService) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Verify — 验证密码
func (s *PasswordService) Verify(passwordHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil
}
