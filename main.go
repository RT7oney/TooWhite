package main

import (
	// "fmt"
	"log"
	"net/http"
)

type Msg struct {
	MsgType int
	From    string
	Target  string
	Data    string
}
type Content struct {
	ContentType int
	From        *Client
	Target      []*Client
	Data        interface{}
}

func main() {
	serv := newServer()
	go serv.run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(serv, w, r)
	})
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
