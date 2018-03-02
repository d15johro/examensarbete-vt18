package pbconv

import (
	"errors"
	"fmt"

	"github.com/d15johro/examensarbete-vt18/osmdecoder"
	"github.com/d15johro/examensarbete-vt18/osmdecoder/pbconv/pb"
)

// Make makes a pb.OSM out of a osmdecoder.OSM.
func Make(osm *osmdecoder.OSM) (*pb.OSM, error) {
	if osm == nil {
		return nil, fmt.Errorf("osm cannot be nil")
	}
	pbOSM := &pb.OSM{
		Version:     osm.Version,
		Generator:   osm.Generator,
		Copyright:   osm.Copyright,
		Attribution: osm.Attribution,
		License:     osm.License,
		Bounds:      makeBounds(osm.Bounds),
	}
	if nodes, err := makeNodes(osm.Nodes); err == nil {
		pbOSM.Nodes = nodes
	}
	if relations, err := makeRelations(osm.Relations); err == nil {
		pbOSM.Relations = relations
	}
	if ways, err := makeWays(osm.Ways); err == nil {
		pbOSM.Ways = ways
	}
	return pbOSM, nil
}

func makeBounds(bounds osmdecoder.Bounds) *pb.Bounds {
	return &pb.Bounds{
		Minlat: bounds.Minlat,
		Maxlat: bounds.Maxlat,
		Minlon: bounds.Minlon,
		Maxlon: bounds.Maxlon,
	}
}

func makeNodes(nodes []osmdecoder.Node) ([]*pb.Node, error) {
	if len(nodes) == 0 {
		return nil, fmt.Errorf("nodes cannot be of length 0")
	}
	pbNodes := []*pb.Node{}
	for i := 0; i < len(nodes); i++ {
		pbSharedAttributes := makeSharedAttributes(nodes[i].SharedAttributes)
		pbNode := &pb.Node{
			Lat:              nodes[i].Lat,
			Lon:              nodes[i].Lon,
			SharedAttributes: pbSharedAttributes,
		}
		if tags, err := makeTags(nodes[i].Tags); err == nil {
			pbNode.Tags = tags
		}
		pbNodes = append(pbNodes, pbNode)
	}
	return pbNodes, nil
}

func makeRelations(relations []osmdecoder.Relation) ([]*pb.Relation, error) {
	if len(relations) == 0 {
		return nil, fmt.Errorf("relations cannot be of length 0")
	}
	pbRelations := []*pb.Relation{}
	for i := 0; i < len(relations); i++ {
		pbSharedAttributes := makeSharedAttributes(relations[i].SharedAttributes)
		pbRelation := &pb.Relation{
			Visible:          relations[i].Visible,
			SharedAttributes: pbSharedAttributes,
		}
		if tags, err := makeTags(relations[i].Tags); err == nil {
			pbRelation.Tags = tags
		}
		if members, err := makeMembers(relations[i].Members); err == nil {
			pbRelation.Members = members
		}
		pbRelations = append(pbRelations, pbRelation)
	}
	return pbRelations, nil
}

func makeWays(ways []osmdecoder.Way) ([]*pb.Way, error) {
	if len(ways) == 0 {
		return nil, fmt.Errorf("ways cannot be of length 0")
	}
	pbWays := []*pb.Way{}
	for i := 0; i < len(ways); i++ {
		pbSharedAttributes := makeSharedAttributes(ways[i].SharedAttributes)
		pbWay := &pb.Way{
			SharedAttributes: pbSharedAttributes,
		}
		if tags, err := makeTags(ways[i].Tags); err == nil {
			pbWay.Tags = tags
		}
		if nds, err := makeNds(ways[i].Nds); err == nil {
			pbWay.Nds = nds
		}
		pbWays = append(pbWays, pbWay)
	}
	return pbWays, nil
}

func makeTags(tags []osmdecoder.Tag) ([]*pb.Tag, error) {
	if len(tags) == 0 {
		return nil, fmt.Errorf("tags cannot be of length 0")
	}
	pbTags := []*pb.Tag{}
	for i := 0; i < len(tags); i++ {
		pbTag := &pb.Tag{
			Key:   tags[i].Key,
			Value: tags[i].Value,
		}
		pbTags = append(pbTags, pbTag)
	}
	return pbTags, nil
}

func makeNds(nds []osmdecoder.Nd) ([]*pb.Nd, error) {
	if len(nds) == 0 {
		return nil, errors.New("nds cannot be of length 0")
	}
	pbNds := []*pb.Nd{}
	for i := 0; i < len(nds); i++ {
		pbNd := &pb.Nd{Ref: nds[i].Ref}
		pbNds = append(pbNds, pbNd)
	}
	return pbNds, nil
}

func makeMembers(members []osmdecoder.Member) ([]*pb.Member, error) {
	if len(members) == 0 {
		return nil, errors.New("members cannot be of length 0")
	}
	pbMembers := []*pb.Member{}
	for i := 0; i < len(members); i++ {
		pbMember := &pb.Member{
			Role: members[i].Role,
			Type: members[i].Type,
			Ref:  members[i].Ref,
		}
		pbMembers = append(pbMembers, pbMember)
	}
	return pbMembers, nil
}

func makeSharedAttributes(sharedAttributes osmdecoder.SharedAttributes) *pb.SharedAttributes {
	return &pb.SharedAttributes{
		Id:        sharedAttributes.ID,
		Version:   sharedAttributes.Version,
		Timestamp: sharedAttributes.Timestamp,
		Uid:       sharedAttributes.UID,
		User:      sharedAttributes.User,
		Changeset: sharedAttributes.ChangeSet,
	}
}
