package zerodown

import "time"

type ServerConfig struct {
	Addr         string
	CertFile     string
	KeyFile      string
	EnableTLS    bool
	ShutdownTime time.Duration // 优雅关闭超时时间
	GracefulEnv  string        // 优雅重启环境变量名
}

func NewDefaultConfig() *ServerConfig {
	return &ServerConfig{
		Addr:         ":8080",
		EnableTLS:    false,
		ShutdownTime: time.Minute,
		GracefulEnv:  "graceful",
	}
}
