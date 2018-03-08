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
	nFiles              = flag.Int("nf", 10, "# of files")
)

type metrics struct {
	id                  uint32
	accessTime          float64
	responseTime        float64
	serializationTime   float64
	deserializationTime float64
	structuringTime     float64
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
	for i := 0; i < (*nFiles)*(*nFiles); i++ { // experimental
		log.Println(i)
		// Request data from server:
		startAccessClock := time.Now()
		startResponseClock := time.Now()
		requestMessage := struct {
			ID                  uint32 `json:"id"`
			SerializationFormat string `json:"serializationFormat"`
			NumberOfFiles       uint32 `json:"numberOfFiles"`
		}{ID: uint32(i), SerializationFormat: *serializationFormat, NumberOfFiles: uint32(*nFiles)}
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
		// Extract structuring and serialization time from data:
		m.structuringTime = extractFloat64FromBytes(data, len(data)-4-8-8, len(data)-4-8)
		m.serializationTime = extractFloat64FromBytes(data, len(data)-4-8, len(data)-4)
		// Extract and validate id from data:
		m.id = extractUint32FromBytes(data, len(data)-4, len(data))
		if m.id != requestMessage.ID {
			log.Println("ID from requestMessage doesn't match ID recieved from server")
			break
		}
		// Extract osm data from data:
		data = data[:len(data)-8-8-4]
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
			if err := deserializeFbs(data); err != nil {
				log.Println(err)
				break
			}
		default:
			log.Fatalln("serialization format not supported")
		}
		m.deserializationTime = time.Since(startDeserializationClock).Seconds() * 1000
		m.accessTime = time.Since(startAccessClock).Seconds() * 1000

		m.log()

	}
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

func (m *metrics) log() {
	s := fmt.Sprintf("%d,%f,%f,%f,%f,%f,%d\n",
		m.id, m.accessTime,
		m.responseTime,
		m.serializationTime,
		m.deserializationTime,
		m.structuringTime,
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

// Unlike pb, deserializing fbs basically means storing the raw binary data in a struct.
// Therefore, the deserialization time for fbs will probably always be 0ms. Since maps
// want access to all data we measure deserialization for fbs as time it takes to access all data
// from memory. This means that we have to call all the getter methods for all elements in the map.
func deserializeFbs(data []byte) error {
	// osm:
	offset := flatbuffers.UOffsetT(0)
	n := flatbuffers.GetUOffsetT(data[offset:])
	osm := &fbs.OSM{}
	osm.Init(data, n+offset)
	// osm attributes:
	osm.Attribution()
	osm.Copyright()
	osm.Generator()
	osm.License()
	osm.Version()
	// nodes:
	for i := 0; i < osm.NodesLength(); i++ {
		// node:
		var node fbs.Node
		ok := osm.Nodes(&node, i)
		if !ok {
			return fmt.Errorf("bad deserialization")
		}
		// node attributes:
		node.Lat()
		node.Lon()
		// shared attributes:
		var sharedAttributes fbs.SharedAttributes
		sA := node.SharedAttributes(&sharedAttributes)
		sA.Version()
		sA.User()
		sA.Uid()
		sA.Timestamp()
		sA.Id()
		sA.Changeset()
		// tags:
		for j := 0; j < node.TagsLength(); j++ {
			// tag:
			var tag fbs.Tag
			ok := node.Tags(&tag, j)
			if !ok {
				return fmt.Errorf("bad deserialization")
			}
			//tag attributes:
			tag.Key()
			tag.Value()
		}
	}
	// ways:
	for i := 0; i < osm.WaysLength(); i++ {
		// way:
		var way fbs.Way
		ok := osm.Ways(&way, i)
		if !ok {
			return fmt.Errorf("bad deserialization")
		}
		// shared attributes:
		var sharedAttributes fbs.SharedAttributes
		sA := way.SharedAttributes(&sharedAttributes)
		sA.Version()
		sA.User()
		sA.Uid()
		sA.Timestamp()
		sA.Id()
		sA.Changeset()
		// nds:
		for j := 0; j < way.NdsLength(); j++ {
			// nd:
			var nd fbs.Nd
			ok := way.Nds(&nd, j)
			if !ok {
				return fmt.Errorf("bad deserialization")
			}
			// nd attributes:
			nd.Ref()
		}
		// tags:
		for j := 0; j < way.TagsLength(); j++ {
			// tag:
			var tag fbs.Tag
			ok := way.Tags(&tag, j)
			if !ok {
				return fmt.Errorf("bad deserialization")
			}
			//tag attributes:
			tag.Key()
			tag.Value()
		}
	}
	// relations
	for i := 0; i < osm.RelationsLength(); i++ {
		// relation:
		var relation fbs.Relation
		ok := osm.Relations(&relation, i)
		if !ok {
			return fmt.Errorf("bad deserialization")
		}
		// relation attributes:
		relation.Visible()
		// shared attributes:
		var sharedAttributes fbs.SharedAttributes
		sA := relation.SharedAttributes(&sharedAttributes)
		sA.Version()
		sA.User()
		sA.Uid()
		sA.Timestamp()
		sA.Id()
		sA.Changeset()
		// members:
		for j := 0; j < relation.MembersLength(); j++ {
			// member:
			var member fbs.Member
			ok := relation.Members(&member, j)
			if !ok {
				return fmt.Errorf("bad deserialization")
			}
			// member attributes:
			member.Ref()
			member.Role()
			member.Type()
		}
		// tags:
		for j := 0; j < relation.TagsLength(); j++ {
			// tag:
			var tag fbs.Tag
			ok := relation.Tags(&tag, j)
			if !ok {
				return fmt.Errorf("bad deserialization")
			}
			//tag attributes:
			tag.Key()
			tag.Value()
		}
	}
	return nil
}
