// auth 包提供认证相关的功能
// 包括JWT token生成、验证和密钥管理
package auth

import (
	"time"
	"video-service/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

// defaultKey 是JWT签名的默认密钥，生产环境应通过配置或Etcd获取
var defaultKey = []byte("change-me-default")

// Claims 定义JWT token的载荷结构
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GetJWTKey 获取JWT签名密钥
// 优先从Etcd配置中获取，如果不存在则使用默认密钥
func GetJWTKey() []byte {
	if config.Secrets != nil && config.Secrets.JWTKey != "" {
		return []byte(config.Secrets.JWTKey)
	}
	return defaultKey
}

// GenerateToken 生成JWT token
// 参数：
//   - username: 用户名
//   - expiration: token过期时间
//
// 返回：
//   - token字符串
//   - 错误信息
func GenerateToken(username string, expiration time.Time) (string, error) {
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(GetJWTKey())
}

// ParseToken 解析并验证JWT token
// 参数：
//   - tokenString: token字符串
//
// 返回：
//   - Claims对象
//   - 错误信息
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return GetJWTKey(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

