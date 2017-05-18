package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
			handleRequest(w, r, appConf.Uris[i].Delay, appConf.Uris[i])
			return
		}
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request, delay int, uri Uri) {
	time.Sleep(time.Duration(delay) * time.Millisecond)

	// from outer
	if uri.Source != "" {
		handleSource(w, r, uri)
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

func handleSource(w http.ResponseWriter, r *http.Request, uri Uri) {
	srcArry := strings.SplitN(uri.Source, " ", 2)
	if strings.HasPrefix(srcArry[1], "http://") {
		handleHttp(w, r, clientHttp, srcArry)
		return
	}

	if strings.HasPrefix(srcArry[1], "https://") {
		handleHttp(w, r, clientHttps, srcArry)
		return
	}

	// handle local file
	fmt.Fprintf(w, loadCert(uri.Body))
	return
}

// handle http and https
func handleHttp(w http.ResponseWriter, r *http.Request, client *http.Client, srcArry []string) {
	var resp *http.Response
	var err error
	if srcArry[0] == "GET" {
		fmt.Println("GET:[", srcArry[1], "]")
		resp, err = client.Get(srcArry[1])
	} else if srcArry[0] == "POST" {
		fmt.Println("POST:[", srcArry[1], "]")
		resp, err = client.Post(srcArry[1], "application/json", r.Body)
	} else {
		err = errors.New("Method not supported yet:" + srcArry[0])
	}

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
	if exist("config/server.cert") && exist("config/server.key") {
		http.ListenAndServeTLS(":20175", "cert/server.cert", "cert/server.key", nil)
	} else {
		http.ListenAndServe(":20175", nil)
	}
}

func loadCert(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	return string(data)
}

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
