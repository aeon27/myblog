package main

import (
	"log"

	"github.com/aeon27/myblog/models"
	"github.com/robfig/cron"
)

func main() {
	log.Println("cron func starting...")

	c := cron.New()
	// 每小时执行一次
	c.AddFunc("0 0 * * * *", func() {
		log.Println("Run models.CleanAlltags...")
		models.CleanAllTags()
	})

	c.AddFunc("0 0 * * * *", func() {
		log.Println("Run models.CleanAllArticles...")
		models.CleanAllArticles()
	})

	c.Run()
}
