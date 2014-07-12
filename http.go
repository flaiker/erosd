package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var webRoot = "web"

func handleWebsocketConnection(conn *websocket.Conn) {
	client := NewWebsocketClientConnection(*conn)
	client.read()
}

func httpServeHomeFunc(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, webRoot+"/index.html")
}

func httpServeWsFunc(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade WS from", r.RemoteAddr, err)
		return
	}

	log.Println("Accepted WS from", ws.RemoteAddr())
	go handleWebsocketConnection(ws)
}

func listenAndServeHTTP(address string) {
	r := mux.NewRouter()
	r.HandleFunc("/ws", httpServeWsFunc).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(webRoot+"/static")))).Methods("GET")
	r.HandleFunc("/{path:.*}", httpServeHomeFunc).Methods("GET")
	log.Println("Listening HTTP on ", address)
	log.Fatalln(http.ListenAndServe(address, r))
}
