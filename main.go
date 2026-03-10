package main

import (
	"context"
	//"google.golang.org/appengine/log"
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
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("starting...")

	if err := util.GetDataViper(); err != nil {
		//fmt.Println("GetDataViper 读取失败", err)
		log.Println("getDataviper have err", err)
		return

	}

	if err := config.Pgsql(); err != nil {
		log.Println("pgsql have err", err)
		return
	}
	// if err := config.Redis(); err != nil {
	// 	log.Println("redis have err", err)
	// 	return
	// }

	defer func() {
		config.Close()
	}()

	r := router.Router()
	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    ":8889",
		Handler: r,
	}

	// 创建一个信号监听器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 启动 HTTP 服务器
	go func() {
		log.Println("HTTP server listening on :8889")
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
