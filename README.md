Client-Server Duo Developed in Go Using Websockets
==================================================

Goals: 
- To have a central server that manages all the clients(a.k.a agents) running on a taget docker machine 
- Server should be able to run docker(and bash) commands remotely to provision, configure and manage the docker containers
- Using websockets is advantageous to overcome any firewall/security restrictions if the docker machine is inside the secure zone
- Will update with more..

How it works:
- Once the client (agent) is instaleld on docker machine, it initiates a http/https connection to central server and upgrades the http/https to websocket protocol (ws) enabling duplex communication between client and server
- Once websocket is established, server sends command payload in byte stream to the remote client, client runs the payload on it and sends back the response to server, server persists the client's response
- websocket is terminted when there's no activity for a while

Fundas:
- websocket protocol: http://tools.ietf.org/html/rfc6455
- go programming language: http://golang.org/