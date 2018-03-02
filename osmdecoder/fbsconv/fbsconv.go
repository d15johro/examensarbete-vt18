package fbsconv

import (
	"errors"

	"github.com/d15johro/examensarbete-vt18/osmdecoder/fbsconv/fbs"

	"github.com/d15johro/examensarbete-vt18/osmdecoder"
	"github.com/google/flatbuffers/go"
)

// Build builds a fbs.OSM out of a osmdecoder.OSM.
func Build(builder *flatbuffers.Builder, osm *osmdecoder.OSM) error {
	if osm == nil {
		return errors.New("osm cannot be nil")
	}
	builder.Reset()
	root := buildOSM(builder, osm)
	builder.Finish(root)
	return nil
}

func buildOSM(builder *flatbuffers.Builder, osm *osmdecoder.OSM) flatbuffers.UOffsetT {
	version := builder.CreateString(osm.Version)
	generator := builder.CreateString(osm.Generator)
	copyright := builder.CreateString(osm.Copyright)
	attribution := builder.CreateString(osm.Attribution)
	license := builder.CreateString(osm.License)
	bounds := buildBounds(builder, osm.Bounds)
	// nodes:
	nodesUOffsetTs := buildNodesUOffsetTs(builder, osm.Nodes)
	fbs.OSMStartNodesVector(builder, len(osm.Nodes))
	for i := len(osm.Nodes) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(nodesUOffsetTs[i])
	}
	nodes := builder.EndVector(len(osm.Nodes))
	// ways:
	waysUOffsetTs := buildWaysUOffsetTs(builder, osm.Ways)
	fbs.OSMStartNodesVector(builder, len(osm.Ways))
	for i := len(osm.Ways) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(waysUOffsetTs[i])
	}
	ways := builder.EndVector(len(osm.Ways))
	// relations:
	relationsUOffsetTs := buildRelationsUOffsetTs(builder, osm.Relations)
	fbs.OSMStartNodesVector(builder, len(osm.Relations))
	for i := len(osm.Relations) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(relationsUOffsetTs[i])
	}
	relations := builder.EndVector(len(osm.Relations))
	// OSM:
	fbs.OSMStart(builder)
	fbs.OSMAddVersion(builder, version)
	fbs.OSMAddGenerator(builder, generator)
	fbs.OSMAddAttribution(builder, attribution)
	fbs.OSMAddCopyright(builder, copyright)
	fbs.OSMAddLicense(builder, license)
	fbs.OSMAddBounds(builder, bounds)
	fbs.OSMAddNodes(builder, nodes)
	fbs.OSMAddWays(builder, ways)
	fbs.OSMAddRelations(builder, relations)
	return fbs.OSMEnd(builder)
}

func buildBounds(builder *flatbuffers.Builder, bounds osmdecoder.Bounds) flatbuffers.UOffsetT {
	fbs.BoundsStart(builder)
	fbs.BoundsAddMaxlat(builder, bounds.Maxlat)
	fbs.BoundsAddMaxlon(builder, bounds.Maxlon)
	fbs.BoundsAddMinlat(builder, bounds.Minlat)
	fbs.BoundsAddMinlon(builder, bounds.Minlon)
	return fbs.BoundsEnd(builder)
}

func buildNodesUOffsetTs(builder *flatbuffers.Builder, nodes []osmdecoder.Node) []flatbuffers.UOffsetT {
	if len(nodes) == 0 {
		return nil
	}
	fbsNodes := make([]flatbuffers.UOffsetT, len(nodes))
	for i := 0; i < len(nodes); i++ {
		sharedAttributes := buildSharedAttributes(builder, osmdecoder.SharedAttributes{
			ID:        nodes[i].ID,
			User:      nodes[i].User,
			UID:       nodes[i].UID,
			ChangeSet: nodes[i].ChangeSet,
			Version:   nodes[i].Version,
			Timestamp: nodes[i].Timestamp,
		})
		// tags:
		tagsUOffsetTs := buildTagsUOffsetTs(builder, nodes[i].Tags)
		fbs.NodeStartTagsVector(builder, len(nodes[i].Tags))
		for i := len(nodes[i].Tags) - 1; i >= 0; i-- {
			builder.PrependUOffsetT(tagsUOffsetTs[i])
		}
		tags := builder.EndVector(len(nodes[i].Tags))
		// node:
		fbs.NodeStart(builder)
		fbs.NodeAddLat(builder, nodes[i].Lat)
		fbs.NodeAddLon(builder, nodes[i].Lon)
		fbs.NodeAddSharedAttributes(builder, sharedAttributes)
		fbs.NodeAddTags(builder, tags)
		fbsNodes[i] = fbs.NodeEnd(builder)
	}
	return fbsNodes
}

func buildWaysUOffsetTs(builder *flatbuffers.Builder, ways []osmdecoder.Way) []flatbuffers.UOffsetT {
	if len(ways) == 0 {
		return nil
	}
	fbsWays := make([]flatbuffers.UOffsetT, len(ways))
	for i := 0; i < len(ways); i++ {
		sharedAttributes := buildSharedAttributes(builder, osmdecoder.SharedAttributes{
			ID:        ways[i].ID,
			User:      ways[i].User,
			UID:       ways[i].UID,
			ChangeSet: ways[i].ChangeSet,
			Version:   ways[i].Version,
			Timestamp: ways[i].Timestamp,
		})
		// nds:
		ndsUOffsetTs := buildNdsUOffsetTs(builder, ways[i].Nds)
		fbs.WayStartNdsVector(builder, len(ways[i].Nds))
		for i := len(ways[i].Nds) - 1; i >= 0; i-- {
			builder.PrependUOffsetT(ndsUOffsetTs[i])
		}
		nds := builder.EndVector(len(ways[i].Nds))
		// tags:
		tagsUOffsetTs := buildTagsUOffsetTs(builder, ways[i].Tags)
		fbs.WayStartTagsVector(builder, len(ways[i].Tags))
		for i := len(ways[i].Tags) - 1; i >= 0; i-- {
			builder.PrependUOffsetT(tagsUOffsetTs[i])
		}
		tags := builder.EndVector(len(ways[i].Tags))
		// way:
		fbs.WayStart(builder)
		fbs.WayAddSharedAttributes(builder, sharedAttributes)
		fbs.WayAddNds(builder, nds)
		fbs.WayAddTags(builder, tags)
		fbsWays[i] = fbs.WayEnd(builder)
	}
	return fbsWays
}

func buildRelationsUOffsetTs(builder *flatbuffers.Builder, relations []osmdecoder.Relation) []flatbuffers.UOffsetT {
	if len(relations) == 0 {
		return nil
	}
	fbsRelation := make([]flatbuffers.UOffsetT, len(relations))
	for i := 0; i < len(relations); i++ {

		visible := byte(0)
		if relations[i].Visible {
			visible = byte(1)
		}
		sharedAttributes := buildSharedAttributes(builder, osmdecoder.SharedAttributes{
			ID:        relations[i].ID,
			User:      relations[i].User,
			UID:       relations[i].UID,
			ChangeSet: relations[i].ChangeSet,
			Version:   relations[i].Version,
			Timestamp: relations[i].Timestamp,
		})
		// members:
		membersUOffsetTs := buildMembersUOffsetTs(builder, relations[i].Members)
		fbs.RelationStartMembersVector(builder, len(relations[i].Members))
		for i := len(relations[i].Members) - 1; i >= 0; i-- {
			builder.PrependUOffsetT(membersUOffsetTs[i])
		}
		members := builder.EndVector(len(relations[i].Members))
		// tags:
		tagsUOffsetTs := buildTagsUOffsetTs(builder, relations[i].Tags)
		fbs.RelationStartTagsVector(builder, len(relations[i].Tags))
		for i := len(relations[i].Tags) - 1; i >= 0; i-- {
			builder.PrependUOffsetT(tagsUOffsetTs[i])
		}
		tags := builder.EndVector(len(relations[i].Tags))
		// relation:
		fbs.RelationStart(builder)
		fbs.RelationAddVisible(builder, visible)
		fbs.RelationAddSharedAttributes(builder, sharedAttributes)
		fbs.RelationAddMembers(builder, members)
		fbs.RelationAddTags(builder, tags)
		fbsRelation[i] = fbs.RelationEnd(builder)
	}
	return fbsRelation
}

func buildSharedAttributes(builder *flatbuffers.Builder, sharedAttributes osmdecoder.SharedAttributes) flatbuffers.UOffsetT {
	timestamp := builder.CreateString(sharedAttributes.Timestamp)
	user := builder.CreateString(sharedAttributes.User)
	fbs.SharedAttributesStart(builder)
	fbs.SharedAttributesAddTimestamp(builder, timestamp)
	fbs.SharedAttributesAddUser(builder, user)
	fbs.SharedAttributesAddId(builder, sharedAttributes.ID)
	fbs.SharedAttributesAddUid(builder, sharedAttributes.UID)
	fbs.SharedAttributesAddChangeset(builder, sharedAttributes.ChangeSet)
	fbs.SharedAttributesAddVersion(builder, sharedAttributes.Version)
	return fbs.SharedAttributesEnd(builder)
}

func buildMembersUOffsetTs(builder *flatbuffers.Builder, members []osmdecoder.Member) []flatbuffers.UOffsetT {
	if len(members) == 0 {
		return nil
	}
	fbsMembers := make([]flatbuffers.UOffsetT, len(members))
	for i := 0; i < len(members); i++ {
		role := builder.CreateString(members[i].Role)
		mType := builder.CreateString(members[i].Type)
		fbs.MemberStart(builder)
		fbs.MemberAddRef(builder, members[i].Ref)
		fbs.MemberAddRole(builder, role)
		fbs.MemberAddType(builder, mType)
		fbsMembers[i] = fbs.MemberEnd(builder)
	}
	return fbsMembers
}

func buildNdsUOffsetTs(builder *flatbuffers.Builder, nds []osmdecoder.Nd) []flatbuffers.UOffsetT {
	if len(nds) == 0 {
		return nil
	}
	fbsNds := make([]flatbuffers.UOffsetT, len(nds))
	for i := 0; i < len(nds); i++ {
		fbs.NdStart(builder)
		fbs.NdAddRef(builder, nds[i].Ref)
		fbsNds[i] = fbs.NdEnd(builder)
	}
	return fbsNds
}

func buildTagsUOffsetTs(builder *flatbuffers.Builder, tags []osmdecoder.Tag) []flatbuffers.UOffsetT {
	if len(tags) == 0 {
		return nil
	}
	fbsTags := make([]flatbuffers.UOffsetT, len(tags))
	for i := 0; i < len(tags); i++ {
		key := builder.CreateString(tags[i].Key)
		value := builder.CreateString(tags[i].Value)
		fbs.TagStart(builder)
		fbs.TagAddKey(builder, key)
		fbs.TagAddValue(builder, value)
		fbsTags[i] = fbs.TagEnd(builder)
	}
	return fbsTags
}
