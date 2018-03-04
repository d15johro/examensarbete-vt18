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
	"github.com/d15johro/examensarbete-vt18/osmdecoder/pbconv"
	"github.com/golang/protobuf/proto"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/gorilla/websocket"
)

var (
	addr = flag.String("addr", "localhost:8080", "http server address")
)

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
			ID                  uint32 `json:"id"`
			SerializationFormat string `json:"serializationFormat"`
		}
		if err := conn.ReadJSON(&requestMessage); err != nil {
			log.Println(err)
			break
		}
		// decode .osm file depending on id from request message:
		file := "../../testdata/test_data" + fmt.Sprintf("%d", requestMessage.ID%6) + ".osm"
		x, err := osmdecoder.DecodeFile(file)
		if err != nil {
			log.Println(err)
			break
		}
		// serialize:
		var (
			data                    []byte
			startSerializationClock time.Time
		)
		switch requestMessage.SerializationFormat {
		case "pb":
			osm, err := pbconv.Make(x)
			if err != nil {
				log.Println(err)
				break
			}
			startSerializationClock = time.Now()
			data, err = proto.Marshal(osm)
			if err != nil {
				log.Println(err)
				break
			}
		case "fbs":
			startSerializationClock = time.Now() // // Unlike pb, fbs serialize data when building the buffer
			builder := flatbuffers.NewBuilder(0)
			err = fbsconv.Build(builder, x)
			if err != nil {
				log.Println(err)
				break
			}
			data = builder.Bytes[builder.Head():]
		default:
			log.Fatalln("serialization format not supported")
		}
		serializationTime := time.Since(startSerializationClock).Seconds() * 1000
		// send data to client:
		data = appendFloat64ToBytes(data, serializationTime)
		data = appendUint32ToBytes(data, requestMessage.ID)
		if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
			log.Println(err)
			break
		}
	}
}

func uint32ToBytes(ui uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, ui)
	return bytes
}

func appendUint32ToBytes(data []byte, ui uint32) []byte {
	buf := uint32ToBytes(ui)
	for i := 0; i < len(buf); i++ {
		data = append(data, buf[i])
	}
	return data
}

func float64ToBytes(f float64) []byte {
	bits := math.Float64bits(f)
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
