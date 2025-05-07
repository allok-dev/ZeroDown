package main

import (
	"log"
	"net/http"

	"github.com/allok-dev/zerodown"
)

func main() {
	// 创建路由
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("收到请求: %s %s", r.Method, r.URL.Path)
		w.Write([]byte("Hello, World!"))
	})

	// 创建服务器配置
	config := &zerodown.ServerConfig{
		Addr:         ":8080",
		GracefulEnv:  "graceful",
		ShutdownTime: 10,
	}

	// 创建服务器
	server := zerodown.NewServer(mux, config)

	// 添加启动钩子
	server.AddStartupHook(func() error {
		log.Println("服务器启动钩子执行")
		return nil
	}, false)

	// 添加关闭钩子
	server.AddShutdownHook(func() error {
		log.Println("服务器关闭钩子执行")
		return nil
	})

	// 启动服务器
	log.Printf("启动服务器于 %s", config.Addr)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
