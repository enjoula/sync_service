// config 包提供配置管理功能
// 支持从本地YAML配置文件、环境变量和Etcd远程配置中心读取配置
package config

import (
	"context"
	"encoding/json"
	"time"

	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Cfg 是全局配置对象，使用Viper管理配置
var Cfg *viper.Viper

// EtcdCli 是Etcd客户端连接，用于从Etcd读取敏感配置信息
var EtcdCli *clientv3.Client

// ConfigKey 是Etcd中存储普通配置的键路径
var ConfigKey = "/video-service/config"

// SecretKey 是Etcd中存储敏感信息（如JWT密钥、数据库密码）的键路径
var SecretKey = "/video-service/secret"

// Secrets 存储从Etcd获取的敏感配置信息
var Secrets = &Secret{}

// Secret 定义敏感配置信息的结构
// JWTKey: JWT签名密钥
// MySQLDsn: MySQL数据库连接字符串（可选，如果配置了会覆盖配置文件中的值）
type Secret struct {
	JWTKey   string `json:"jwt_key"`
	MySQLDsn string `json:"mysql_dsn,omitempty"`
}

// InitConfig 初始化配置系统
// 1. 创建Viper配置对象
// 2. 读取本地YAML配置文件
// 3. 启用环境变量自动读取
// 4. 如果配置了Etcd地址，连接Etcd并读取敏感信息
func InitConfig() {
	// 创建新的Viper实例
	Cfg = viper.New()

	// 设置配置文件路径
	Cfg.SetConfigFile("configs/config.yaml")

	// 启用环境变量自动读取，环境变量会覆盖配置文件中的值
	Cfg.AutomaticEnv()

	// 读取配置文件（忽略错误，允许配置文件不存在）
	_ = Cfg.ReadInConfig()

	// 获取Etcd地址配置
	etcdAddr := Cfg.GetString("etcd.addr")

	// 如果未配置Etcd地址，直接返回
	if etcdAddr == "" {
		return
	}

	// 创建Etcd客户端连接
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdAddr},
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	EtcdCli = cli

	// 创建带超时的上下文，用于从Etcd获取配置
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 从Etcd获取敏感配置信息
	resp, err := EtcdCli.Get(ctx, SecretKey)
	if err == nil && len(resp.Kvs) > 0 {
		var s Secret
		// 解析JSON格式的敏感配置
		_ = json.Unmarshal(resp.Kvs[0].Value, &s)
		// 合并敏感配置到全局Secrets对象
		mergeSecrets(&s)
	}
}

// mergeSecrets 合并从Etcd获取的敏感配置到全局配置
// 如果Etcd中的配置不为空，则覆盖默认配置
func mergeSecrets(s *Secret) {
	// 如果Etcd中配置了JWT密钥，则使用Etcd中的值
	if s.JWTKey != "" {
		Secrets.JWTKey = s.JWTKey
	}
	// 如果Etcd中配置了MySQL连接字符串，则覆盖配置文件中的值
	if s.MySQLDsn != "" {
		Cfg.Set("mysql.dsn", s.MySQLDsn)
	}
}
