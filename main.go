package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/src/gee"
)

func main() {
	r := gee.New()

	r.GET("/",func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w,"URL.Path = %q\n",r.URL.Path)
	})

	go r.GET("/getUser",func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})
	r.Run(":9998")
}

// func sendMsg(ch chan int) {
// 	fmt.Println("等待 5 秒后, 往 ch 发送 10 <- 10")
// 	time.Sleep(time.Second * 5)
//     ch <- 10
// }

// func recv(ch chan int){
// 	v := <-ch
// 	fmt.Printf("v = %v\n",v)
// }

// func main() {
// 	//定义一个通道 ch
//     ch := make(chan int)
// 	ch2 := make(chan int)

// 	//暂停 3 秒
// 	// 开启异步 go 关键字 是把这个函数变为异步执行
//     go sendMsg(ch)
// 	fmt.Println("开始监听 ch")
// 	go recv(ch)
// 	// fmt.Println("开始监听 ch")
// 	for i := range ch2 {
// 		fmt.Printf("收到:%v\n",i)
// 	}

// }