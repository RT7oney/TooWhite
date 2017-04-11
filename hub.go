package main

import (
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
			for client := range serv.clients {
				/*===========================*/
				// fmt.Println("开始分段广播", client)
				message := content.Data
				fmt.Println("message from servub", message)
				if content.ContentType == 0 {
					// 广播给用户
					if content.ToUser == client.uid {
						message = content.From + "对" + content.ToUser + "说：" + message
						select {
						case client.send <- []byte(message):
						default:
							close(client.send)
							delete(serv.clients, client)
						}
					}
				} else {
					// 广播给分组
					if content.ToGroup == client.gid {
						message = content.From + "在群" + content.ToGroup + "说：" + message
						fmt.Println("开始分组广播", message)
						select {
						case client.send <- []byte(message):
						default:
							close(client.send)
							delete(serv.clients, client)
						}
					}
				}

				/*===========================*/
			}
		}
	}
}
