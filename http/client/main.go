package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/d15johro/examensarbete-vt18/metrics"
	"github.com/d15johro/examensarbete-vt18/osmdecoder/fbsconv/fbs"
	"github.com/d15johro/examensarbete-vt18/osmdecoder/pbconv/pb"
	"github.com/golang/protobuf/proto"
	flatbuffers "github.com/google/flatbuffers/go"
)

var (
	serializationFormat = flag.String("sf", "pb", "Serialization format")
)

func init() {
	flag.Parse()
}

func main() {
	m := metrics.New()
	m.Filepath = "./http_" + *serializationFormat + ".txt"
	if err := m.Setup(); err != nil {
		log.Fatalln(err)
	}
	c := http.Client{}
	for i := 0; i < 10; i++ { // experimental
		log.Println(i)
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
		m.ID = uint32(id)
		serializationDuration, err := time.ParseDuration(resp.Header.Get("serializationDuration"))
		if err != nil {
			log.Fatalln(err)
		}
		m.SerializationTime = serializationDuration.Seconds() * 1000
		structuringDuration, err := time.ParseDuration(resp.Header.Get("structuringDuration"))
		if err != nil {
			log.Fatalln(err)
		}
		m.StructuringTime = structuringDuration.Seconds() * 1000
		// Read response data:
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		m.ResponseTime = time.Since(startResponseClock).Seconds() * 1000
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
