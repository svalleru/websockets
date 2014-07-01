//package client
package main

import (
"log"
"code.google.com/p/go.net/websocket"
"encoding/json"
"os/exec"
	"time"
	"fmt"
)
var HOST = "localhost"
var ORIGIN = "http://"+HOST+"/"
var PORT = "8000"
//JSON request template
type template struct {
	Id, Origin, Timestamp, Controller, Payload string
}


func sockhandler(url string, data []byte) []byte{
//conf, err := websocket.NewConfig("ws://localhost:8000/", "http://localhost:8000/")
//log.Print(conf)
ws, err := websocket.Dial(url, "", ORIGIN)
if err != nil {
log.Fatal(err)
}
//	if _, err := ws.Write([]byte("date")); err != nil {
//		log.Fatal(err)
//	}
websocket.Message.Send(ws, data)
var msg = make([]byte, 512)
var n int
if n, err = ws.Read(msg); err != nil {
log.Fatal(err)
}
//print whatever you received frm server
return msg[:n]
}

func main() {
	var url = "ws://"+HOST+":"+PORT+"/readconfig"
	var data = []byte("client id: c1")
	configdata := sockhandler(url, data)
	//log.Print(string(config))
	var config template

	err := json.Unmarshal(configdata, &config)

	if err!= nil {
		log.Fatal(err)
	}

	log.Print("received payload :: " + config.Payload)
    log.Print("executing payload..")
	out, err := exec.Command(string(config.Payload)).Output()
	if err != nil {
		log.Fatal(err)
	}
	config.Payload = string(out)
	//log.Print(config)
	log.Print("sending response to server..")
	var url_resp = "ws://"+HOST+":"+PORT+"/storeconfig"
	//TBD JSON API for client->server communication
	var data_resp = []byte(config.Payload)
	log.Print(string(sockhandler(url_resp, data_resp)))

	log.Print("Entering heartbeat loop..")
	for{
		time.Sleep(2 * time.Second)
		var heartbeat_url = "ws://"+HOST+":"+PORT+"/"
		var heartbeat_resp = []byte(fmt.Sprint(time.Now()))
		log.Print(string(sockhandler(heartbeat_url, heartbeat_resp)))
	}
}

