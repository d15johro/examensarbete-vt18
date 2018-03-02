package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/d15johro/examensarbete-vt18/osmdecoder/fbsconv/fbs"
	"github.com/d15johro/examensarbete-vt18/osmdecoder/pbconv/pb"
	"github.com/golang/protobuf/proto"
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
	bufferReadSize int32 = 1024 * 3 // TODO: increase when sending bigger chunks of data in experiment
)

type client struct {
	// conn is the websocket connection to which the client is connected.
	conn *websocket.Conn
	// request is a channel used to syncronize read and write operations (read --> write).
	request chan int32
}

func init() {
	flag.Parse()
	// TODO: clean log file with metrics data.
}

func main() {
	conn, err := websocket.Dial(*url, "", *origin)
	if err != nil {
		log.Fatalln("could not dial websocket server:", err)
	}
	defer conn.Close()
	c := client{
		conn:    conn,
		request: make(chan int32),
	}
	c.conn.PayloadType = websocket.BinaryFrame
	defer close(c.request)
	go func() { c.request <- 0 }() // init write operation
	go c.write()                   // spawn write function on a different user-space thread
	c.read()
}

// read reads data received from the server. Once the data is read and
// deserialized, read append all collected metrics to a log file.
// read then sends a bool to the client.request channel which triggers
// the write function to start operating.
func (c *client) read() {
	defer c.conn.Close()
	for {
		// read data from server:
		var buf = make([]byte, bufferReadSize)
		n, err := c.conn.Read(buf)
		if err != nil {
			log.Println("read:", err)
			break
		}
		data := buf[:n]                                                      // skip zeros
		responseTimeDuration := time.Since(startAccessTime).Seconds() * 1000 // response time in ms
		// deserialize data:
		var pbOSM pb.OSM
		deserializationTime, err := deserializeGetTime(data, &pbOSM, deserialize)
		if err != nil {
			log.Println("read:", err)
			break
		}
		accessTimeDuration := time.Since(startAccessTime).Seconds() * 1000 // access time in ms
		// TODO: save metrics to a log file.
		log.Printf("\n---\nmetrics:\n\tID: %d\n\taccess time: %fms\n\tresponse time %fms\n\tdeserialization time: %fms\n\tpbOSM.Generator: %s\n\tbytes read (data size): %d\n---",
			pbOSM.Id,
			accessTimeDuration,
			responseTimeDuration,
			deserializationTime,
			pbOSM.Generator,
			n)

		c.request <- pbOSM.Id + 1
	}
}

// write sends a message to server requesting more data.
func (c *client) write() {
	defer c.conn.Close()
	for id := range c.request {
		if id >= 20 { // limit requests
			break
		}
		msg := struct {
			ID int32 `json:"id"`
		}{ID: id}
		startAccessTime = time.Now() // start access time clock
		if err := websocket.JSON.Send(c.conn, &msg); err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func deserializeGetTime(data []byte, v interface{}, f func(data []byte, v interface{}) error) (ms float64, err error) {
	start := time.Now()
	err = f(data, v)
	if err != nil {
		return
	}
	elapsed := time.Since(start)
	ms = elapsed.Seconds() * 1000
	return
}

func deserialize(data []byte, v interface{}) error {
	if _, ok := v.(fbs.OSM); ok {
		v = fbs.GetRootAsOSM(data, 0)
		return nil
	}
	if msg, ok := v.(proto.Message); ok {
		return proto.Unmarshal(data, msg)
	}

	return fmt.Errorf("invalid type")
}
