package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var clientHttps *http.Client
var clientHttp *http.Client

func init() {
	initConf()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	clientHttps = &http.Client{Transport: tr}
	clientHttp = &http.Client{}
}

// handler entrance
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

	// from outer
	if uri.Source != "" {
		handleSource(w, uri)
		return
	}

	// from json
	for _, entry := range uri.Headers {
		for key, value := range entry {
			w.Header().Set(key, value)
		}
	}

	fmt.Fprintf(w, uri.Body)
}

func handleSource(w http.ResponseWriter, uri Uri) {
	if strings.HasPrefix(uri.Source, "http://") {
		handleHttp(w, clientHttp, uri)
		return
	}

	if strings.HasPrefix(uri.Source, "https://") {
		handleHttp(w, clientHttps, uri)
		return
	}

	// handle local file
	fmt.Fprintf(w, loadCert(uri.Body))
	return
}

// handle http and https
func handleHttp(w http.ResponseWriter, client *http.Client, uri Uri) {
	resp, err := client.Get(uri.Source)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	for key, entry := range resp.Header {
		w.Header().Set(key, strings.Join(entry, ","))
	}
	fmt.Fprintf(w, string(body))
}

func main() {
	for i := 0; i < len(appConf.Uris); i++ {
		http.HandleFunc(appConf.Uris[i].Uri, handler)
	}
	//http.HandleFunc("/", handler)
	//http.ListenAndServeTLS(":20175", "cert/server.cert", "cert/server.key", nil)
	http.ListenAndServe(":20175", nil)
}

func loadCert(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	return string(data)
}
