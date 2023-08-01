package main

import (
	"fmt"
	"net/http"

	"github.com/aeon27/myblog/pkg/setting"
	router "github.com/aeon27/myblog/routers"
)

func main() {
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router.InitRouter(),
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
