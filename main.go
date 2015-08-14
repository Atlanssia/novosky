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
            handleRequest(w, appConf.Uris[i].Delay, appConf.Uris[i].Body)
            return
        }
    }
}

func handleRequest(w http.ResponseWriter, delay int, body string) {
    time.Sleep(time.Duration(delay) * time.Millisecond)
    fmt.Fprintf(w, body)
}

func main() {
    for i := 0; i < len(appConf.Uris); i++ {
        http.HandleFunc(appConf.Uris[i].Uri, handler)
    }
//    http.HandleFunc("/", handler)
    http.ListenAndServe(":7456", nil)
}
