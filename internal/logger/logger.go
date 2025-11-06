// logger 包提供日志系统初始化功能
// 使用zap作为日志库，支持文件和控制台双输出
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger 初始化日志系统
// 功能：
// 1. 创建日志目录
// 2. 配置文件日志输出（使用lumberjack实现日志轮转）
// 3. 配置控制台日志输出
// 4. 设置日志级别（文件：Info级别，控制台：Debug级别）
// 5. 使用JSON格式输出到文件，控制台使用可读格式
func InitLogger() {
	// 创建logs目录（如果不存在）
	_ = os.MkdirAll("logs", 0755)

	// 配置lumberjack日志轮转器
	// Filename: 日志文件路径
	// MaxSize: 单个日志文件最大大小（MB），超过后自动轮转
	// MaxBackups: 保留的旧日志文件数量
	// MaxAge: 日志文件保留天数
	lj := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    50, // 50MB
		MaxBackups: 3,  // 保留3个备份文件
		MaxAge:     28, // 保留28天
	}

	// 创建文件日志写入器
	writeSyncer := zapcore.AddSync(lj)

	// 创建控制台日志写入器
	consoleSyncer := zapcore.AddSync(os.Stdout)

	// 配置日志编码器
	encoderCfg := zap.NewProductionEncoderConfig()
	// 设置时间字段名为"ts"
	encoderCfg.TimeKey = "ts"
	// 使用ISO8601时间格式
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// 创建双核心日志系统
	// 文件输出：JSON格式，Info级别及以上
	// 控制台输出：可读格式，Debug级别及以上
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), writeSyncer, zap.InfoLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), consoleSyncer, zap.DebugLevel),
	)

	// 创建logger实例
	logger := zap.New(core)

	// 替换全局logger，使zap.L()可以获取到配置好的logger
	zap.ReplaceGlobals(logger)
}
