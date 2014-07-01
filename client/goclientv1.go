//package client
package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	_ "fmt"
	"log"
	"os/exec"
	"time"
	"strings"
)

var HOST = "localhost"
var ORIGIN = "http://" + HOST + "/"
var PORT = "8000"

//JSON request template
type template struct {
	Id, Origin, Timestamp, Controller, Payload string
}

func sockhandler(url string, data []byte) []byte {
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
	var msg = make([]byte, 2048 * 10)
	var n int
	if n, err = ws.Read(msg); err != nil {
		log.Fatal(err)
	}
	//print whatever you received frm server
	return msg[:n]
}

//async func
func async_exec(c chan string, cdata string) string{
	//log.Print("cdata is...", cdata)
	var config template
	err := json.Unmarshal([]byte(cdata), &config)
	if err != nil {
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
	var url_resp = "ws://" + HOST + ":" + PORT + "/storeconfig"
	//TBD JSON API for client->server communication
	var data_resp = []byte(config.Payload)
	c <- string(sockhandler(url_resp, data_resp))
	return <-c
}


func payload_handler() {

	var url = "ws://" + HOST + ":" + PORT + "/readconfig"
	data, _ := exec.Command("hostname").Output()
	configdata := sockhandler(url, data)
	//log.Print(string(config))

    res := strings.Split(string(configdata), "##")
	res = res[:len(res)-1]
	start := time.Now()
	c := make(chan string)
	var results []string

	for i := range res {
		cdata := res[i]
		//for every v, sping a async call
		go func() {
			c <- async_exec(c, cdata)
		}()
	}

	//Disabled timeouts for now
	//timeout := time.After(10 * time.Millisecond)

	for i := 0; i < len(res); i++ {
		select {
		case result := <-c:
			results = append(results, result)
			//		case <-timeout:
			//			fmt.Println("timed out.")
			//			continue
		}
	}
	log.Print("bash jobs ran: ", len(results))
	elapsed := time.Since(start)
	log.Print("time elapsed: ", elapsed)
	log.Print("--------------")
	return
}

func main() {

	log.Print("Entering heartbeat loop..")
	for {
		time.Sleep(2 * time.Second)
		var heartbeat_url = "ws://" + HOST + ":" + PORT + "/"
		out, _ := exec.Command("hostname").Output()
		var heartbeat_resp = []byte(out)
		response := string(sockhandler(heartbeat_url, heartbeat_resp))

		switch response {
		case "001": //client key just initialized
			log.Print("001: client initialized..")
		case "002": //client initialized but no jobs
			log.Print("002: no jobs..yet..")
		case "003": //client initialized and jobs in queue
			log.Print("003: receiving payload..running handler..")
			payload_handler()
		default:
			log.Print("no response code match")

		}
	}
}
