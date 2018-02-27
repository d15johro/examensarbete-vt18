package main

import (
	"flag"
	"log"
	"time"

	"github.com/d15johro/examensarbete-vt18/websocket/codec"
	"github.com/d15johro/examensarbete-vt18/websocket/pb_send"

	"golang.org/x/net/websocket"
)

var (
	url    = flag.String("url", "ws://localhost:8080/websocket", "websocket server url")
	origin = flag.String("origin", "http://localhost/", "websocket origin")
)

var (
	startAccessTime time.Time // global since timer starts and ends in different functions
)

const (
	bufferReadSize int32 = 1024 // TODO: increase when sending bigger chunks of data in experiment
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
		var buf = make([]byte, bufferReadSize)
		n, err := c.conn.Read(buf)
		if err != nil {
			log.Println("read:", err)
			break
		}
		data := buf[:n]                        // skip zeros
		startDeserializationTime := time.Now() // start deserialization time clock
		var x pb_send.Send
		if err = codec.PB.Unmarshal(data, websocket.BinaryFrame, &x); err != nil {
			log.Println("read:", err)
			break
		}
		deserializationDuration := time.Since(startDeserializationTime).Seconds() * 1000 // deserialization time
		accessTimeDuration := time.Since(startAccessTime).Seconds() * 1000               // access time

		log.Printf("\n---\nmetrics:\n\tID: %d\n\taccess time: %f ms\n\tdeserialization time: %f\n\tmessage: %s\n---",
			c.message.ID,
			accessTimeDuration,
			deserializationDuration,
			x.Data)

		c.request <- true
	}
}

func (c *client) write() {
	defer c.conn.Close()
	for range c.request {
		if c.message.ID >= 5 {
			break
		}
		startAccessTime = time.Now() // start access time clock
		if err := websocket.JSON.Send(c.conn, c.message); err != nil {
			log.Println("write:", err)
			break
		}
		c.message.ID++
	}
}
