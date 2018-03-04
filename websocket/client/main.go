package main

import (
	"encoding/binary"
	"flag"
	"log"
	"math"
	"time"

	"github.com/d15johro/examensarbete-vt18/osmdecoder/fbsconv/fbs"
	"github.com/d15johro/examensarbete-vt18/osmdecoder/pbconv/pb"
	"github.com/golang/protobuf/proto"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/gorilla/websocket"
)

var (
	dialURL             = flag.String("du", "ws://localhost:8080/websocket", "url to dial websocket server")
	serializationFormat = flag.String("sf", "pb", "Serialization format")
)

func init() {
	flag.Parse()
}

func main() {
	log.Printf("dialing websocket server on %s using %s as serialisering format...", *dialURL, *serializationFormat)
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
			ID                  uint32 `json:"id"`
			SerializationFormat string `json:"serializationFormat"`
		}{ID: uint32(i), SerializationFormat: *serializationFormat}
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
		serializationTime := extractFloat64FromBytes(data, len(data)-12, len(data)-4)
		data = data[:len(data)-8-4]
		// deserialize data:
		startDeserializationClock := time.Now()
		switch *serializationFormat {
		case "pb":
			osm := &pb.OSM{}
			if err := proto.Unmarshal(data, osm); err != nil {
				log.Println(err)
				break
			}
			log.Println(osm.Copyright)
		case "fbs":
			offset := flatbuffers.UOffsetT(0)
			n := flatbuffers.GetUOffsetT(data[offset:])
			osm := &fbs.OSM{}
			osm.Init(data, n+offset)
			log.Println(string(osm.Copyright()))
		default:
			log.Fatalln("serialization format not supported")
		}
		deserializationTime := time.Since(startDeserializationClock).Seconds() * 1000
		// print collected metrics:
		log.Printf("ID: %d, serialization time: %f, deserialization time: %f\n---\n",
			id,
			serializationTime,
			deserializationTime,
		)
	}
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
