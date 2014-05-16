package main

import (
"io"
"net/http"

"code.google.com/p/go.net/websocket"
)

// Echo the data received on the WebSocket.
func EchoServer(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

// This example demonstrates a trivial echo server.
func main() {
	http.Handle("/echo", websocket.Handler(EchoServer))
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

