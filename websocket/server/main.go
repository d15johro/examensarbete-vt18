package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/d15johro/examensarbete-vt18/expcodec"
	"github.com/d15johro/examensarbete-vt18/fs"
	"github.com/d15johro/examensarbete-vt18/osmdecoder"
	"github.com/d15johro/examensarbete-vt18/osmdecoder/pbconv"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/gorilla/websocket"
)

var (
	addr          = flag.String("addr", "localhost:8080", "http server address")
	mapsDir       = "../../data/maps/"
	numberOfFiles uint32
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 10000,
}

func init() {
	flag.Parse()
	var err error
	numberOfFiles, err = fs.FileCount(mapsDir)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	log.Println(numberOfFiles)
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
		// Read client request:
		var requestMessage struct {
			ID                  uint32 `json:"id"`
			SerializationFormat string `json:"serializationFormat"`
		}
		if err := conn.ReadJSON(&requestMessage); err != nil {
			log.Println(err)
			break
		}
		// Decode .osm file depending on id from request message:
		filename := mapsDir + "map" + fmt.Sprintf("%d", requestMessage.ID%numberOfFiles) + ".osm"
		file, err := os.Open(filename)
		if err != nil {
			log.Println(err)
			break
		}
		x, err := osmdecoder.Decode(file)
		if err != nil {
			log.Println(err)
			break
		}
		file.Close()
		// get original file size:
		fileinfo, err := os.Stat(filename)
		if err != nil {
			log.Println(err)
			break
		}
		originalDataSize := uint64(fileinfo.Size())
		// Serialize object:
		var (
			data                    []byte
			startSerializationClock time.Time
			startStructuringClock   time.Time
			structuringTime         float64
		)
		switch requestMessage.SerializationFormat {
		case "pb":
			startStructuringClock = time.Now()
			osm, err := pbconv.Make(x)
			if err != nil {
				log.Println(err)
				break
			}
			structuringTime = time.Since(startStructuringClock).Seconds() * 1000
			startSerializationClock = time.Now()
			data, err = expcodec.SerializePB(osm)
			if err != nil {
				log.Println(err)
				break
			}
		case "fbs":
			// Since flatbuffers already stores data in a "serialized" form, serialization basically
			// means getting a pointer to the internal storage. Therefore, unlike pb where we start the
			// serialize clock after the pb.OSM object has been structured, full build/serialize cycle
			// is measured. This is done in osmcodec.SerializeFBS(x *osmdecoder.OSM).
			startSerializationClock = time.Now()
			builder := flatbuffers.NewBuilder(0)
			data, err = expcodec.SerializeFBS(builder, x)
			if err != nil {
				log.Println(err)
				break
			}
			// Structuring time is not meaused in fbs since its the same as serialization time. We send
			// the default value of time.Duration back to the client.
		default:
			log.Fatalln("serialization format not supported")
		}
		serializationTime := time.Since(startSerializationClock).Seconds() * 1000

		// Send data to client:
		data = appendUint64ToBytes(data, originalDataSize)
		data = appendFloat64ToBytes(data, structuringTime)
		data = appendFloat64ToBytes(data, serializationTime)
		data = appendUint32ToBytes(data, requestMessage.ID)
		if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
			log.Println(err)
			break
		}
	}
}

func uint64ToBytes(ui uint64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, ui)
	return bytes
}

func appendUint64ToBytes(data []byte, ui uint64) []byte {
	buf := uint64ToBytes(ui)
	for i := 0; i < len(buf); i++ {
		data = append(data, buf[i])
	}
	return data
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
