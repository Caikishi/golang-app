package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":9999", nil))
	fmt.Println("启动成功")
}


// func recv(ch chan int) {
// 	fmt.Println("往 ch 发送 10 <- 10")
//     ch <- 10
// }
// func main() {
// 	//定义一个通道 ch
//     ch := make(chan int)
// 	//暂停 3 秒
// 	// 开启异步 go 关键字 是把这个函数变为异步执行 
//     go recv(ch)

// 	time.Sleep(time.Second * 3)
// 		//	fmt.Println("开始监听 ch")
// 	for i := range ch {
// 		fmt.Println(i)
// 	}
    
// }