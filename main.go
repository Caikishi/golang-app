package main

import (
	"example/src/automatic"
	"example/src/gee"
	"log"
)

func init() {
	log.Println("启动项目")
	// go automatic.BuildJava(nil)
}

func main() {
	service := gee.NewService()

	service.POST("/pushReact", func(ctx *gee.Context) {
		go automatic.BuildReact(ctx)
	})

	service.POST("/pushBs", func(ctx *gee.Context) {
		go automatic.BuildJava(ctx)
	})
	service.Run(":9999")
}
