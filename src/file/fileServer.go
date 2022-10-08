package file

import (
	"example/src/WebSocketHandler"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fsnotify/fsnotify"
)

func FileServer() {
	str, _ := os.Getwd()
	file := http.FileServer(http.Dir(str))
	fmt.Printf("当前目录：%v\n", file)
	http.Handle("/", http.StripPrefix("/", file))
	fmt.Printf("http://127.0.0.1:8181\n")
	err := http.ListenAndServe(":8181", nil)
	if err != nil {
		log.Println(err)
	}
}

func FsnotifyWatch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		defer close(done)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				//是否等于创建文件
				if event.Op.String() == "CREATE" {
					WebSocketHandler.BroadcastUsers(1, event.Name)
					// log.Printf("%s %s\n", event.Name, event.Op)
				}
				log.Printf("%s %s\n", event.Name, event.Op)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	//监听目录
	err = watcher.Add("./")
	if err != nil {
		log.Fatal("Add failed:", err)
	}
	<-done
}
