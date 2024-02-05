package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"stock/biz/dal/sql"
	"stock/biz/router"
	"stock/util"
	"syscall"
	"time"
)

func main() {

	if err := util.GetDataViper(); err != nil {
		fmt.Println("GetDataViper 读取失败", err)
		return
	}

	if err := config.Mysql(); err != nil {
		fmt.Println("mysql连接失败", err)
		return
	}
	if err := config.Redis(); err != nil {
		fmt.Println("newRedis连接失败", err)
		return
	}

	defer func() {
		config.Close()
	}()

	r := router.Router()
	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// 创建一个信号监听器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 启动 HTTP 服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	// 创建一个5秒的超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅地关闭 HTTP 服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped gracefully")
}
