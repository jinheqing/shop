package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// MagicLinkService — 用户魔法链接（magic link）认证
// 128 位 URL-safe token，10 分钟有效，单次使用
type MagicLinkService struct {
	rdb          *redis.Client
	MailSvc      *MailService // 公开：handlers 需要用它发邮件
	domain       string       // 魔法链接域名，如 https://www.ourdomain.com
	expireWindow time.Duration
}

func NewMagicLinkService(rdb *redis.Client, mailSvc *MailService, domain string) *MagicLinkService {
	return &MagicLinkService{
		rdb:          rdb,
		MailSvc:      mailSvc,
		domain:       domain,
		expireWindow: 10 * time.Minute,
	}
}

var (
	ErrMagicLinkNotFound = errors.New("magic link not found or expired")
	ErrMagicLinkUsed     = errors.New("magic link already used")
)

// Generate — 生成 128 位 URL-safe token，存 Redis
// 返回 token + magic link URL
func (s *MagicLinkService) Generate(ctx context.Context, userID uint64, email string) (token string, magicURL string, err error) {
	// 1. 生成 128 位随机（16 字节 base64url ≈ 22 字符，取 24 字节 = 192 位 ≈ 32 字符 URL-safe）
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", "", fmt.Errorf("token gen: %w", err)
	}
	// URL-safe base64 编码 + 去掉 padding
	token = base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(buf)
	// 确保 128 位熵（24 字节 = 192 位，足够了）

	// 2. 存 Redis: key = magic:<token>, value = user:id:email, TTL = 10min
	key := fmt.Sprintf("magic:%s", token)
	value := fmt.Sprintf("%d:%s", userID, email)
	if err := s.rdb.Set(ctx, key, value, s.expireWindow).Err(); err != nil {
		return "", "", fmt.Errorf("redis set: %w", err)
	}

	// 3. 生成完整 URL
	magicURL = fmt.Sprintf("%s/auth/magic?token=%s", s.domain, token)

	// 4. 记录日志（不打印 token）
	log.Info().Str("email", email).Str("token_fp", hex.EncodeToString(buf[:4])).Msg("magic link generated")

	return token, magicURL, nil
}

// Verify — 验证 token，返回 user info，验证后立即删除（单次使用）
func (s *MagicLinkService) Verify(ctx context.Context, token string) (userID uint64, email string, err error) {
	key := fmt.Sprintf("magic:%s", token)

	// 1. 原子 GET + DEL（Lua 脚本）
	// 防止并发重复使用
	script := `
local val = redis.call('GET', KEYS[1])
if val then
  redis.call('DEL', KEYS[1])
  return val
end
return nil
`
	result, err := s.rdb.Eval(ctx, script, []string{key}).Result()
	if err != nil {
		return 0, "", err
	}
	if result == nil {
		return 0, "", ErrMagicLinkNotFound
	}

	// 2. 解析 "userID:email"
	val := result.(string)
	var uid uint64
	var em string
	if _, err := fmt.Sscanf(val, "%d:%s", &uid, &em); err != nil {
		return 0, "", fmt.Errorf("corrupted magic link data")
	}

	log.Info().Str("email", em).Msg("magic link verified + consumed")
	return uid, em, nil
}
