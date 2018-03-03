package main

import (
	"encoding/binary"
	"flag"
	"log"
	"math"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http server address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 5,
	WriteBufferSize: 1024 * 5,
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
		var requestMessage struct {
			ID uint32 `json:"id"`
		}
		if err := conn.ReadJSON(&requestMessage); err != nil {
			log.Println(err)
			break
		}
		var serializationTime float64 = 42

		data := []byte("wsup")
		data = appendFloat64ToBytes(data, serializationTime)
		data = appendUint32ToBytes(data, requestMessage.ID)
		log.Println(requestMessage.ID)
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
