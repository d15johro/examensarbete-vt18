package main

import (
	"flag"
	"log"

	"github.com/d15johro/examensarbete-vt18/websocket/codec"
	"github.com/d15johro/examensarbete-vt18/websocket/pb_send"

	"golang.org/x/net/websocket"
)

var (
	url    = flag.String("url", "ws://localhost:8080/websocket", "websocket server url")
	origin = flag.String("origin", "http://localhost/", "websocket origin")
)

type client struct {
	conn    *websocket.Conn
	request chan bool
	message message
}

type message struct {
	ID int32 `json:"id"`
}

func init() {
	flag.Parse()
}

func main() {
	conn, err := websocket.Dial(*url, "", *origin)
	if err != nil {
		log.Fatalln("could not dial websocket server:", err)
	}
	defer conn.Close()
	c := client{
		conn:    conn,
		request: make(chan bool),
		message: message{ID: 0},
	}
	defer close(c.request)
	go func() { c.request <- true }() // init request
	go c.write()
	c.read()
}

func (c *client) read() {
	defer c.conn.Close()
	for {
		var x pb_send.Send
		if err := codec.PB.Receive(c.conn, &x); err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("%d: %s\n", c.message.ID, x.Data)
		c.request <- true
	}
}

func (c *client) write() {
	defer c.conn.Close()
	for range c.request {
		if c.message.ID >= 30 {
			break
		}
		if err := websocket.JSON.Send(c.conn, c.message); err != nil {
			log.Println("write:", err)
			break
		}
		c.message.ID++
	}
}
