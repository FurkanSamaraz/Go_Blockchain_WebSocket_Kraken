package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var aa = make(chan string)
var err error
var messageType int
var message []byte
var upgrader = websocket.Upgrader{} // use default options
func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()
	// The event loop
	for {
		messageType, message, err = conn.ReadMessage()
		if err != nil {
			//	log.Println("Error during message reading:", err)
			break
		}

		log.Printf("Alinan: %s", message)
		aa <- string(message)
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Error during message writing:", err)
			break
		}
	}
}

var sss string

func dea(w http.ResponseWriter, r *http.Request) {
	fmt.Println("bu:", string(message))

	sss := <-aa
	b, _ := json.Marshal(sss)
	w.Write(b)
}
func main() {
	http.HandleFunc("/socket", socketHandler)
	http.HandleFunc("/", dea)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

/*package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/webdeveloppro/golang-websocket-client/pkg/server"
)

var addr = flag.String("addr", ":8000", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()

	hub := server.NewHub()
	go hub.Run()
	http.HandleFunc("/frontend", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("got new connection")
		server.ServeWs(hub, w, r)
	})

	fmt.Println("server started ... ")
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		panic(err)
	}

}
*/
