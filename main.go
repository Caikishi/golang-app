package main

import (
	"example/src/WebSocketHandler"
	"example/src/file"
	"example/src/gee"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/gorilla/websocket"
)

// type Person interface {
// 	getName() string
// }

// type Student struct {
// 	name string
// 	age  int
// }

// // func (stu *Student) getName() string {
// // 	return stu.name
// // }

// type Worker struct {
// 	name   string
// 	gender string
// }

// func (w *Worker) getName() string {
// 	fmt.Println("我是 worker 我处理一些 worker 在 getName 的逻辑")
// 	return w.name
// }

// func ex(p Person) {
// 	fmt.Println("我是 ex 我处理一些 ex 的逻辑 然后调用了接口 getName")
// 	fmt.Println(p.getName())

// }
// func ex2(s *Student) {
// 	fmt.Println(s.name)
// }

// type id int64

// func Sum[T id](args ...T) T {
// 	var sum T
// 	for i := 0; i < len(args); i++ {
// 		sum += args[i]
// 	}
// 	return sum
// }

// type QueryWrapper struct {
// }

// func (query *QueryWrapper) eq(s string, q string) *QueryWrapper {
// 	fmt.Printf("添加条件查询条件,%v == %v\n", s, q)
// 	return query
// }

// // 查询
// func (q *QueryWrapper) query() {
// 	fmt.Println("执行sql查询")
// }

// func test(args ...any) { // 可以接受任意个string参数
// 	for _, v := range args {
// 		log.Println(v)
// 	}
// }

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func main() {
	r := gee.NewService()
	r.GET("/hello", func(c *gee.Context) {
		// expect /hello?name=geektutu
		c.HTML(http.StatusOK, "<!DOCTYPE html><html><head><meta  charset='UTF-8' /><title>菜鸟教程(runoob.com)</title><script src='https://cdnjs.cloudflare.com/ajax/libs/lodash.js/4.17.21/lodash.min.js'></script><style>html,body{height:100%;margin:0;background-color: #000;color: #fff;padding:0}#sse{height:100%}#iframe{position:absolute;width:100%;height:100%;top:0;left:0}</style><script type='text/javascript'>var host='http://192.168.62.98:8181/';var change=function(data){document.getElementById('iframe').src=host+data};function strToByte(str){var bytes=new Array();var len,c;len=str.length;for(var i=0;i<len;i++){c=str.charCodeAt(i);if(c>=0x010000&&c<=0x10ffff){bytes.push(((c>>18)&0x07)|0xf0);bytes.push(((c>>12)&0x3f)|0x80);bytes.push(((c>>6)&0x3f)|0x80);bytes.push((c&0x3f)|0x80)}else if(c>=0x000800&&c<=0x00ffff){bytes.push(((c>>12)&0x0f)|0xe0);bytes.push(((c>>6)&0x3f)|0x80);bytes.push((c&0x3f)|0x80)}else if(c>=0x000080&&c<=0x0007ff){bytes.push(((c>>6)&0x1f)|0xc0);bytes.push((c&0x3f)|0x80)}else{bytes.push(c&0xff)}}return bytes}function WebSocketTest(){if('WebSocket'in window){var ws=new WebSocket('ws://192.168.62.98:8080/ws');ws.onopen=function(){ws.send('testing')};ws.onmessage=(evt)=>{let receivedMsg=evt.data;console.log('接收数据：'+receivedMsg);let d=document.getElementById('test');d.innerHTML=d.innerHTML+`<div>${receivedMsg}</div>`};ws.onclose=function(){console.log('连接已关闭...')}}else{console.log('您的浏览器不支持 WebSocket!')}}WebSocketTest();</script></head><body><div id='test'></div><div id='sse'><!--<span>最新文件:<span id='t'></span></span>--><!--<div style='height: 100%'><iframe id='iframe'src=''frameborder='0'height='100%'width='100%'seamless></iframe></div>--><!--<a href='javascript:WebSocketTest()'>运行WebSocket</a>--></div></body></html>")

		// c.HTML(http.StatusOK, "<!DOCTYPE html><html><head><meta charset='utf-8'/><title>测试</title><script src='https://cdnjs.cloudflare.com/ajax/libs/lodash.js/4.17.21/lodash.min.js'></script><style>html,body{height:100%;margin:0;padding:0}#sse{height:100%}#iframe{position:absolute;width:100%;height:100%;top:0;left:0}</style><script type='text/javascript'>var host='http://127.0.0.1:8181/';var change=function(data){document.getElementById('iframe').src=host+data};function WebSocketTest(){if('WebSocket'in window){var ws=new WebSocket('ws://127.0.0.1:8080/ws');ws.onopen=function(){ws.send('testing')};ws.onmessage=_.debounce((evt)=>{let receivedMsg=evt.data;console.log('接收数据：'+receivedMsg);change(receivedMsg)},300);ws.onclose=function(){console.log('连接已关闭...')}}else{console.log('您的浏览器不支持 WebSocket!')}}WebSocketTest();</script></head><body><div id='sse'><!--<span>最新文件:<span id='t'></span></span>--><div style='height: 100%'><iframe id='iframe'src=''frameborder='0'height='100%'width='100%'seamless></iframe></div><!--<a href='javascript:WebSocketTest()'>运行WebSocket</a>--></div></body></html>")
	})
	go r.Run(":9998")
	go file.FsnotifyWatch()
	go WebSocketHandler.StartWebsocket("0.0.0.0:8899")
	// open("http://127.0.0.1:9998/hello")
	service := gee.NewService()

	// service.POST("/pushReact", func(ctx *gee.Context) {
	// 	go automatic.BuildReact(ctx)
	// })

	// service.POST("/pushBs", func(ctx *gee.Context) {
	// 	go automatic.BuildJava(ctx)
	// })
	go service.Run(":9999")
	file.FileServer()
	// proxy.Start()
}
