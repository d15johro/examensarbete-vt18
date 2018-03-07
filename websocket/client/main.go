package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
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

type metrics struct {
	id                  uint32
	accessTime          float64
	responseTime        float64
	serializationTime   float64
	deserializationTime float64
	dataSize            int
	filepath            string
}

func init() {
	flag.Parse()
}

func main() {
	m := metrics{filepath: "./websocket_" + *serializationFormat + ".txt"}
	if err := m.setup(); err != nil {
		log.Fatalln(err)
	}
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
	for i := 0; i < 1000; i++ {
		// Request data from server:
		startAccessClock := time.Now()
		startResponseClock := time.Now()
		requestMessage := struct {
			ID                  uint32 `json:"id"`
			SerializationFormat string `json:"serializationFormat"`
		}{ID: uint32(i), SerializationFormat: *serializationFormat}
		if err := conn.WriteJSON(&requestMessage); err != nil {
			log.Println(err)
			break
		}
		// Read response data from server:
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		m.responseTime = time.Since(startResponseClock).Seconds() * 1000
		// Extract id and serialization time from data:
		m.id = extractUint32FromBytes(data, len(data)-4, len(data))
		if m.id != requestMessage.ID {
			log.Println("ID from requestMessage doesn't match ID recieved from server")
		}
		m.serializationTime = extractFloat64FromBytes(data, len(data)-12, len(data)-4)
		data = data[:len(data)-8-4]
		m.dataSize = len(data)
		// Deserialize data:
		startDeserializationClock := time.Now()
		switch *serializationFormat {
		case "pb":
			osm := &pb.OSM{}
			if err := proto.Unmarshal(data, osm); err != nil {
				log.Println(err)
				break
			}
		case "fbs":
			// Unlike pb, deserializing fbs basically means storing the raw binary data in a struct.
			// Therefore, the deserialization time for fbs will probably always be 0ms.
			offset := flatbuffers.UOffsetT(0)
			n := flatbuffers.GetUOffsetT(data[offset:])
			osm := &fbs.OSM{}
			osm.Init(data, n+offset)
		default:
			log.Fatalln("serialization format not supported")
		}
		m.deserializationTime = time.Since(startDeserializationClock).Seconds() * 1000
		m.accessTime = time.Since(startAccessClock).Seconds() * 1000
		if i > 499 {
			m.log()
		}
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

func (m *metrics) log() {
	s := fmt.Sprintf("%d,%f,%f,%f,%f,%d\n",
		m.id, m.accessTime,
		m.responseTime,
		m.serializationTime,
		m.deserializationTime,
		m.dataSize)
	file, err := os.OpenFile(m.filepath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	if _, err = file.WriteString(s); err != nil {
		log.Fatalln(err)
	}
}

func (m *metrics) setup() error {
	_, err := os.Stat(m.filepath)
	if err != nil {
		if !os.IsNotExist(err) {
			if err := os.Remove(m.filepath); err != nil {
				return err
			}
		}
	}
	_, err = os.Create(m.filepath)
	return err
}
