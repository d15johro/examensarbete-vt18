package osmdecoder

type OSM struct {
	Version     string `xml:"version,attr,omitempty"`
	Generator   string `xml:"generator,attr,omitempty"`
	Copyright   string `xml:"copyright,attr,omitempty"`
	Attribution string `xml:"attribution,attr,omitempty"`
	License     string `xml:"license,attr,omitempty"`
	Bounds      Bounds
	Nodes       []Node
	Ways        []Way
	Relations   []Relation
}

type Bounds struct {
	Minlat float64 `xml:"minlat,attr,omitempty"`
	Minlon float64 `xml:"minlon,attr,omitempty"`
	Maxlat float64 `xml:"maxlat,attr,omitempty"`
	Maxlon float64 `xml:"maxlon,attr,omitempty"`
}

type Node struct {
	sharedAttributes
	Lat  float64 `xml:"lat,attr,omitempty"`
	Lon  float64 `xml:"lon,attr,omitempty"`
	Tags []Tag   `xml:"tag,omitempty"`
}

type Relation struct {
	sharedAttributes
	Visible bool     `xml:"visible,attr,omitempty"`
	Members []Member `xml:"member,omitempty"`
	Tags    []Tag    `xml:"tag,omitempty"`
}

type Way struct {
	sharedAttributes
	Tags []Tag `xml:"tag,omitempty"`
	Nds  []Nd  `xml:"nd,omitempty"`
}

type sharedAttributes struct {
	ID        int64  `xml:"id,attr,omitempty"`
	Version   int32  `xml:"version,attr,omitempty"`
	Timestamp string `xml:"timestamp,attr,omitempty"`
	UID       int64  `xml:"uid,attr,omitempty"`
	User      string `xml:"user,attr,omitempty"`
	ChangeSet int64  `xml:"changeset,attr,omitempty"`
}

type Tag struct {
	Key   string `xml:"k,attr,omitempty"`
	Value string `xml:"v,attr,omitempty"`
}

type Nd struct {
	Ref int64 `xml:"ref,attr,omitempty"`
}

type Member struct {
	Type string `xml:"type,attr,omitempty"`
	Ref  int64  `xml:"ref,attr,omitempty"`
	Role string `xml:"role,attr,omitempty"`
}
