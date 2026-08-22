package main

import (
	"Blog/pkg/logger"
	"Blog/pkg/snowflake"
	"Blog/repository/database"
	"Blog/repository/redis"
	"Blog/router"
	"Blog/settings"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := settings.Init(); err != nil {
		fmt.Printf("init viper error: %s ", err.Error())
		return
	}

	if err := snowflake.Init(); err != nil {
		fmt.Printf("init snowflake error: %s", err.Error())
		return
	}
	if err := database.Init(); err != nil {
		fmt.Printf("init mysql connect error: %s", err.Error())
		return
	}
	if err := redis.Init(); err != nil {
		fmt.Printf("init redis connect error: %s", err.Error())
		return
	}
	if err := logger.Init("development"); err != nil {
		fmt.Printf("init logger error: %s", err.Error())
	}

	//优雅关机
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", settings.Cnf.App.Port),
		Handler: router.Init(),
	}
	//用协程实现持续监听
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	//定义五秒超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//函数退出前调用cancel，释放资源
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown: ", err)
	}

}
