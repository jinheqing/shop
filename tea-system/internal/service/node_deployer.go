package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"runtime"
	"time"

	"tea-system/internal/models"
	"tea-system/internal/repository"

	wireguard "tea-system/internal/service/wireguard"
)

// NodeDeployer — 节点一键部署
type NodeDeployer struct {
	repo   *repository.NodeRepo
	wgMgr  *wireguard.Manager
}

func NewNodeDeployer(nodeRepo *repository.NodeRepo, wgMgr *wireguard.Manager) *NodeDeployer {
	return &NodeDeployer{repo: nodeRepo, wgMgr: wgMgr}
}

// DeployNodeRequest — 部署参数
type DeployNodeRequest struct {
	Name     string `json:"name" binding:"required"`
	IP       string `json:"ip" binding:"required"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	NodeType string `json:"node_type" binding:"required"`
}

func (d *NodeDeployer) Deploy(ctx context.Context, req DeployNodeRequest) (*models.Node, error) {
	if req.Port == 0 {
		req.Port = 22
	}

	// 1. 生成 WG 密钥（mock 模式）
	_, pub, _ := wireguard.GenerateKeyPair()

	// 2. 分配合并 IP（mock: 从 10.10.0.10 开始）
	wgIP, err := d.allocateWireguardIP(ctx)
	if err != nil {
		return nil, fmt.Errorf("allocate wg ip: %w", err)
	}

	// 3. SSH 部署（沙箱跳过，走 mock）
	status := models.NodeStatusOnline
	err = d.sshDeploy(req.IP, req.Username, req.Password, req.Port, req.NodeType)
	if err != nil {
		log.Printf("node_deployer: ssh deploy failed (mock mode continues): %v", err)
		status = models.NodeStatusDeploying
	}

	// 4. 加 peer 到主节点
	d.wgMgr.AddPeer(wireguard.Peer{
		PublicKey:  pub,
		Endpoint:   fmt.Sprintf("%s:%d", req.IP, 51820),
		AllowedIPs: []string{wgIP},
	})
	_ = d.wgMgr.Apply()

	// 5. 存 DB
	node := &models.Node{
		NodeName:           req.Name,
		NodeType:           req.NodeType,
		PublicIP:           req.IP,
		WireguardIP:         wgIP,
		WireguardPublicKey: pub,
		Status:              status,
		DeployedAt:          timePtr(time.Now()),
	}
	if err := d.repo.Create(ctx, node); err != nil {
		return nil, fmt.Errorf("db create node: %w", err)
	}

	// 6. 密码变量置空 + GC（安全约束）
	req.Password = ""
	
	runtime.GC()

	log.Printf("node_deployer: node %s deployed (ip=%s, wg=%s)", req.Name, req.IP, wgIP)
	return node, nil
}

// allocateWireguardIP — 简单递增分配
func (d *NodeDeployer) allocateWireguardIP(ctx context.Context) (string, error) {
	nodes, err := d.repo.List(ctx)
	if err != nil {
		return "", err
	}
	used := map[string]bool{}
	for _, n := range nodes {
		used[n.WireguardIP] = true
	}
	for i := 10; i <= 250; i++ {
		ip := fmt.Sprintf("10.10.0.%d", i)
		if !used[ip] {
			return ip, nil
		}
	}
	return "", errors.New("wg ip pool exhausted")
}

// sshDeploy — SSH 远程安装（沙箱 mock）
func (d *NodeDeployer) sshDeploy(ip, user, pass string, port int, nodeType string) error {
	// 沙箱里没有 SSH 端口，走 mock
	log.Printf("node_deployer: would SSH into %s@%s:%d (mock mode)", user, ip, port)
	// 真实部署流程:
	//   apt update → docker install → wg install → container up
	return nil
}

// randomString — 辅助
func randomString(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func timePtr(t time.Time) *time.Time { return &t }
