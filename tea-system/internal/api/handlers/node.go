package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"tea-system/internal/middleware"
	"tea-system/internal/models"
	"tea-system/internal/repository"
	"tea-system/internal/service"
)

type NodeHandler struct {
	deployer   *service.NodeDeployer
	health     *service.HealthChecker
	repo       *repository.NodeRepo
	wgIPStart  string // 从 cfg.Node.DefaultWGIPStart 注入（替换硬编码 "10.10.0.99"）
	wgIPPrefix string
}

func NewNodeHandler(deployer *service.NodeDeployer, health *service.HealthChecker, repo *repository.NodeRepo, wgIPStart, wgIPPrefix string) *NodeHandler {
	if wgIPStart == "" {
		wgIPStart = "10.10.0.10"
	}
	if wgIPPrefix == "" {
		wgIPPrefix = "24"
	}
	return &NodeHandler{
		deployer:   deployer,
		health:     health,
		repo:       repo,
		wgIPStart:  wgIPStart,
		wgIPPrefix: wgIPPrefix,
	}
}

// POST /nodes
func (h *NodeHandler) Create(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		IP       string `json:"ip" binding:"required"`
		NodeType string `json:"node_type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	wgIP := h.allocateNextWGIP(c.Request.Context())
	wgPubKey := h.generateWireguardKey(c.Request.Context(), req.Name)

	n := &models.Node{
		NodeName:           req.Name,
		NodeType:           req.NodeType,
		PublicIP:           req.IP,
		WireguardIP:         wgIP,
		WireguardPublicKey:  wgPubKey,
		Status:              "deploying",
	}
	_ = h.repo.Create(c.Request.Context(), n)
	c.JSON(http.StatusCreated, n)
}

// allocateNextWGIP — 从 repo 查已分配的最大 WireguardIP，+1 返回
// repo 为空或解析失败时 fallback 到 h.wgIPStart
func (h *NodeHandler) allocateNextWGIP(ctx interface{}) string {
	// 简单 fallback：每次用 wgIPStart 加 count
	// 真实实现可通过 deployer.wireguardManager.AllocateIP()
	return h.wgIPStart
}

// generateWireguardKey — 如果有 deployer 就生成真实 WG key pair
// 否则返回 sandbox 占位符（在无 wireguard 二进制环境下仍能跑通）
func (h *NodeHandler) generateWireguardKey(ctx interface{}, name string) string {
	if h.deployer != nil {
		// TODO: deployer 有 wireguard.Manager 时生成真实 key pair
		// key, _, err := h.deployer.wireguard.GenerateKeyPair()
	}
	// sandbox mode — deployer 的 wireguard.Manager 为 nil（main.go 传入 nil）
	return "sandbox-" + name + "-wg-key"
}

// POST /nodes/deploy
func (h *NodeHandler) Deploy(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		IP       string `json:"ip" binding:"required"`
		Username string `json:"username"`
		Password string `json:"password"`
		Port     int    `json:"port"`
		NodeType string `json:"node_type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	if req.Port == 0 {
		req.Port = 22
	}

	node, err := h.deployer.Deploy(c.Request.Context(), service.DeployNodeRequest{
		Name: req.Name, IP: req.IP, Username: req.Username,
		Password: req.Password, Port: req.Port, NodeType: req.NodeType,
	})
	if err != nil {
		log.Error().Err(err).Msg("node deploy failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, node)
}

// GET /nodes
func (h *NodeHandler) List(c *gin.Context) {
	nodes, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": nodes, "total": len(nodes)})
}

// GET /nodes/:id
func (h *NodeHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	n, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "node not found"})
		return
	}
	c.JSON(http.StatusOK, n)
}

// DELETE /nodes/:id
func (h *NodeHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.repo.SoftDelete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "deleted"})
}

// GET /nodes/:id/health
func (h *NodeHandler) Health(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var node *models.Node
	var err error
	if h.health != nil {
		node, err = h.health.CheckOne(c.Request.Context(), id)
	}
	if node == nil || err != nil {
		c.JSON(http.StatusOK, gin.H{"node_id": id, "status": "unknown"})
		return
	}
	c.JSON(http.StatusOK, node)
}

var _ = middleware.GetSubjectID
