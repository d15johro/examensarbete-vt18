package osmdecoder_test

import (
	"bytes"
	"testing"

	"github.com/d15johro/examensarbete-vt18/osmdecoder"
)

func TestDecode(t *testing.T) {
	osm, err := osmdecoder.Decode(bytes.NewReader(data))
	if err != nil {
		t.Error("expected no error")
	}
	if osm.Attribution != "some attribution" {
		t.Errorf("expected %s, got %s", "some attribution", osm.Attribution)
	}
	if osm.License != "" {
		t.Errorf("expected %s, got %s", "", osm.License)
	}
	if len(osm.Nodes) != 3 {
		t.Errorf("expected %d, got %d", 2, len(osm.Nodes))
	}
	if len(osm.Ways) != 2 {
		t.Errorf("expected %d, got %d", 2, len(osm.Ways))
	}
	if len(osm.Relations) != 2 {
		t.Errorf("expected %d, got %d", 2, len(osm.Relations))
	}
	if len(osm.Nodes[0].Tags) != 2 {
		t.Errorf("expected %d, got %d", 2, len(osm.Nodes[0].Tags))
	}
	if osm.Nodes[0].Tags[0].Key != "ele" {
		t.Errorf("expected %s, got %s", "ele", osm.Nodes[0].Tags[0].Key)
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
