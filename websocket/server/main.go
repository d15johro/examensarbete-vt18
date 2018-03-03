package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/d15johro/examensarbete-vt18/osmdecoder"
	"github.com/d15johro/examensarbete-vt18/osmdecoder/fbsconv"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http server address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 10,
	WriteBufferSize: 1024 * 10,
}

func init() {
	flag.Parse()
}

func main() {
	http.HandleFunc("/websocket", handler)
	log.Println("server listening on", *addr)
	log.Fatalln(http.ListenAndServe(*addr, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	for {
		// read client request:
		var requestMessage struct {
			ID uint32 `json:"id"`
		}
		if err := conn.ReadJSON(&requestMessage); err != nil {
			log.Println(err)
			break
		}
		// decode .osm file depending on id from request message:
		file := "../../testdata/test_data" + fmt.Sprintf("%d", requestMessage.ID%6) + ".osm"
		x, err := osmdecoder.DecodeFile(file)
		if err != nil {
			log.Println("write:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Println("nodes l", len(x.Nodes))
		// serialize:
		startSerializationClock := time.Now()
		builder := flatbuffers.NewBuilder(0)
		err = fbsconv.Build(builder, x)
		if err != nil {
			log.Println("write:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		data := builder.Bytes[builder.Head():]
		serializationTime := time.Since(startSerializationClock).Seconds() * 1000
		// send data to client:
		data = appendFloat64ToBytes(data, serializationTime)
		data = appendUint32ToBytes(data, requestMessage.ID)
		log.Println(len(data))
		if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
			log.Println(err)
			break
		}
	}
}

func uin32ToBytes(i uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, i)
	return bytes
}

func appendUint32ToBytes(data []byte, i uint32) []byte {
	buf := uin32ToBytes(i)
	for i := 0; i < len(buf); i++ {
		data = append(data, buf[i])
	}
	return data
}

func float64ToBytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func appendFloat64ToBytes(data []byte, f float64) []byte {
	buf := float64ToBytes(f)
	for i := 0; i < len(buf); i++ {
		data = append(data, buf[i])
	}
	return data
}
