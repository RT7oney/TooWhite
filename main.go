package main

import (
	// "fmt"
	"log"
	"net/http"
)

/*===========广播的协议===============*/
type Msg struct {
	MsgType int
	UserId  string
	GroupId string
	Content Content
}
type Content struct {
	ContentType int
	From        string
	ToUser      string
	ToGroup     string
	Data        string
}

/*==================================*/

func main() {
	serv := newServer()
	go serv.run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(serv, w, r)
	})
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
