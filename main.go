package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func reader(conn *websocket.Conn) {
	defer conn.Close()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading: %v", err)
		}

		log.Printf("Message received: %v", string(message))

		if err := conn.WriteMessage(websocket.TextMessage, []byte("Received message")); err != nil {
			log.Println(err)
			return
		}
	}
}

func writer(conn *websocket.Conn) {
	defer conn.Close()
	for {
		err := conn.WriteMessage(websocket.TextMessage,
			[]byte(`
			<td id="1_1" style="border:1px solid LightGrey;background:black;"></td>
			<td id="2_2" style="border:1px solid LightGrey;background:black;"></td>
			`))
		if err != nil {
			log.Printf("Error writing message: %v", err)
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func gameOfLifeWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading ws: %s", err)
		return
	}
	defer func() {
		log.Println("closing connection")
		c.Close()
	}()
	go writer(c)
	go reader(c)

	//Keep handler alive
	select {}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/gameoflife", gameOfLifeWS)
	log.Println("Listening on port 3000")
	http.ListenAndServe(":3000", nil)
}
