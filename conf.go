package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Uri struct {
	Uri     string              `json:"uri"`
	Delay   int                 `json:"delay"`
	Headers []map[string]string `json:"headers"`
	Body    string              `json:"body"`
}

type App struct {
	Uris []Uri `json:"uris"`
}

var appConf App

func initConf() {

	file, _ := os.Open("app.json")
	decoder := json.NewDecoder(file)
	appConf = App{}
	err := decoder.Decode(&appConf)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(appConf)
}
