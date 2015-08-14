package main
import (
    "os"
    "encoding/json"
    "fmt"
)

type Uri struct {
    Uri string `json:"uri"`
    Delay int `json:"delay"`
    Body string `json:"body"`
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
