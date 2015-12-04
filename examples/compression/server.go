package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/euforia/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
		Extensions: []string{
			"permessage-deflate; server_no_context_takeover; client_no_context_takeover",
		},
	}
	webroot, _ = filepath.Abs("./")
	listenAddr = "0.0.0.0:12345"
)

func ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	log.Printf("Client connected: %s\n", r.RemoteAddr)

	var (
		msgType  int
		msgBytes []byte
	)

	for {
		// Blocking call
		if msgType, msgBytes, err = ws.ReadMessage(); err != nil {
			log.Printf("Client disconnected %s: %s\n", r.RemoteAddr, err)
			break
		}

		log.Printf("type: %d; payload: %d bytes;\n", msgType, len(msgBytes))

		ws.WriteMessage(msgType, msgBytes)
	}
}

func main() {
	// Serve index.html
	http.Handle("/", http.FileServer(http.Dir(webroot)))
	// Websocket endpoint
	http.HandleFunc("/ws", ServeWebSocket)

	log.Printf("Starting server on: %s\n", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatalln(err)
	}
}
