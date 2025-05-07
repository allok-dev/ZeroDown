# ZeroDown
Go语言Web服务平滑重启

## 介绍
ZeroDown是一个用于Go语言Web服务平滑重启的库。它提供了一种简单的方法来在不中断服务的情况下进行升级和重启。
最初是参考`github.com/fvbock/endless`实现的，代码托管于`github.com/with-zeal/restart`，但是由于原账号`with-zeal`丢失二次验证App的动态码，且原来的代码存在一些问题，例如：变量名不明确，且缺少一些必要的注释，所以重新实现了一个。

## 特性
- 平滑重启：在不中断服务的情况下进行升级和重启。
- 优雅关闭：在关闭前等待一段时间，以确保所有请求都已完成处理。
- 启动钩子：在服务器启动时执行自定义的钩子函数。

## 安装
```shell
go get github.com/allok-dev/zerodown
```

## 用法
```go
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
```
同时支持`gin`等基于`net/http`的框架。