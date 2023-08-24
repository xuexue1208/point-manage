package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"point-manage/config"
	"point-manage/db"
	"point-manage/middle"
	"point-manage/pprof"
	"point-manage/router"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"go.uber.org/zap"
)

func main() {

	//加载配置文件
	config.Initialize()

	//线上关闭gin的bebug
	ginMode := config.Instance.EnvParams.Mode
	//ginMode := os.Getenv("gin_mode")
	if ginMode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	//初始化数据库
	db.NewInit()
	// //插入excel数据
	//utils.JsonFile.ReadJsonFile("./3.0.40.json")
	//time.Sleep(100000000000)

	//初始化路由配置
	r := gin.New()
	r.Use(
		//跨域配置
		middle.Cors(),
		middle.Recovery(),
		//prometheus
		middle.PromMiddleware(&middle.PromOpts{ExcludeRegexEndpoint: "^/metrics|/favicon.ico"}),
	)
	//初始化路由
	router.Router.InitApiRouter(r)
	//gin server启动
	srv := &http.Server{
		Addr:    config.Instance.EnvParams.ListenAddr,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(
				"监听失败",                         // 输出自定义错误提示
				zap.String("listen", srv.Addr), // 有关该错误的关键信息
				zap.Error(err),
			)
		}
	}()
	logger.Info("gin server 启动成功", srv.Addr)
	//初始化pprof
	pprof.InitPprof()

	//等待中断信号，优雅关闭所有server及DB
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	//设置ctx超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//cancel用于释放ctx
	defer cancel()

	//关闭gin server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Gin Server关闭异常:", zap.Error(err))
	}
	logger.Info("Gin Server退出成功")
}
