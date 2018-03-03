package main

import (
	"encoding/binary"
	"flag"
	"log"
	"math"

	"github.com/gorilla/websocket"
)

var dialURL = flag.String("du", "ws://localhost:8080/websocket", "url to dial websocket server")

func main() {
	log.Println("dialing websocket server on", *dialURL)
	conn, resp, err := websocket.DefaultDialer.Dial(*dialURL, nil)
	if err == websocket.ErrBadHandshake {
		log.Fatalln("handshake failed with status %d\n", resp.StatusCode)
	}
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	for i := 0; i < 10; i++ {
		requestMessage := struct {
			ID uint32 `json:"id"`
		}{ID: uint32(i)}
		if err := conn.WriteJSON(&requestMessage); err != nil {
			log.Println(err)
			break
		}
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		id := extractUint32FromBytes(data, len(data)-4, len(data))
		log.Println("id:", id)
		if id != requestMessage.ID {
			log.Println("ID from requestMessage doesn't match ID recieved from server")
		}
		serializationTime := extractFloat64FromBytes(data, len(data)-12, len(data)-4+1)
		log.Println("sf:", serializationTime)
		data = data[:8+4]
		log.Println("data:", string(data))
		log.Println("---")
	}
}

func extractFloat64FromBytes(data []byte, start, end int) float64 {
	return float64FromBytes(data[start:end])
}

func float64FromBytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func extractUint32FromBytes(data []byte, start, end int) uint32 {
	return uint32FromBytes(data[start:end])
}

func uint32FromBytes(bytes []byte) uint32 {
	return binary.LittleEndian.Uint32(bytes)
}
