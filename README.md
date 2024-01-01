## A simple chat service implemented using websockets
### Install packages
go get "golang.org/x/net/websocket"

### Start Server
go run main.go

### Send message to server
- Open multiple new tabs in chrome
- Open devtools in all of those
- Run the following script in console
```
const socket = new WebSocket("ws://localhost:3000/new-socket-conn")
socket.onmessage = e => console.log("received from server: ", e.data)
socket.send("hello from client")
```
