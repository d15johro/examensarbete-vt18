package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/d15johro/examensarbete-vt18/websocket/codec"
	"github.com/d15johro/examensarbete-vt18/websocket/pb_send"
	"golang.org/x/net/websocket"
)

var addr = flag.String("addr", ":8080", "the address to host the websocket server")

type client struct {
	conn *websocket.Conn
	send chan message
}

type message struct {
	ID int32 `json:"id"`
}

func init() {
	flag.Parse()
}

func main() {
	http.Handle("/websocket", websocket.Handler(func(conn *websocket.Conn) {
		defer conn.Close()
		c := client{
			conn: conn,
			send: make(chan message),
		}
		defer func() { close(c.send) }()
		go c.write()
		c.read()
	}))
	log.Fatalln(http.ListenAndServe(*addr, nil))
}

func (c *client) read() {
	defer c.conn.Close()
	for {
		var msg message
		if err := websocket.JSON.Receive(c.conn, &msg); err != nil {
			log.Println("read:", err)
			break
		}
		c.send <- msg
	}
}

func (c *client) write() {
	defer c.conn.Close()
	for msg := range c.send {
		if msg.ID%10 == 0 {
			log.Println("message ID:", msg.ID)
		}
		send := pb_send.Send{Data: "some chunk of text"}
		if err := codec.PB.Send(c.conn, &send); err != nil {
			log.Println("write:", err)
			break
		}
	}
}
