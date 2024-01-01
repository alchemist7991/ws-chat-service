package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func registerHandler(server *Server) {
	http.Handle("/new-socket-conn", websocket.Handler(server.handleWS))
	go http.Handle("/live-feed", websocket.Handler(server.liveFeed))
	http.ListenAndServe(":3000", nil)
}

func (s *Server) liveFeed(ws *websocket.Conn) {
	for {
		// some continuous data
		payload := fmt.Sprintf("%s", time.Now().UnixMicro())
		ws.Write([]byte(payload))
		time.Sleep(1 * time.Second)
	}
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
		s.broadcast(msg)
		fmt.Println(string(msg))
	}
}

func (s *Server) broadcast(msg []byte) {
	for ws := range s.conns {
		go func (ws *websocket.Conn) {
			_, err := ws.Write(msg)
			if err != nil {
				fmt.Println("Error in writing message")
			}
		}(ws)
	}
}

func main() {
	server := NewServer()
	registerHandler(server)
}