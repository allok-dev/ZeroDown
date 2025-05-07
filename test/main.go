package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/allok-dev/zerodown"
)

func main() {
	pid := os.Getpid()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// 创建路由
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%d 收到请求: %s %s", pid, r.Method, r.URL.Path)
		time.Sleep(50 * time.Second)
		w.Write([]byte(fmt.Sprintf("Hello, I'm from %d!", pid)))
	})

	// 创建服务器配置
	config := &zerodown.ServerConfig{
		Addr:         ":8080",
		GracefulEnv:  "graceful",
		ShutdownTime: 60,
	}

	// 创建服务器
	server := zerodown.NewServer(mux, config)

	// 添加启动钩子
	server.AddStartupHook(func() error {
		log.Printf("%d 服务器启动钩子执行", pid)
		return nil
	}, false)

	// 添加关闭钩子
	server.AddShutdownHook(func() error {
		log.Printf("%d 服务器关闭钩子执行", pid)
		return nil
	})

	// 启动服务器
	log.Printf("启动服务器于 Port:%s PID:%d", config.Addr, pid)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
