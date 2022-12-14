package WebSocketHandler

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/process"
)

const (
	// 允许等待的写入时间
	writeWait = 100 * time.Second

	// 允许从对等方读取下一个pong消息的时间。
	pongWait = 600 * time.Second

	// 将ping发送到此时间段的对等方。必须小于pongWait。
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// 最大的连接ID，每次连接都加1 处理
var maxConnId int64

// 客户端读写消息
type wsMessage struct {
	// websocket.TextMessage 消息类型
	messageType int
	data        []byte
}

// ws 的所有连接
// 用于广播
var wsConnAll map[int64]*wsConnection

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有的CORS 跨域请求，正式环境可以关闭
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 客户端连接
type wsConnection struct {
	wsSocket *websocket.Conn // 底层websocket
	inChan   chan *wsMessage // 读队列
	outChan  chan *wsMessage // 写队列

	mutex     sync.Mutex // 避免重复关闭管道,加锁处理
	isClosed  bool
	closeChan chan byte // 关闭通知
	id        int64
}

func wsHandler(resp http.ResponseWriter, req *http.Request) {
	// 应答客户端告知升级连接为websocket
	wsSocket, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Println("升级为websocket失败", err.Error())
		return
	}
	maxConnId++
	// TODO 如果要控制连接数可以计算，wsConnAll长度
	// 连接数保持一定数量，超过的部分不提供服务
	wsConn := &wsConnection{
		wsSocket:  wsSocket,
		inChan:    make(chan *wsMessage, 1000),
		outChan:   make(chan *wsMessage, 1000),
		closeChan: make(chan byte),
		isClosed:  false,
		id:        maxConnId,
	}
	wsConnAll[maxConnId] = wsConn

	pwd, _ := os.Getwd()
	//获取文件或目录相关信息
	fileInfoList, err := ioutil.ReadDir(pwd)
	if err != nil {
		log.Fatal(err)
	}
	// 排序
	sort.SliceStable(fileInfoList, func(i, j int) bool {
		return fileInfoList[i].ModTime().Unix() > fileInfoList[j].ModTime().Unix()
	})

	// addrs, err := net.InterfaceAddrs()
	// for _, address := range addrs {
	// 	// 检查ip地址判断是否回环地址
	// 	if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
	// 		if ipnet.IP.To4() != nil {
	// 			fmt.Println(ipnet.IP.String())
	// 		}
	// 	}
	// }
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// log.Println("当前在线人数", len(wsConnAll))

	// 处理器,发送定时信息，避免意外关闭
	go wsConn.processLoop()
	// 读协程
	go wsConn.wsReadLoop()
	// 写协程
	go wsConn.wsWriteLoop()

	pNames := ProcessName(int32(os.Getpid()))
	var pName string
	for _, v := range pNames {
		pName = v
	}
	for i := range fileInfoList {
		//去除文件夹
		if !fileInfoList[i].IsDir() && fileInfoList[i].Name() != pName {
			// fmt.Printf("文件名：%v,文件修改时间:%v,文件mode%v\n", fileInfoList[i].Name(), fileInfoList[i].ModTime(), addrs) //打印当前文件或目录下的文件或目录名
			go BroadcastUsers(1, fileInfoList[i].Name())
			return
		}

	}

}

/*
获取所有进程名，以数组返回
*/

func ProcessName(s int32) (pname []string) {
	pids, _ := process.Pids()
	for _, pid := range pids {
		if s == pid {
			pn, _ := process.NewProcess(pid)
			pName, _ := pn.Name()
			return append(pname, pName)
		}
	}
	return pname
}

// 处理队列中的消息
func (wsConn *wsConnection) processLoop() {
	// 处理消息队列中的消息
	// 获取到消息队列中的消息，处理完成后，发送消息给客户端
	for {
		msg, err := wsConn.wsRead()
		if err != nil {
			// log.Println("获取消息出现错误", err.Error())
			break
		}
		log.Println("接收到消息", string(msg.data), msg.messageType)
		// 修改以下内容把客户端传递的消息传递给处理程序
		// err = wsConn.wsWrite(msg.messageType, []byte("新建文件。。。"))
		if err != nil {
			log.Println("发送消息给客户端出现错误", err.Error())
			break
		}
	}
}

// 处理消息队列中的消息
func (wsConn *wsConnection) wsReadLoop() {
	// 设置消息的最大长度
	wsConn.wsSocket.SetReadLimit(maxMessageSize)
	wsConn.wsSocket.SetReadDeadline(time.Now().Add(pongWait))
	for {
		// 读一个message
		msgType, data, err := wsConn.wsSocket.ReadMessage()
		if err != nil {
			websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure)
			log.Println("消息读取出现错误", err.Error())
			wsConn.close()
			return
		}
		req := &wsMessage{
			msgType,
			data,
		}
		// 放入请求队列,消息入栈
		select {
		case wsConn.inChan <- req:
		case <-wsConn.closeChan:
			return
		}
	}
}

// 发送消息给客户端
func (wsConn *wsConnection) wsWriteLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		// 取一个应答
		case msg := <-wsConn.outChan:
			// 写给websocket
			if err := wsConn.wsSocket.WriteMessage(msg.messageType, msg.data); err != nil {
				// log.Println("发送消息给客户端发生错误", err.Error())
				// 切断服务
				wsConn.close()
				return
			}
		case <-wsConn.closeChan:
			// 获取到关闭通知
			return
		case <-ticker.C:
			// 出现超时情况
			wsConn.wsSocket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := wsConn.wsSocket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// 写入消息到队列中
func (wsConn *wsConnection) wsWrite(messageType int, data []byte) error {
	select {
	case wsConn.outChan <- &wsMessage{messageType, data}:
	case <-wsConn.closeChan:
		return errors.New("连接已经关闭")
	}
	return nil
}

// 读取消息队列中的消息
func (wsConn *wsConnection) wsRead() (*wsMessage, error) {
	select {
	case msg := <-wsConn.inChan:
		// 获取到消息队列中的消息
		return msg, nil
	case <-wsConn.closeChan:

	}
	return nil, errors.New("连接已经关闭")
}

// 关闭连接
func (wsConn *wsConnection) close() {
	log.Println("关闭连接被调用了")
	wsConn.wsSocket.Close()
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if wsConn.isClosed == false {
		wsConn.isClosed = true
		// 删除这个连接的变量
		delete(wsConnAll, wsConn.id)
		close(wsConn.closeChan)
	}
}

// 对所有用户进行广播
func BroadcastUsers(messageType int, data string) {
	for _, ws := range wsConnAll {
		ws.wsWrite(messageType, []byte(data))
	}
}

// 对所有用户进行广播
func BroadcastUsers2(messageType int, data []byte) {
	for _, ws := range wsConnAll {
		ws.wsWrite(messageType, data)
	}
}

// 启动程序
func StartWebsocket(addrPort string) {
	wsConnAll = make(map[int64]*wsConnection)
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(addrPort, nil)
}
