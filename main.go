package main

import (
	"fmt"
	"net/http"

	"github.com/docker/distribution/uuid"
	"github.com/eulbyvan/go-real-time-chat-app/model"
	"github.com/gorilla/websocket"
)

var manager = model.ClientManager{
	Broadcast:  make(chan []byte),
	Register:   make(chan *model.Client),
	Unregister: make(chan *model.Client),
	Clients:    make(map[*model.Client]bool),
}

func main() {
	fmt.Println("Starting application...")

	go manager.Start()

	http.HandleFunc("/ws", wsPage)

	http.ListenAndServe(":12345", nil)
}

func wsPage(res http.ResponseWriter, req *http.Request) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		return
	}

	client := &model.Client{Id: uuid.Generate().String(), Socket: conn, Send: make(chan []byte)}

	manager.Register <- client

	go client.Read(&manager)
	go client.Write()
}
