package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"tea-system/internal/im"
	"tea-system/internal/middleware"
	"tea-system/internal/repository"
	"tea-system/internal/service"
)

// IMWSHandler — WebSocket 入口
type IMWSHandler struct {
	hub          *im.Hub
	messageRepo  *repository.MessageRepo
	translateSvc *service.TranslateService
	jwtSecret    string
}

func NewIMWSHandler(
	hub *im.Hub,
	messageRepo *repository.MessageRepo,
	translateSvc *service.TranslateService,
	jwtSecret string,
) *IMWSHandler {
	return &IMWSHandler{
		hub:          hub,
		messageRepo:  messageRepo,
		translateSvc: translateSvc,
		jwtSecret:    jwtSecret,
	}
}

// Serve — GET /ws/im?conversation_id=1&token=xxx
// 不走 JWTAuth 中间件，自己解析 JWT（从 query token 或 Authorization header）
func (h *IMWSHandler) Serve(c *gin.Context) {
	tokenString := ""

	// 1. 优先 query token
	if q := c.Query("token"); q != "" {
		tokenString = q
	}
	// 2. 次选 Authorization header
	if tokenString == "" {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			tokenString = authHeader
		}
	}

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "missing token"})
		return
	}

	// 解析 JWT
	claims := &middleware.AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(h.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "invalid token"})
		return
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "token expired"})
		return
	}

	// 升级为 WebSocket
	conn, err := im.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Warn().Err(err).Msg("ws upgrade failed")
		return
	}

	client := im.NewClient(
		h.hub,
		conn,
		claims.SubjectID,
		claims.SubjectType,
		h.messageRepo,
		h.translateSvc,
	)

	// 注册到 hub
	h.hub.Register(client)

	// 启动读写 pump（阻塞直到连接断开）
	go client.WritePump()
	go client.ReadPump()
}
