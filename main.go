package main

import (
	"fmt"
	"net/http"
	"time"
)

func init() {
	initConf()
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	for i := 0; i < len(appConf.Uris); i++ {
		if r.URL.Path == appConf.Uris[i].Uri {
			handleRequest(w, appConf.Uris[i].Delay, appConf.Uris[i])
			return
		}
	}
}

func handleRequest(w http.ResponseWriter, delay int, uri Uri) {
	time.Sleep(time.Duration(delay) * time.Millisecond)

	for _, entry := range uri.Headers {
		for key, value := range entry {
			w.Header().Set(key, value)
		}
	}
	fmt.Fprintf(w, uri.Body)
}

func main() {
	for i := 0; i < len(appConf.Uris); i++ {
		http.HandleFunc(appConf.Uris[i].Uri, handler)
	}
	//    http.HandleFunc("/", handler)
	http.ListenAndServe(":2456", nil)
}
