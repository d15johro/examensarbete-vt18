package codec_test

import (
	"bytes"
	"testing"

	"github.com/d15johro/examensarbete-vt18/osmdecoder/pbconv"

	"github.com/d15johro/examensarbete-vt18/osmdecoder"
	"github.com/d15johro/examensarbete-vt18/osmdecoder/fbsconv"
	"github.com/d15johro/examensarbete-vt18/osmdecoder/fbsconv/fbs"
	"github.com/d15johro/examensarbete-vt18/osmdecoder/pbconv/pb"
	"github.com/d15johro/examensarbete-vt18/websocket/codec"
	"golang.org/x/net/websocket"
)

func TestPBCodecMarshal(t *testing.T) {
	osm, _ := osmdecoder.Decode(bytes.NewReader(data))
	pbOSM, _ := pbconv.Make(osm)
	_, payloadType, err := codec.PB.Marshal(pbOSM)
	if err != nil {
		t.Error("expected no error:", err)
	}
	if payloadType != websocket.BinaryFrame {
		t.Errorf("expected payload to be %d, got %d", payloadType, websocket.BinaryFrame)
	}
	_, _, err = codec.PB.Marshal("bad input")
	if err == nil {
		t.Error("expected error")
	}
}

func TestPBCodecUnmarshal(t *testing.T) {
	osm, _ := osmdecoder.Decode(bytes.NewReader(data))
	x, _ := pbconv.Make(osm)
	data, _, _ := codec.PB.Marshal(x)
	var pbOSM pb.OSM
	err := codec.PB.Unmarshal(data, websocket.BinaryFrame, &pbOSM)
	if err != nil {
		t.Error("expected no error:", err)
	}
	if pbOSM.Attribution != "some attribution" {
		t.Errorf("expected %s, got %s", "some attribution", pbOSM.Attribution)
	}
	err = codec.FBS.Unmarshal(data, websocket.BinaryFrame, "bad input")
	if err == nil {
		t.Error("expected error")
	}
}

func TestFBSCodecMarshal(t *testing.T) {
	osm, _ := osmdecoder.Decode(bytes.NewReader(data))
	builder, _ := fbsconv.Make(osm)
	data, payloadType, err := codec.FBS.Marshal(builder)
	if err != nil {
		t.Errorf("expected no error:", err)
	}
	if payloadType != websocket.BinaryFrame {
		t.Errorf("expected payload to be %d, got %d", payloadType, websocket.BinaryFrame)
	}
	if len(data) == 0 {
		t.Error("expected length of data to be greater than 0")
	}
	_, _, err = codec.FBS.Marshal("bad type")
	if err == nil {
		t.Error("expected error")
	}
}

func TestFBSCodecUnmarshal(t *testing.T) {
	osm, _ := osmdecoder.Decode(bytes.NewReader(data))
	builder, _ := fbsconv.Make(osm)
	data, _, _ := codec.FBS.Marshal(builder)

	var fbsOSM fbs.OSM
	err := codec.FBS.Unmarshal(data, websocket.BinaryFrame, &fbsOSM)
	if err != nil {
		t.Error("expected no error:", err)
	}
	if string(fbsOSM.Attribution()) != "some attribution" {
		t.Errorf("expected %s, got %s", "some attribution", string(fbsOSM.Attribution()))
	}
	err = codec.FBS.Unmarshal(data, websocket.BinaryFrame, "bad input")
	if err == nil {
		t.Error("expected error")
	}
}

var data = []byte(`
	<?xml version="1.0" encoding="UTF-8"?>
	<osm version="0.6" generator="" copyright="some copyright" attribution="some attribution">
		<bounds minlat="39.0149000" minlon="-105.8183000" maxlat="39.0825000" maxlon="-105.6836000"/>
		<node id="151399886" visible="true" version="2" changeset="145015" timestamp="2008-18-06T13:43:46Z" user="johndoe" uid="42" lat="39.0519362" lon="-105.7205623">
			<tag k="ele" v="2738"/>
			<tag k="gnis:Class" v="Populated Place"/>
		</node>
		<node id="177751860" visible="true" version="4" changeset="27496079" timestamp="2018-12-15T23:40:01Z" user="johndoe" uid="42" lat="39.0930075" lon="-105.7188191"/>
		<node id="177751888" visible="true" version="4" changeset="27642703" timestamp="2018-12-23T02:43:00Z" user="johndoe" uid="42" lat="39.0966927" lon="-105.7169881"/>
		<way id="17116248" visible="true" version="4" changeset="27900617" timestamp="2018-01-04T00:50:21Z" user="johndoe" uid="42">
			<nd ref="2654521390"/>
			<nd ref="177754587"/>
			<tag k="tiger:cfcc" v="A41"/>
			<tag k="tiger:county" v="Park, CO"/>
		</way>
		<way id="17119329" visible="true" version="1" changeset="354864" timestamp="2008-12-19T01:03:23Z" user="johndoe" uid="42">
			<nd ref="177784222"/>
		</way>
		<relation id="83671" visible="true" version="109" changeset="54497178" timestamp="2018-12-09T21:42:37Z" user="johndoe" uid="42">
			<member type="way" ref="470381973" role=""/>
			<member type="way" ref="17050765" role=""/>
			<member type="way" ref="310298778" role=""/>
		</relation>
		<relation id="83671" visible="true" version="109" changeset="54497178" timestamp="2018-12-09T21:42:37Z" user="johndoe" uid="42">
			<tag k="cycle_network" v="US:ACA"/>
			<tag k="is_in:state" v="CO"/>
		</relation>
	</osm>
`)
