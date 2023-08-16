package main

import (
	"fmt"
	"log"
	"syscall"

	"github.com/aeon27/myblog/models"
	"github.com/aeon27/myblog/pkg/logging"
	"github.com/aeon27/myblog/pkg/setting"
	router "github.com/aeon27/myblog/routers"
	"github.com/fvbock/endless"
)

func main() {
	setting.Setup()
	models.Setup()
	logging.Setup()

	// 借助 endless 实现服务优雅启停
	// endless server 监听以下几种信号量
	// syscall.SIGHUP: 触发 fork 子进程和重新启动
	// syscall.SIGUSR1/syscall.SIGTSTP: 被监听，但不会触发任何动作
	// syscall.SIGUSR2: 触发hammerTime
	// syscall.SIGINT/sycall.SIGTERM: 触发服务器关闭（优雅关闭）
	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20 // 2^20 = 1M
	addr := fmt.Sprintf(":%d", setting.ServerSetting.HTTPPort)

	server := endless.NewServer(addr, router.InitRouter())
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid: %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}
