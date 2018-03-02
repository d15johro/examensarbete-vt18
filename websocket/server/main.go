package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/d15johro/examensarbete-vt18/osmdecoder"
	"github.com/d15johro/examensarbete-vt18/osmdecoder/pbconv"
	"github.com/golang/protobuf/proto"
	flatbuffers "github.com/google/flatbuffers/go"
	"golang.org/x/net/websocket"
)

var addr = flag.String("addr", ":8080", "the address to host the websocket server")

type client struct {
	// conn is the websocket connection to which the client connects.
	conn *websocket.Conn
	// send is a channel used to syncronize read and write operations (read --> write).
	send chan int32
}

func init() {
	flag.Parse()
	// TODO: clean log file with metrics data.
}

func main() {
	http.Handle("/websocket", websocket.Handler(func(conn *websocket.Conn) {
		defer conn.Close()
		c := client{
			conn: conn,
			send: make(chan int32),
		}
		c.conn.PayloadType = websocket.BinaryFrame
		defer func() { close(c.send) }()
		go c.write() // spawn write function on a different user-space thread
		c.read()
	}))
	log.Println("starting websocket server on:", *addr)
	log.Fatalln(http.ListenAndServe(*addr, nil))
}

// read reads messages received from the client. Once the message is read and
// deserialized, read sends the message to the client.send channel which triggers
// the write function to start operating.
func (c *client) read() {
	defer c.conn.Close()
	for {
		msg := struct {
			ID int32 `json:"id"`
		}{}
		if err := websocket.JSON.Receive(c.conn, &msg); err != nil {
			log.Println("read:", err)
			break
		}
		c.send <- msg.ID
	}
}

// write writes and sends data to the client depending on the ID of the message
// being sent to the client.send channel that triggered the write operation.
// Collected metrics are appended to a log file.
func (c *client) write() {
	defer c.conn.Close()
	for id := range c.send {
		// TODO: select what data to send depending on msg.ID.
		// read .osm file and deserialize into a osmdecoder.OSM struct:
		fileNumber := fmt.Sprintf("%d", id%6)
		osm, err := osmdecoder.DecodeFile("../../testdata/test_data" + fileNumber + ".osm")
		if err != nil {
			log.Println("write:", err)
			break
		}
		// make a pb.OSM out of osm:
		pbOSM, err := pbconv.Make(osm)
		if err != nil {
			log.Println("write:", err)
			break
		}
		pbOSM.Id = id
		// serialize pbOSM:
		data, serializationTime, err := serializeGetTime(pbOSM, serialize)
		if err != nil {
			log.Println("write:", err)
			break
		}
		// write data to client:
		n, err := c.conn.Write(data)
		if err != nil {
			log.Println("write:", err)
			break
		}
		// TODO: append data to a log file.
		log.Printf("\n---\nmetrics:\n\tfile number: %s\n\tID: %d\n\tserialize time: %f ms\n\tbytes written: %d\n---",
			fileNumber,
			id,
			serializationTime,
			n)
	}
}

func serializeGetTime(v interface{}, f func(v interface{}) (data []byte, err error)) (data []byte, ms float64, err error) {
	start := time.Now()
	data, err = f(v)
	if err != nil {
		return
	}
	elapsed := time.Since(start)
	ms = elapsed.Seconds() * 1000
	return
}

func serialize(v interface{}) (data []byte, err error) {
	if b, ok := v.(flatbuffers.Builder); ok {
		return b.Bytes[b.Head():], nil
	}
	if osm, ok := v.(proto.Message); ok {
		data, err = proto.Marshal(osm)
		return
	}
	return nil, fmt.Errorf("invalid type")
}
