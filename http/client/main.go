package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/d15johro/examensarbete-vt18/osmdecoder/fbsconv/fbs"
	"github.com/d15johro/examensarbete-vt18/osmdecoder/pbconv/pb"
	"github.com/golang/protobuf/proto"
	flatbuffers "github.com/google/flatbuffers/go"
)

var serializationFormat = flag.String("sf", "pb", "Serialization format")

type metrics struct {
	id                  int
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
	m := metrics{filepath: "./" + *serializationFormat + ".txt"}
	if err := m.setup(); err != nil {
		log.Fatalln(err)
	}
	c := http.Client{}
	for i := 0; i < 13*10; i++ {
		startAccessClock := time.Now()
		startResponseClock := time.Now()
		// Send GET request to server:
		url := fmt.Sprintf("http://localhost:8080/%d", i)
		resp, err := c.Get(url)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		// Validate response:
		if resp.StatusCode != http.StatusOK {
			log.Fatalln(http.StatusText(resp.StatusCode))
		}
		// Collect metrics from response header:
		id, err := strconv.Atoi(resp.Header.Get("id"))
		if err != nil {
			log.Fatalln(err)
		}
		m.id = id
		serializationDuration, err := time.ParseDuration(resp.Header.Get("serializationDuration"))
		if err != nil {
			log.Fatalln(err)
		}
		m.serializationTime = serializationDuration.Seconds() * 1000
		// Read response data:
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		m.responseTime = time.Since(startResponseClock).Seconds() * 1000
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
		m.log()
		// Print collected metrics to consol:
		log.Printf("%+v\n", m)
	}
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
