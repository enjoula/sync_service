// utils 包提供通用工具函数
package utils

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetRealIP 获取客户端真实IP地址
// 优先级顺序：
// 1. X-Real-IP (Nginx代理常用)
// 2. X-Forwarded-For 的第一个IP (标准代理头)
// 3. CF-Connecting-IP (Cloudflare)
// 4. True-Client-IP (Akamai, Cloudflare企业版)
// 5. RemoteAddr (直连情况)
func GetRealIP(c *gin.Context) string {
	// 1. 尝试从 X-Real-IP 获取
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		ip = strings.TrimSpace(ip)
		if isValidIP(ip) {
			return ip
		}
	}

	// 2. 尝试从 X-Forwarded-For 获取（可能包含多个IP，取第一个）
	if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
		// X-Forwarded-For 格式: client, proxy1, proxy2
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if isValidIP(ip) {
				return ip
			}
		}
	}

	// 3. 尝试从 CF-Connecting-IP 获取 (Cloudflare)
	if ip := c.GetHeader("CF-Connecting-IP"); ip != "" {
		ip = strings.TrimSpace(ip)
		if isValidIP(ip) {
			return ip
		}
	}

	// 4. 尝试从 True-Client-IP 获取
	if ip := c.GetHeader("True-Client-IP"); ip != "" {
		ip = strings.TrimSpace(ip)
		if isValidIP(ip) {
			return ip
		}
	}

	// 5. 使用 Gin 的 ClientIP 方法作为后备
	return c.ClientIP()
}

// isValidIP 验证IP地址是否有效
func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// IsPrivateIP 判断是否为内网IP
func IsPrivateIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// 私有IP地址段：
	// 10.0.0.0/8
	// 172.16.0.0/12
	// 192.168.0.0/16
	// 127.0.0.0/8 (localhost)
	privateBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"169.254.0.0/16", // link-local
		"::1/128",        // IPv6 localhost
		"fc00::/7",       // IPv6 unique local addr
		"fe80::/10",      // IPv6 link-local
	}

	for _, block := range privateBlocks {
		_, subnet, _ := net.ParseCIDR(block)
		if subnet != nil && subnet.Contains(parsedIP) {
			return true
		}
	}

	return false
}

// GetIPInfo 获取IP地址的详细信息
func GetIPInfo(c *gin.Context) map[string]string {
	realIP := GetRealIP(c)

	return map[string]string{
		"real_ip":          realIP,
		"is_private":       boolToString(IsPrivateIP(realIP)),
		"x_real_ip":        c.GetHeader("X-Real-IP"),
		"x_forwarded_for":  c.GetHeader("X-Forwarded-For"),
		"cf_connecting_ip": c.GetHeader("CF-Connecting-IP"),
		"remote_addr":      c.Request.RemoteAddr,
	}
}

// boolToString 将布尔值转换为字符串
func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
