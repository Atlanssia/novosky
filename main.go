package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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
	if strings.HasSuffix(uri.Body, "cert") {
		fmt.Fprintf(w, loadCert(uri.Body))
		return
	}
	fmt.Fprintf(w, uri.Body)
}

func main() {
	for i := 0; i < len(appConf.Uris); i++ {
		http.HandleFunc(appConf.Uris[i].Uri, handler)
	}
	//    http.HandleFunc("/", handler)
	http.ListenAndServeTLS(":20175", "cert/server.cert", "cert/server.key", nil)
}

func loadCert(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	return string(data)
}
