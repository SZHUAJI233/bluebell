package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web/controller"
	"web/router"
	"web/setting"

	"go.uber.org/zap"
)

func main() {

	// 初始化
	controller.Init()
	defer controller.Close()

	// 注册路由
	r := router.Setup(setting.Conf.Mode)

	// 启动服务
	srv := &http.Server{ // 将r注册到server的handler中
		Addr:    ":8080",
		Handler: r,
	}

	// 开启goruntine去提供web服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Error("listen:%s\n", zap.Error(err))
		}
	}()

	// 优雅关机
	quit := make(chan os.Signal, 1)                      // 创建一个chan通道用于接收信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 监听这两种关闭信号，发生则传输到quit中
	<-quit                                               // 从quit中读取数据，quit为空，则会阻塞
	zap.L().Info("Shutdown server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("Server down: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")

}
