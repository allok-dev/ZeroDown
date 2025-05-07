package zerodown

import (
	"net"
	"net/http"
	"os"
	"syscall"
)

type zerodownServer struct {
	server   *http.Server
	listener net.Listener
	config   *ServerConfig
	*taskManager
}

func NewServer(router http.Handler, config *ServerConfig) *zerodownServer {
	if config == nil {
		config = NewDefaultConfig()
	}

	return &zerodownServer{
		server: &http.Server{
			Addr:    config.Addr,
			Handler: router,
		},
		config:      config,
		taskManager: newTaskManager(),
	}
}

func (ser *zerodownServer) initListener() error {
	var err error
	if os.Getenv(ser.config.GracefulEnv) == "true" {
		f := os.NewFile(3, "")
		ser.listener, err = net.FileListener(f)
		if err == nil {
			syscall.Kill(os.Getppid(), syscall.SIGQUIT)
		}
	} else {
		ser.listener, err = net.Listen("tcp", ser.server.Addr)
	}
	return err
}

// Run 启动 HTTP 服务器并处理系统信号
//
// 信号处理：
//   - SIGHUP: 忽略该信号
//   - SIGINT: 触发服务器重启，执行以下步骤：
//     1. 执行所有关闭钩子函数
//     2. 保持当前监听器状态
//     3. 启动新的进程接管服务
//   - SIGQUIT: 触发服务器优雅关闭，执行以下步骤：
//     1. 停止接收新的连接
//     2. 等待现有连接处理完成
//     3. 执行关闭钩子函数
//     4. 关闭服务器
//
// 返回值：
//   - 如果服务器启动失败或在运行过程中发生错误，返回对应的错误
//   - 如果是正常关闭，返回 nil
func (ser *zerodownServer) Run() error {
	if err := ser.initListener(); err != nil {
		return err
	}

	if err := ser.executeStartupHooks(); err != nil {
		return err
	}

	handler := newSignalHandler(ser.listener, ser.server, ser.config, ser.taskManager)
	errChan := make(chan error, 1)

	go func() {
		var err error
		if ser.config.EnableTLS {
			err = ser.server.ServeTLS(ser.listener, ser.config.CertFile, ser.config.KeyFile)
		} else {
			err = ser.server.Serve(ser.listener)
		}
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	go func() {
		errChan <- handler.handleSignals()
	}()

	return <-errChan
}

// AddStartupHook 添加启动钩子函数
//
// 参数：
//   - handler: 钩子函数，返回值为 error 类型
//   - async: 是否异步执行钩子函数
//
// 钩子函数在服务器启动前执行，用于执行一些操作。例如log.Printf("服务启动于 PID: %d", os.Getpid())
// 如果 async 为 true，则钩子函数将在单独的 goroutine 中执行，不会阻塞服务器启动。
// 如果 async 为 false，则钩子函数将在服务器启动时同步执行，可能会阻塞服务器启动。
func (ser *zerodownServer) AddStartupHook(handler func() error, async bool) {
	ser.addStartupHook(handler, async)
}

// AddShutdownHook 添加关闭钩子函数
//
// 参数：
//   - handler: 钩子函数，返回值为 error 类型
//
// 钩子函数在服务器关闭/重启前执行，用于执行一些关闭操作。例如db.Close()
// 钩子函数将在服务器关闭/重启前同步执行，会阻塞服务器关闭。
func (ser *zerodownServer) AddShutdownHook(handler func() error) {
	ser.addShutdownHook(handler, false)
}
