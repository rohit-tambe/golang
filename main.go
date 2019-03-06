package main

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Hello there")
	server, err := socketio.NewServer(nil)

	if err != nil {
		log.Fatal(err)
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})
	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})
	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})
	go server.Serve()
	defer server.Close()
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.Handle("/socket.io/", server)
	log.Println("Serving at localhost:9000...")
	log.Fatal(http.ListenAndServe(":9000", nil))

}
