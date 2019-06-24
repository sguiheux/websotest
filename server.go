package main

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"log"
	"net/http"
	"time"
)

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/socket.io", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
	server.OnEvent("/", "send-me", func(s socketio.Conn, msg string) {
		i := 0
		for i < 10 {
			i++
			s.Emit("reply", fmt.Sprintf("Message %d",i))
			time.Sleep(1 * time.Second)
		}

	})
	server.OnError("/", func(e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		fmt.Println("closed", msg)
	})
	go server.Serve()
	defer server.Close()

	wsServer := new(customServer)
	wsServer.Server = server

	http.Handle("/socket.io/", wsServer)
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

//Custom server which basically only contains a socketio variable
//But we need it to enhance it with functions
type customServer struct {
	Server *socketio.Server
}

//Header handling, this is necessary to adjust security and/or header settings in general
//Please keep in mind to adjust that later on in a productive environment!
//Access-Control-Allow-Origin will be set to whoever will call the server
func (s *customServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	origin := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	s.Server.ServeHTTP(w, r)
}
