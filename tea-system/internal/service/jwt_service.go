package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"tea-system/internal/middleware"
)

// JWTService — JWT 签发 + 验证 + Refresh
type JWTService struct {
	secret          string
	staffExpireMin  int
	userExpireMin   int
	refreshExpireDay int
}

func NewJWTService(secret string, staffExpireMin, userExpireMin, refreshExpireDay int) *JWTService {
	return &JWTService{
		secret:           secret,
		staffExpireMin:   staffExpireMin,
		userExpireMin:    userExpireMin,
		refreshExpireDay: refreshExpireDay,
	}
}

// GenerateStaffToken — 给 Staff 签发 access token
func (s *JWTService) GenerateStaffToken(staffID uint64, email, role string) (string, error) {
	now := time.Now()
	claims := &middleware.AuthClaims{
		SubjectType: "staff",
		SubjectID:   staffID,
		Email:       email,
		Role:        role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(s.staffExpireMin) * time.Minute)),
			Issuer:    "tea-system",
			Audience:  jwt.ClaimStrings{"staff-api"},
		},
	}
	return s.sign(claims)
}

// GenerateUserToken — 给 User 签发 access token
func (s *JWTService) GenerateUserToken(userID uint64, email string) (string, error) {
	now := time.Now()
	claims := &middleware.AuthClaims{
		SubjectType: "user",
		SubjectID:   userID,
		Email:       email,
		Role:        "customer",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(s.userExpireMin) * time.Minute)),
			Issuer:    "tea-system",
			Audience:  jwt.ClaimStrings{"public-site"},
		},
	}
	return s.sign(claims)
}

// GenerateRefreshToken — 长时效 refresh token
func (s *JWTService) GenerateRefreshToken(subjectType string, subjectID uint64) (string, error) {
	now := time.Now()
	claims := &middleware.AuthClaims{
		SubjectType: subjectType,
		SubjectID:   subjectID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(s.refreshExpireDay) * 24 * time.Hour)),
			Issuer:    "tea-system",
			Subject:   fmt.Sprintf("%s:%d", subjectType, subjectID),
		},
	}
	return s.sign(claims)
}

func (s *JWTService) sign(claims *middleware.AuthClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

// Parse — 解析并验证 token
func (s *JWTService) Parse(tokenString string) (*middleware.AuthClaims, error) {
	claims := &middleware.AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
