package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/d15johro/examensarbete-vt18/osmdecoder/fbsconv/fbs"
	"github.com/d15johro/examensarbete-vt18/osmdecoder/pbconv/pb"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

var (
	dialURL             = flag.String("du", "ws://localhost:8080/websocket", "url to dial websocket server")
	serializationFormat = flag.String("sf", "pb", "Serialization format")
)

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
		// request data from server:
		requestMessage := struct {
			ID uint32 `json:"id"`
		}{ID: uint32(i)}
		if err := conn.WriteJSON(&requestMessage); err != nil {
			log.Println(err)
			break
		}
		// read response data from server:
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		// extract id and serialization time from data:
		id := extractUint32FromBytes(data, len(data)-4, len(data))
		if id != requestMessage.ID {
			log.Println("ID from requestMessage doesn't match ID recieved from server")
		}
		serializationTime := extractFloat64FromBytes(data, len(data)-12, len(data)-4+1)
		data = data[:len(data)-8-4]
		// deserialize data:
		startDeserializationClock := time.Now()
		var osm pb.OSM // change type depending on serialization format being used
		if err := deserialize(data, &osm); err != nil {
			log.Fatalln(err)
		}
		deserializationTime := time.Since(startDeserializationClock).Seconds() * 1000
		// print collected metrics:
		log.Printf("ID: %d, serialization time: %f, deserialization time: %f, #nodes: %d\n---\n",
			id,
			serializationTime,
			deserializationTime,
			len(osm.Nodes),
		)
	}
}

func deserialize(data []byte, v interface{}) (err error) {
	if osm, ok := v.(*pb.OSM); ok {
		return proto.Unmarshal(data, osm)
	}
	if _, ok := v.(*fbs.OSM); ok {
		v = fbs.GetRootAsOSM(data, 0)
		return
	}
	err = fmt.Errorf("deserialize: type not supported")
	return
}

func uint32FromBytes(bytes []byte) uint32 {
	return binary.LittleEndian.Uint32(bytes)
}

func extractUint32FromBytes(data []byte, start, end int) uint32 {
	return uint32FromBytes(data[start:end])
}

func extractFloat64FromBytes(data []byte, start, end int) float64 {
	return float64FromBytes(data[start:end])
}

func float64FromBytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	f := math.Float64frombits(bits)
	return f
}
