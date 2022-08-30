package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/src/calc"
	"github.com/gorilla/src/gee"
	"rsc.io/quote"
)

type Person interface {
	getName() string
}

type Student struct {
	name string
	age  int
}

// func (stu *Student) getName() string {
// 	return stu.name
// }

type Worker struct {
	name   string
	gender string
}

func (w *Worker) getName() string {
	return w.name
}

func ex(p Person) {
	fmt.Println(p.getName())
}
func ex2(s *Student) {
	fmt.Println(s.name)
}

func main() {
	//接口测试
	var _ Person = (*Worker)(nil)
	stu := &Student{
		name: "小明",
		age: 14,
	}
	ex2(stu)
	worker := &Worker{
		name: "小李",
		gender: "男",
	}
	//调用不同包路径下的函数
	fmt.Printf("calc.Add(1, 3): %v\n", calc.Add(1, 3))
	//调用第三方包
	fmt.Printf("quote.Hello(): %v\n", quote.Hello())
	ex(worker)

	//启动服务
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	r.GET("/hello", func(c *gee.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}