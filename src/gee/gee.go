package gee

import (
	"fmt"
	"net/http"
)

//定义一个函数类型
type HandlerFunc func (http.ResponseWriter,*http.Request) 

//顶一个类型 Engine 
//属性 router 是一个 map key 为 string value 是 一个函数类型
type Engine struct{
	router map[string]HandlerFunc
}

func New() *Engine{
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) addRoute(method string,pattern string,handlerFunc HandlerFunc){
	key := method + "-" +pattern
	engine.router[key] = handlerFunc
}

func (engine *Engine) GET(palette string,handlerFunc HandlerFunc){	
	engine.addRoute("GET",palette,handlerFunc)
}

func (engine *Engine) POST(palette string,handlerFunc HandlerFunc){
	engine.addRoute("POST",palette,handlerFunc)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}