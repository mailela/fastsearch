package core

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"gitee.com/rachel_os/fastsearch/global"
	"gitee.com/rachel_os/fastsearch/negative"
	"gitee.com/rachel_os/fastsearch/searcher"
	"gitee.com/rachel_os/fastsearch/searcher/model"
	"gitee.com/rachel_os/fastsearch/searcher/words"
	"gitee.com/rachel_os/fastsearch/web/controller"
	"gitee.com/rachel_os/fastsearch/web/router"

	//_ "net/http/pprof"
	"os"
	"os/signal"

	//"runtime"
	"syscall"
	"time"
)

// NewContainer 创建一个容器
func NewContainer(tokenizer *words.Tokenizer, neg *negative.Engine) *searcher.Container {
	container := &searcher.Container{
		Dir:       global.CONFIG.Data,
		Debug:     global.CONFIG.Debug,
		AllowDrop: global.CONFIG.AllowDrop,
		Neg:       neg,
		Tokenizer: tokenizer,
		Shard:     global.CONFIG.Shard,
		Timeout:   global.CONFIG.Timeout,
		BufferNum: global.CONFIG.BufferNum,
	}
	if err := container.Init(); err != nil {
		panic(err)
	}

	return container
}

func NewTokenizer(dictionaryPath string) *words.Tokenizer {
	return words.NewTokenizer(dictionaryPath)
}

// Initialize 初始化
func Initialize() {

	//runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪
	//runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪

	//go func() { http.ListenAndServe("0.0.0.0:6060", nil) }()

	global.CONFIG = Parser()

	if !global.CONFIG.Debug {
		log.SetOutput(os.Stdout) //将记录器的输出设置为os.Stdout
	}

	defer func() {

		if r := recover(); r != nil {
			fmt.Printf("panic: %s\n", r)
		}
	}()

	//初始化分词器
	tokenizer := NewTokenizer(global.CONFIG.Dictionary)
	neg := negative.NewNegative(global.CONFIG.Negative_data, tokenizer)
	neg.AllKeys(&model.NegSearch{})
	global.Container = NewContainer(tokenizer, neg)

	// 初始化业务逻辑
	controller.NewServices()

	// 注册路由
	r := router.SetupRouter()
	// 启动服务
	srv := &http.Server{
		Addr:    global.CONFIG.Addr,
		Handler: r,
	}
	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("listen:", err)
			os.Exit(0)
		}
	}()

	// 优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
