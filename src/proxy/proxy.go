package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var replaceStr = "/api_proxy/"

func proxy2(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse("http://172.16.26.61:86")
	r.URL.Path = strings.Replace(r.URL.Path, "/api_proxy", "", -1)
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func proxy(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse("http://172.16.26.61:666")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func Start() {
	http.HandleFunc(replaceStr, proxy2)
	http.HandleFunc("/", proxy)

	err := http.ListenAndServe(":9095", nil)
	if err != nil {
		fmt.Println("HTTP server failed,err:", err)
		return
	}
}
