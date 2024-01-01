package main

import (
	"fmt"
	"io"
	"net/http"
	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func registerHandler(server *Server) {
	http.Handle("/new-socket-conn", websocket.Handler(server.handleWS))
	http.ListenAndServe(":3000", nil)
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Print("new incoming connection request from", ws.RemoteAddr(), "\n")
	s.conns[ws] = true
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 2048)
	for {
		n, err := ws.Read(buf);
		if err != nil {
			if err == io.EOF {
				break;
			}
			fmt.Println("Error while reading message: ", err)
			continue
		}
		msg := buf[:n]
		fmt.Println(string(msg))
		ws.Write([]byte("Received new message"))
	}
}

func main() {
	server := NewServer()
	registerHandler(server)
}