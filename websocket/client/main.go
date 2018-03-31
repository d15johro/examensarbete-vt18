package main

import (
	"encoding/binary"
	"flag"
	"log"
	"math"
	"time"

	"github.com/d15johro/examensarbete-vt18/expcodec"
	"github.com/d15johro/examensarbete-vt18/metrics"
	"github.com/gorilla/websocket"
)

var (
	dialURL             = flag.String("du", "ws://localhost:8080/websocket", "url to dial websocket server")
	serializationFormat = flag.String("sf", "", "Serialization format")
	iterations          = flag.Int("itr", 0, "# iterations")
	warmups             = flag.Int("wm", 0, "# warmups")
)

func init() {
	flag.Parse()
	// Flags are required:
	if *serializationFormat == "" || *iterations == 0 || *warmups == 0 {
		log.Fatal("The following flags must be provided:\nflag\t\tvalue\t\tmeaning\nsf\t\tpb or fb\tThe serializationformat to use\nitr\t\tgreater than 0\t# of iterations\nwm\t\tgreater than 0\t# of warmup iterations")
	}
}

func main() {
	m := metrics.New()
	m.Filepath = "./websocket_" + *serializationFormat + ".txt"
	if err := m.Setup(); err != nil {
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
	for i := 0; i < *iterations; i++ {
		log.Println(i)
		// Request data from server:
		startAccessClock := time.Now()
		startResponseClock := time.Now()
		requestMessage := struct {
			ID                  uint32 `json:"id"`
			SerializationFormat string `json:"serializationFormat"`
		}{ID: uint32(i + 1), SerializationFormat: *serializationFormat}
		if err := conn.WriteJSON(&requestMessage); err != nil {
			log.Println(err)
			break
		}
		log.Println("reading...")
		// Read response data from server:
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("msg read: %d bytes\n", len(data))
		m.ResponseTime = time.Since(startResponseClock).Seconds() * 1000
		// Extract structuring time, serialization time, and original data size from data:
		m.OriginalDataSize = extractUint64FromBytes(data, len(data)-4-8-8-8, len(data)-4-8-8)
		m.StructuringTime = extractFloat64FromBytes(data, len(data)-4-8-8, len(data)-4-8)
		m.SerializationTime = extractFloat64FromBytes(data, len(data)-4-8, len(data)-4)
		// Extract and validate id from data:
		m.ID = extractUint32FromBytes(data, len(data)-4, len(data))
		if m.ID != requestMessage.ID {
			log.Println("ID from requestMessage doesn't match ID recieved from server")
			break
		}
		// Extract osm data from data:
		data = data[:len(data)-8-8-8-4]
		m.SerializedDataSize = len(data)
		// Deserialize data:
		startDeserializationClock := time.Now()
		switch *serializationFormat {
		case "pb":
			if err := expcodec.DeserializePB(data); err != nil {
				log.Println(err)
				break
			}
		case "fbs":
			if err := expcodec.DeserializeFbs(data); err != nil {
				log.Println(err)
				break
			}
		default:
			log.Fatalln("serialization format not supported")
		}
		m.DeserializationTime = time.Since(startDeserializationClock).Seconds() * 1000
		m.AccessTime = time.Since(startAccessClock).Seconds() * 1000
		// Log data to file after warmup iterations
		if i >= *warmups {
			m.ID = m.ID - uint32(*warmups)
			m.Log()
		}
	}
}

func extractUint64FromBytes(data []byte, start, end int) uint64 {
	return binary.LittleEndian.Uint64(data[start:end])
}

func extractUint32FromBytes(data []byte, start, end int) uint32 {
	return binary.LittleEndian.Uint32(data[start:end])
}

func extractFloat64FromBytes(data []byte, start, end int) float64 {
	return float64FromBytes(data[start:end])
}

func float64FromBytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	f := math.Float64frombits(bits)
	return f
}
