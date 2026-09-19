package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// GeoIPInfo — IP 归属地信息
type GeoIPInfo struct {
	IP          string `json:"ip"`
	City        string `json:"city"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Region      string `json:"region"`
}

type GeoIPService struct {
	rdb         *redis.Client
	httpClient  *http.Client
	cacheTTL    time.Duration
}

func NewGeoIPService(rdb *redis.Client) *GeoIPService {
	return &GeoIPService{
		rdb:        rdb,
		httpClient: &http.Client{Timeout: 5 * time.Second},
		cacheTTL:   24 * time.Hour,
	}
}

// Lookup — 根据 IP 查询归属地，结果缓存 24 小时
func (s *GeoIPService) Lookup(ctx context.Context, ip string) *GeoIPInfo {
	if ip == "" || isPrivateIP(ip) {
		return &GeoIPInfo{IP: ip, City: "Local Network", Country: "", CountryCode: ""}
	}

	// 先查 Redis 缓存
	if s.rdb != nil {
		key := fmt.Sprintf("geoip:%s", ip)
		if cached, err := s.rdb.Get(ctx, key).Bytes(); err == nil {
			var info GeoIPInfo
			if json.Unmarshal(cached, &info) == nil {
				return &info
			}
		}
	}

	// 调用 ip-api.com 免费 API（无需 key，限 45 req/min）
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country,countryCode,regionName,city,query", ip)
	resp, err := s.httpClient.Get(url)
	if err != nil {
		log.Warn().Err(err).Str("ip", ip).Msg("geoip lookup failed")
		return &GeoIPInfo{IP: ip}
	}
	defer resp.Body.Close()

	var apiResp struct {
		Status      string `json:"status"`
		Country     string `json:"country"`
		CountryCode string `json:"countryCode"`
		RegionName  string `json:"regionName"`
		City        string `json:"city"`
		Query       string `json:"query"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		log.Warn().Err(err).Str("ip", ip).Msg("geoip decode failed")
		return &GeoIPInfo{IP: ip}
	}

	if apiResp.Status != "success" {
		return &GeoIPInfo{IP: ip}
	}

	info := &GeoIPInfo{
		IP:          ip,
		City:        apiResp.City,
		Country:     apiResp.Country,
		CountryCode: apiResp.CountryCode,
		Region:      apiResp.RegionName,
	}

	// 写入 Redis 缓存
	if s.rdb != nil {
		data, _ := json.Marshal(info)
		s.rdb.Set(ctx, fmt.Sprintf("geoip:%s", ip), data, s.cacheTTL)
	}

	return info
}

// isPrivateIP — 判断是否为内网/私有 IP
func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return true
	}
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() {
		return true
	}
	return false
}
