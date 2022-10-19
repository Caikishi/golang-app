package file

import (
	"example/src/WebSocketHandler"
	"log"
	"net/http"
	"os"

	"github.com/fsnotify/fsnotify"
)

func FileServer() {
	str, _ := os.Getwd()
	file := http.FileServer(http.Dir(str))
	http.Handle("/", http.StripPrefix("/", file))
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
				// if event.Op.String() == "CREATE" {

				// 	// log.Printf("%s %s\n", event.Name, event.Op)
				// }
				WebSocketHandler.BroadcastUsers(1, event.Name)
				log.Printf("%s %s\n", event.Name, event.Op)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	err = watcher.Add("./")
	if err != nil {
		log.Fatal("Add failed:", err)
	}
	<-done
}
