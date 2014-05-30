//package server
package main
import (
	"io/ioutil"
	"net/http"
	"log"
	"code.google.com/p/go.net/websocket"
	"io"
	"bytes"
)

var SERVER_PORT = "443"
var CLIENT_CONFIG = "/Users/svalleru/Desktop/golang/src/server/c1_req.json"
var CLIENT_RESP = "/Users/svalleru/Desktop/golang/src/server/c1_resp.json"

func readconfig(ws *websocket.Conn) {
log.Print("reading config..")
content, err := ioutil.ReadFile(CLIENT_CONFIG)
if err == nil {
	//log.Print(string(ws))
	var incoming []byte
	if err := websocket.Message.Receive(ws, &incoming); err != nil {
      log.Print("client stream ended..")//break
	}
	log.Print("data received from client: ", string(incoming))

	io.Copy(ws, bytes.NewBuffer(content))
	log.Print("config sent..")
 } else{
	log.Print(err)
	panic(err)
 }
}

func storeconfig(ws *websocket.Conn) {
	log.Print("processing client's response..")

	var incoming []byte
	if err := websocket.Message.Receive(ws, &incoming); err != nil {
		log.Print("client stream ended..")//break
	}
	err := ioutil.WriteFile(CLIENT_RESP, incoming, 0644)
	if err == nil {
		log.Print("response stored..")
		io.Copy(ws, bytes.NewBuffer([]byte("response stored at server..")))
	} else{
		log.Print(err)
		io.Copy(ws, bytes.NewBuffer([]byte("unable to persist..")))
		panic(err)
	}
}

func main(){
	http.Handle("/readconfig", websocket.Handler(readconfig))
	http.Handle("/storeconfig", websocket.Handler(storeconfig))
	log.Print("server running on port "+SERVER_PORT+"..")
	err := http.ListenAndServe(":"+SERVER_PORT, nil)
	if err != nil {
		panic("unable to serve.. " + err.Error())
	}

}

