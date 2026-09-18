package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// UploadHandler — 文件上传（通用：图片/视频/附件）
type UploadHandler struct {
	storageDir string
	publicURL  string
}

// NewUploadHandler — storageDir: 本地存储目录；publicURL: 对外访问基地址，如 /uploads
func NewUploadHandler(storageDir, publicURL string) *UploadHandler {
	_ = os.MkdirAll(storageDir, 0o755)
	return &UploadHandler{storageDir: storageDir, publicURL: publicURL}
}

type uploadLimit struct {
	maxSize     int64
	allowedExts map[string]bool
}

var uploadRules = map[string]uploadLimit{
	"image": {
		maxSize:     10 * 1024 * 1024, // 10MB
		allowedExts: map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true},
	},
	"video": {
		maxSize:     100 * 1024 * 1024, // 100MB
		allowedExts: map[string]bool{".mp4": true, ".webm": true, ".mov": true},
	},
	"file": {
		maxSize:     20 * 1024 * 1024, // 20MB
		allowedExts: map[string]bool{".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true, ".txt": true, ".csv": true, ".zip": true},
	},
}

// Upload — POST /api/v1/upload?type=image|video|file
func (h *UploadHandler) Upload(c *gin.Context) {
	fileType := c.DefaultQuery("type", "file")
	if fileType != "image" && fileType != "video" && fileType != "file" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid type: image | video | file"})
		return
	}

	rule := uploadRules[fileType]
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, rule.maxSize+1)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "no file or too large"})
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !rule.allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "file type not allowed"})
		return
	}

	// 生成随机文件名
	var b [12]byte
	_, _ = rand.Read(b[:])
	filename := fmt.Sprintf("%s-%s%s", time.Now().Format("20060102"), hex.EncodeToString(b[:]), ext)

	// 按类型分子目录
	subDir := filepath.Join(h.storageDir, fileType)
	_ = os.MkdirAll(subDir, 0o755)

	dstPath := filepath.Join(subDir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		log.Error().Err(err).Str("path", dstPath).Msg("failed to create upload file")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "storage error"})
		return
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "write failed"})
		return
	}

	publicURL := fmt.Sprintf("%s/%s/%s", h.publicURL, fileType, filename)

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "uploaded",
		"file": gin.H{
			"url":  publicURL,
			"type": fileType,
			"name": header.Filename,
			"size": size,
			"mime": header.Header.Get("Content-Type"),
		},
	})
}
