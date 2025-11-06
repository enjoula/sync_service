// metrics 包提供Prometheus监控指标功能
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	// HTTP请求总数计数器
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)
	
	// HTTP请求耗时直方图
	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

// InitMetrics 初始化Prometheus指标
// 注册所有自定义指标
func InitMetrics() {
	// 注册自定义指标
	prometheus.MustRegister(HTTPRequestsTotal)
	prometheus.MustRegister(HTTPRequestDuration)
}

// Handler 返回Prometheus指标HTTP处理器
// 用于暴露/metrics端点
func Handler() http.Handler {
	return promhttp.Handler()
}

