package osmdecoder

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

// DecodeFile decodes a .osm file given path.
func DecodeFile(path string) (*OSM, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("DecodeFile: %s", err.Error())
	}
	defer file.Close()
	return Decode(file)
}

// Decode decodes a stream of .osm data into the OSM struct.
func Decode(reader io.Reader) (*OSM, error) {
	osm := &OSM{}
	decoder := xml.NewDecoder(reader)
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch typedToken := token.(type) {
		case xml.StartElement:
			switch typedToken.Name.Local {
			case "osm":
				attrs := typedToken.Attr
				for i := 0; i < len(attrs); i++ {
					attr := attrs[i]
					switch attr.Name.Local {
					case "version":
						osm.Version = attr.Value
					case "generator":
						osm.Generator = attr.Value
					case "copyright":
						osm.Copyright = attr.Value
					case "attribution":
						osm.Attribution = attr.Value
					case "license":
						osm.License = attr.Value
					}
				}
			case "node":
				var node Node
				err := decoder.DecodeElement(&node, &typedToken)
				if err != nil {
					return nil, err
				}
				osm.Nodes = append(osm.Nodes, node)

			case "way":
				var way Way
				err := decoder.DecodeElement(&way, &typedToken)
				if err != nil {
					return nil, err
				}
				osm.Ways = append(osm.Ways, way)

			case "relation":
				var relation Relation
				err := decoder.DecodeElement(&relation, &typedToken)
				if err != nil {
					return nil, err
				}
				osm.Relations = append(osm.Relations, relation)
			}
		}
	}
	return osm, nil
}
