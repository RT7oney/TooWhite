package main

import (
	"encoding/json"
	"fmt"
)

type Server struct {
	clients map[*Client]bool

	broadcast chan *Content

	register chan *Client

	unregister chan *Client
}

func newServer() *Server {
	return &Server{
		broadcast:  make(chan *Content),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (serv *Server) run() {
	for {
		select {
		case client := <-serv.register:
			serv.clients[client] = true
			fmt.Println("把客户端放在线程池里", serv.clients)
		case client := <-serv.unregister:
			if _, ok := serv.clients[client]; ok {
				delete(serv.clients, client)
				close(client.send)
			}
		case content := <-serv.broadcast:
			fmt.Println("所有的客户端", serv.clients)
			fmt.Println("TARGET", content.Target)
			// 根据content传过来的需要广播的target来遍历广播
			if content.Target != nil {
				for _, client := range content.Target {
					message, _ := json.Marshal(content.Data)
					if _, ok := serv.clients[client]; ok {
						fmt.Println("被广播的用户", client)
						select {
						case client.send <- []byte(message):
						default:
							close(client.send)
							delete(serv.clients, client)
						}
					}
				}
			}
		}
	}
}
