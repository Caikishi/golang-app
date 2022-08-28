package main

import (
	"fmt"
	"net"
)

type Server struct{
    Ip string
    Port int 
} 

func (s *Server) Handler(conn net.Conn){
    fmt.Println("建立连接成功")
}

//创建一个 server
func NewServer(ip string,port int) *Server{
    //创建对象
     server := &Server{
        Ip: ip,
        Port: port,
     }
     return server
}

func (s *Server) Start(){
    listen,err := net.Listen("tcp",fmt.Sprintf("%s:%d",s.Ip,s.Port))
    if err != nil {
        fmt.Println("net.error",err)
        return
    }
    //defer 会在最后执行
    defer listen.Close()

    for {
        conn,err := listen.Accept()
        if(err != nil){
            continue
        }
        go s.Handler(conn)
    }

}
