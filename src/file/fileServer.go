package file

import (
	"example/src/WebSocketHandler"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bep/debounce"
	"github.com/fsnotify/fsnotify"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func cors(prefix string, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, prefix)
		rp := strings.TrimPrefix(r.URL.RawPath, prefix)
		if len(p) < len(r.URL.Path) && (r.URL.RawPath == "" || len(rp) < len(r.URL.RawPath)) {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = p
			r2.URL.RawPath = rp
			w.Header().Set("Access-Control-Allow-Origin", "*")                                                            // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token") //header的类型
			w.Header().Add("Access-Control-Allow-Credentials", "true")                                                    //设置为true，允许ajax异步请求带cookie信息
			w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")                             //允许请求方法
			// w.Header().Set("content-type", "application/json;charset=UTF-8")                                              //返回数据格式是json
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			h.ServeHTTP(w, r2)
		} else {
			http.NotFound(w, r)
		}
	}
}

func FileServer() {
	str, _ := os.Getwd()
	file := http.FileServer(http.Dir(str + "/video"))
	http.Handle("/", cors("/", file))
	err := http.ListenAndServe(":8181", nil)
	if err != nil {
		log.Println(err)
	}
}

// 添加监听文件
func FsnotifyWatch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}
	defer watcher.Close()

	done := make(chan bool)

	debounced := debounce.New(1000 * time.Millisecond)

	go func() {
		defer close(done)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Println("event.Op.String()", event.Name)
				if event.Op.String() == "CREATE" || event.Op.String() == "WRITE" || event.Op.String() == "Chmod" {
					// err := ffmpeg.Input("./"+event.Name).
					// 	Output("./out1.mp4", ffmpeg.KwArgs{"c:v": "libx264"}).
					// 	OverWriteOutput().ErrorToStdOut().Run()
					// log.Fatal("ffmpeg failed: ", err)
					f := func() {
						// str := event.Name[1:len(event.Name)]
						fmt.Println("test:", event.Name)
						err := ffmpeg.Input(event.Name).
							Output("./video/sample_data/test.mp4", ffmpeg.KwArgs{"c:v": "h264", "preset": "fast"}).
							OverWriteOutput().ErrorToStdOut().Run()
						fmt.Println(err)
						WebSocketHandler.BroadcastUsers(1, event.Name)
					}
					debounced(f)

				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	err = watcher.Add("./video")
	if err != nil {
		log.Fatal("Add failed:", err)
	}
	<-done
}
