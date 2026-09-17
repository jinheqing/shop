package wireguard

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

// Peer — 远端 peer 配置对象
type Peer struct {
	PublicKey    string // 32 bytes base64
	IP           string // peer 的 WireGuard IP (CIDR)，如 "10.0.0.2/32"
	Endpoint     string // "ip:port"，可为空（仅出站）
	AllowedIPs   []string
	PersistentKeepalive int // 秒；0 表示不设置
}

// Manager — WireGuard 主节点配置管理（配置字符串拼装 + 可选 wg 命令）
// 沙箱模式下不真正执行 wg，只返回可落地的配置字符串。
type Manager struct {
	mu sync.Mutex

	InterfaceName string // 如 "wg0"
	Address       string // 主节点 WG IP/CIDR，如 "10.0.0.1/24"
	ListenPort    int
	PrivateKey    string // 主节点私钥
	PublicKey     string // 主节点公钥（由 private key 推导）
	DNS           string // 可选
	MTU           int    // 可选，默认 1280

	Peers []Peer

	// WireGuard 可执行文件路径；为空或 exec 失败会走 mock
	wgBin string
}

// NewManager — 创建主节点 WireGuard 管理器
// privateKey 为空时自动生成（mock 密钥）。
func NewManager(iface, address string, listenPort int, privateKey string) (*Manager, error) {
	if iface == "" {
		iface = "wg0"
	}
	if address == "" {
		address = "10.0.0.1/24"
	}
	if listenPort == 0 {
		listenPort = 51820
	}
	m := &Manager{
		InterfaceName: iface,
		Address:       address,
		ListenPort:    listenPort,
		MTU:           1280,
		wgBin:         "wg",
	}
	if privateKey != "" {
		m.PrivateKey = privateKey
		// 尝试用 wg pubkey 推导
		m.PublicKey = derivePublicKey(privateKey)
	} else {
		priv, pub := mockKeyPair()
		m.PrivateKey = priv
		m.PublicKey = pub
	}
	return m, nil
}

// GenerateKeyPair — 生成一对 WireGuard 密钥（优先用 wg genkey，否则 mock）
func GenerateKeyPair() (privateKey, publicKey string, err error) {
	if wgBinPath, err := exec.LookPath("wg"); err == nil {
		out, err := exec.Command(wgBinPath, "genkey").CombinedOutput()
		if err == nil {
			priv := strings.TrimSpace(string(out))
			pubCmd := exec.Command(wgBinPath, "pubkey")
			pubCmd.Stdin = strings.NewReader(priv)
			pubOut, err := pubCmd.CombinedOutput()
			if err == nil {
				return priv, strings.TrimSpace(string(pubOut)), nil
			}
		}
	}
	priv, pub := mockKeyPair()
	return priv, pub, nil
}

// AddPeer — 添加一个 peer
func (m *Manager) AddPeer(p Peer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Peers = append(m.Peers, p)
}

// RemovePeer — 按公钥移除 peer
func (m *Manager) RemovePeer(pubKey string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := m.Peers[:0]
	for _, p := range m.Peers {
		if p.PublicKey != pubKey {
			out = append(out, p)
		}
	}
	m.Peers = out
}

// BuildInterfaceConfig — 生成主节点 Interface section
func (m *Manager) BuildInterfaceConfig() string {
	var b strings.Builder
	b.WriteString("[Interface]\n")
	fmt.Fprintf(&b, "Address = %s\n", m.Address)
	fmt.Fprintf(&b, "ListenPort = %d\n", m.ListenPort)
	fmt.Fprintf(&b, "PrivateKey = %s\n", m.PrivateKey)
	if m.MTU > 0 {
		fmt.Fprintf(&b, "MTU = %d\n", m.MTU)
	}
	if m.DNS != "" {
		fmt.Fprintf(&b, "DNS = %s\n", m.DNS)
	}
	return b.String()
}

// BuildPeerConfig — 生成单个 peer 的 Peer section
func (p Peer) BuildPeerConfig() string {
	var b strings.Builder
	b.WriteString("\n[Peer]\n")
	fmt.Fprintf(&b, "PublicKey = %s\n", p.PublicKey)
	if p.Endpoint != "" {
		fmt.Fprintf(&b, "Endpoint = %s\n", p.Endpoint)
	}
	if len(p.AllowedIPs) > 0 {
		fmt.Fprintf(&b, "AllowedIPs = %s\n", strings.Join(p.AllowedIPs, ", "))
	} else if p.IP != "" {
		fmt.Fprintf(&b, "AllowedIPs = %s\n", p.IP)
	}
	if p.PersistentKeepalive > 0 {
		fmt.Fprintf(&b, "PersistentKeepalive = %d\n", p.PersistentKeepalive)
	}
	return b.String()
}

// BuildFullConfig — 生成完整的 wg0.conf 字符串
func (m *Manager) BuildFullConfig() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	var b strings.Builder
	b.WriteString(m.BuildInterfaceConfig())
	for _, p := range m.Peers {
		b.WriteString(p.BuildPeerConfig())
	}
	return b.String()
}

// Apply — 如果 wg 命令可用就尝试同步 peer（沙箱里一般会跳过）
func (m *Manager) Apply() error {
	if _, err := exec.LookPath(m.wgBin); err != nil {
		log.Warn().Msg("wireguard: wg binary not found, skipping Apply (mock mode)")
		return nil
	}
	log.Debug().Msg("wireguard: wg Apply() would sync peers here")
	return nil
}

// PublicKey — 主节点公钥
func (m *Manager) PublicKeyGetter() string { return m.PublicKey }

// mockKeyPair — 沙箱模式下生成假密钥（base64 32 bytes，但非真 WG 密钥）
func mockKeyPair() (string, string) {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	priv := base64.StdEncoding.EncodeToString(b)
	// 简单 mock：把 private 再做一次变换作为 public，便于区分
	pubBytes := make([]byte, 32)
	for i, v := range b {
		pubBytes[i] = v ^ 0x55
	}
	pub := base64.StdEncoding.EncodeToString(pubBytes)
	return priv, pub
}

// derivePublicKey — 优先尝试用 wg pubkey 推导，否则返回基于 privateKey 的 mock
func derivePublicKey(privateKey string) string {
	if wgBinPath, err := exec.LookPath("wg"); err == nil {
		cmd := exec.Command(wgBinPath, "pubkey")
		cmd.Stdin = strings.NewReader(privateKey)
		out, err := cmd.CombinedOutput()
		if err == nil && len(strings.TrimSpace(string(out))) > 0 {
			return strings.TrimSpace(string(out))
		}
	}
	// mock fallback
	h := []byte(privateKey)
	out := make([]byte, 32)
	for i := 0; i < 32; i++ {
		out[i] = h[i%len(h)] ^ byte(i)
	}
	return base64.StdEncoding.EncodeToString(out)
}
