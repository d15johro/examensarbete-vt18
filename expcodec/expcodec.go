// Package expcodec contains codec functions for experiment.
package expcodec

import (
	"fmt"

	"github.com/d15johro/examensarbete-vt18/osmdecoder"
	"github.com/d15johro/examensarbete-vt18/osmdecoder/fbsconv"
	"github.com/d15johro/examensarbete-vt18/osmdecoder/fbsconv/fbs"
	"github.com/d15johro/examensarbete-vt18/osmdecoder/pbconv/pb"
	"github.com/golang/protobuf/proto"
	flatbuffers "github.com/google/flatbuffers/go"
)

func SerializePB(x proto.Message) ([]byte, error) {
	return proto.Marshal(x)
}

// unlike pb, fbs serialize when building the buffer. You don't have to structure the data in a pb.OSM equivalent
// struct before serializing, the buffer ends up in a serialized form.
func SerializeFBS(builder *flatbuffers.Builder, x *osmdecoder.OSM) ([]byte, error) {
	builder.Reset()
	err := fbsconv.Build(builder, x)
	if err != nil {
		return nil, err
	}
	return builder.Bytes[builder.Head():], nil
}

// Deserialize functions are not returning any pointers to structs (*fbs.OSM or *pb.OSM) since we are not doing anything with the data,
// we are only measuring how long it takes to deserialize.

func DeserializePB(data []byte) error {
	osm := &pb.OSM{}
	return proto.Unmarshal(data, osm)
}

// Unlike pb, deserializing fbs basically means storing the raw binary data in a struct.
// Therefore, the deserialization time for fbs will probably always be 0ms. Since maps
// want access to all data we measure deserialization for fbs as time it takes to access all data
// from memory. This means that we have to call all the getter methods for all elements in the map.
func DeserializeFbs(data []byte) error {
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
