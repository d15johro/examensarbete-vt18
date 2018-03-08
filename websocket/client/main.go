package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/d15johro/examensarbete-vt18/metrics"
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
	for i := 0; i < 10; i++ { // experimental
		log.Println(i)
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
		m.ResponseTime = time.Since(startResponseClock).Seconds() * 1000
		// Extract structuring and serialization time from data:
		m.StructuringTime = extractFloat64FromBytes(data, len(data)-4-8-8, len(data)-4-8)
		m.SerializationTime = extractFloat64FromBytes(data, len(data)-4-8, len(data)-4)
		// Extract and validate id from data:
		m.ID = extractUint32FromBytes(data, len(data)-4, len(data))
		if m.ID != requestMessage.ID {
			log.Println("ID from requestMessage doesn't match ID recieved from server")
			break
		}
		// Extract osm data from data:
		data = data[:len(data)-8-8-4]
		m.DataSize = len(data)
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
		m.DeserializationTime = time.Since(startDeserializationClock).Seconds() * 1000
		m.AccessTime = time.Since(startAccessClock).Seconds() * 1000

		m.Log()

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
