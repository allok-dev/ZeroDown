package zerodown

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

type signalHandler struct {
	listener net.Listener
	server   *http.Server
	config   *ServerConfig
	tm       *taskManager
}

func newSignalHandler(listener net.Listener, server *http.Server, config *ServerConfig, tm *taskManager) *signalHandler {
	return &signalHandler{
		listener: listener,
		server:   server,
		config:   config,
		tm:       tm,
	}
}

func (sh *signalHandler) reload() error {
	tl, ok := sh.listener.(*net.TCPListener)
	if !ok {
		return errors.New("listener is not tcp listener")
	}

	f, err := tl.File()
	if err != nil {
		return err
	}
	defer f.Close()

	os.Setenv(sh.config.GracefulEnv, "true")
	cmd := exec.Command(os.Args[0])
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = []*os.File{f}

	return cmd.Start()
}

func (sh *signalHandler) shutdown(ctx context.Context) error {
	sh.server.SetKeepAlivesEnabled(false)
	return sh.server.Shutdown(ctx)
}

func (sh *signalHandler) handleSignals() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
	defer signal.Stop(quit)

	for {
		select {
		case sig := <-quit:
			switch sig {
			case syscall.SIGHUP:
				continue
			case syscall.SIGINT:
				if err := sh.tm.executeShutdownHooks(); err != nil {
					return err
				}
				if err := sh.reload(); err != nil {
					return err
				}
			case syscall.SIGQUIT:
				if err := sh.tm.executeShutdownHooks(); err != nil {
					return err
				}
				shutdownCtx, cancel := context.WithTimeout(context.Background(), sh.config.ShutdownTime)
				defer cancel()
				return sh.shutdown(shutdownCtx)
			}
		}
	}
}
