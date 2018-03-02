// Code generated by protoc-gen-go. DO NOT EDIT.
// source: osm.proto

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	osm.proto

It has these top-level messages:
	OSM
	Bounds
	SharedAttributes
	Node
	Way
	Nd
	Relation
	Tag
	Member
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type OSM struct {
	Version     string      `protobuf:"bytes,2,opt,name=version" json:"version,omitempty"`
	Generator   string      `protobuf:"bytes,3,opt,name=generator" json:"generator,omitempty"`
	Copyright   string      `protobuf:"bytes,4,opt,name=copyright" json:"copyright,omitempty"`
	Attribution string      `protobuf:"bytes,5,opt,name=attribution" json:"attribution,omitempty"`
	License     string      `protobuf:"bytes,6,opt,name=license" json:"license,omitempty"`
	Bounds      *Bounds     `protobuf:"bytes,7,opt,name=bounds" json:"bounds,omitempty"`
	Nodes       []*Node     `protobuf:"bytes,8,rep,name=nodes" json:"nodes,omitempty"`
	Ways        []*Way      `protobuf:"bytes,9,rep,name=ways" json:"ways,omitempty"`
	Relations   []*Relation `protobuf:"bytes,10,rep,name=relations" json:"relations,omitempty"`
}

func (m *OSM) Reset()                    { *m = OSM{} }
func (m *OSM) String() string            { return proto.CompactTextString(m) }
func (*OSM) ProtoMessage()               {}
func (*OSM) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *OSM) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *OSM) GetGenerator() string {
	if m != nil {
		return m.Generator
	}
	return ""
}

func (m *OSM) GetCopyright() string {
	if m != nil {
		return m.Copyright
	}
	return ""
}

func (m *OSM) GetAttribution() string {
	if m != nil {
		return m.Attribution
	}
	return ""
}

func (m *OSM) GetLicense() string {
	if m != nil {
		return m.License
	}
	return ""
}

func (m *OSM) GetBounds() *Bounds {
	if m != nil {
		return m.Bounds
	}
	return nil
}

func (m *OSM) GetNodes() []*Node {
	if m != nil {
		return m.Nodes
	}
	return nil
}

func (m *OSM) GetWays() []*Way {
	if m != nil {
		return m.Ways
	}
	return nil
}

func (m *OSM) GetRelations() []*Relation {
	if m != nil {
		return m.Relations
	}
	return nil
}

type Bounds struct {
	Minlat float64 `protobuf:"fixed64,1,opt,name=minlat" json:"minlat,omitempty"`
	Minlon float64 `protobuf:"fixed64,2,opt,name=minlon" json:"minlon,omitempty"`
	Maxlat float64 `protobuf:"fixed64,3,opt,name=maxlat" json:"maxlat,omitempty"`
	Maxlon float64 `protobuf:"fixed64,4,opt,name=maxlon" json:"maxlon,omitempty"`
}

func (m *Bounds) Reset()                    { *m = Bounds{} }
func (m *Bounds) String() string            { return proto.CompactTextString(m) }
func (*Bounds) ProtoMessage()               {}
func (*Bounds) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Bounds) GetMinlat() float64 {
	if m != nil {
		return m.Minlat
	}
	return 0
}

func (m *Bounds) GetMinlon() float64 {
	if m != nil {
		return m.Minlon
	}
	return 0
}

func (m *Bounds) GetMaxlat() float64 {
	if m != nil {
		return m.Maxlat
	}
	return 0
}

func (m *Bounds) GetMaxlon() float64 {
	if m != nil {
		return m.Maxlon
	}
	return 0
}

type SharedAttributes struct {
	Id        int64  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Version   int32  `protobuf:"varint,2,opt,name=version" json:"version,omitempty"`
	Timestamp string `protobuf:"bytes,3,opt,name=timestamp" json:"timestamp,omitempty"`
	Uid       int64  `protobuf:"varint,4,opt,name=uid" json:"uid,omitempty"`
	User      string `protobuf:"bytes,5,opt,name=user" json:"user,omitempty"`
	Changeset int64  `protobuf:"varint,6,opt,name=changeset" json:"changeset,omitempty"`
}

func (m *SharedAttributes) Reset()                    { *m = SharedAttributes{} }
func (m *SharedAttributes) String() string            { return proto.CompactTextString(m) }
func (*SharedAttributes) ProtoMessage()               {}
func (*SharedAttributes) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *SharedAttributes) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *SharedAttributes) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *SharedAttributes) GetTimestamp() string {
	if m != nil {
		return m.Timestamp
	}
	return ""
}

func (m *SharedAttributes) GetUid() int64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *SharedAttributes) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *SharedAttributes) GetChangeset() int64 {
	if m != nil {
		return m.Changeset
	}
	return 0
}

type Node struct {
	Lat              float64           `protobuf:"fixed64,1,opt,name=lat" json:"lat,omitempty"`
	Lon              float64           `protobuf:"fixed64,2,opt,name=lon" json:"lon,omitempty"`
	SharedAttributes *SharedAttributes `protobuf:"bytes,3,opt,name=sharedAttributes" json:"sharedAttributes,omitempty"`
	Tags             []*Tag            `protobuf:"bytes,4,rep,name=tags" json:"tags,omitempty"`
}

func (m *Node) Reset()                    { *m = Node{} }
func (m *Node) String() string            { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()               {}
func (*Node) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Node) GetLat() float64 {
	if m != nil {
		return m.Lat
	}
	return 0
}

func (m *Node) GetLon() float64 {
	if m != nil {
		return m.Lon
	}
	return 0
}

func (m *Node) GetSharedAttributes() *SharedAttributes {
	if m != nil {
		return m.SharedAttributes
	}
	return nil
}

func (m *Node) GetTags() []*Tag {
	if m != nil {
		return m.Tags
	}
	return nil
}

type Way struct {
	SharedAttributes *SharedAttributes `protobuf:"bytes,1,opt,name=sharedAttributes" json:"sharedAttributes,omitempty"`
	Tags             []*Tag            `protobuf:"bytes,2,rep,name=tags" json:"tags,omitempty"`
	Nds              []*Nd             `protobuf:"bytes,3,rep,name=nds" json:"nds,omitempty"`
}

func (m *Way) Reset()                    { *m = Way{} }
func (m *Way) String() string            { return proto.CompactTextString(m) }
func (*Way) ProtoMessage()               {}
func (*Way) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Way) GetSharedAttributes() *SharedAttributes {
	if m != nil {
		return m.SharedAttributes
	}
	return nil
}

func (m *Way) GetTags() []*Tag {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *Way) GetNds() []*Nd {
	if m != nil {
		return m.Nds
	}
	return nil
}

type Nd struct {
	Ref int64 `protobuf:"varint,1,opt,name=ref" json:"ref,omitempty"`
}

func (m *Nd) Reset()                    { *m = Nd{} }
func (m *Nd) String() string            { return proto.CompactTextString(m) }
func (*Nd) ProtoMessage()               {}
func (*Nd) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Nd) GetRef() int64 {
	if m != nil {
		return m.Ref
	}
	return 0
}

type Relation struct {
	Visible          bool              `protobuf:"varint,1,opt,name=visible" json:"visible,omitempty"`
	SharedAttributes *SharedAttributes `protobuf:"bytes,2,opt,name=sharedAttributes" json:"sharedAttributes,omitempty"`
	Members          []*Member         `protobuf:"bytes,3,rep,name=members" json:"members,omitempty"`
	Tags             []*Tag            `protobuf:"bytes,4,rep,name=tags" json:"tags,omitempty"`
}

func (m *Relation) Reset()                    { *m = Relation{} }
func (m *Relation) String() string            { return proto.CompactTextString(m) }
func (*Relation) ProtoMessage()               {}
func (*Relation) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *Relation) GetVisible() bool {
	if m != nil {
		return m.Visible
	}
	return false
}

func (m *Relation) GetSharedAttributes() *SharedAttributes {
	if m != nil {
		return m.SharedAttributes
	}
	return nil
}

func (m *Relation) GetMembers() []*Member {
	if m != nil {
		return m.Members
	}
	return nil
}

func (m *Relation) GetTags() []*Tag {
	if m != nil {
		return m.Tags
	}
	return nil
}

type Tag struct {
	Key   string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *Tag) Reset()                    { *m = Tag{} }
func (m *Tag) String() string            { return proto.CompactTextString(m) }
func (*Tag) ProtoMessage()               {}
func (*Tag) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *Tag) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Tag) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type Member struct {
	Type string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	Ref  int64  `protobuf:"varint,2,opt,name=ref" json:"ref,omitempty"`
	Role string `protobuf:"bytes,3,opt,name=role" json:"role,omitempty"`
}

func (m *Member) Reset()                    { *m = Member{} }
func (m *Member) String() string            { return proto.CompactTextString(m) }
func (*Member) ProtoMessage()               {}
func (*Member) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *Member) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Member) GetRef() int64 {
	if m != nil {
		return m.Ref
	}
	return 0
}

func (m *Member) GetRole() string {
	if m != nil {
		return m.Role
	}
	return ""
}

func init() {
	proto.RegisterType((*OSM)(nil), "pb.OSM")
	proto.RegisterType((*Bounds)(nil), "pb.Bounds")
	proto.RegisterType((*SharedAttributes)(nil), "pb.SharedAttributes")
	proto.RegisterType((*Node)(nil), "pb.Node")
	proto.RegisterType((*Way)(nil), "pb.Way")
	proto.RegisterType((*Nd)(nil), "pb.Nd")
	proto.RegisterType((*Relation)(nil), "pb.Relation")
	proto.RegisterType((*Tag)(nil), "pb.Tag")
	proto.RegisterType((*Member)(nil), "pb.Member")
}

func init() { proto.RegisterFile("osm.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 532 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0xc1, 0x8e, 0xd3, 0x30,
	0x10, 0x86, 0x95, 0xb8, 0x4d, 0x9b, 0x29, 0x42, 0x95, 0xb5, 0x5a, 0x59, 0x02, 0xa1, 0x2a, 0xe2,
	0x50, 0x21, 0xd1, 0x43, 0x79, 0x01, 0xd8, 0xfb, 0x2e, 0x92, 0x77, 0xa5, 0x3d, 0x3b, 0x1b, 0x93,
	0x46, 0x24, 0x71, 0x64, 0xbb, 0x0b, 0xb9, 0x70, 0xe7, 0x25, 0x38, 0x71, 0xe1, 0x2d, 0xd1, 0x4c,
	0x9c, 0x16, 0x4a, 0x85, 0x04, 0xb7, 0x99, 0x6f, 0xa6, 0x9e, 0xf1, 0xff, 0xbb, 0x81, 0xd4, 0xb8,
	0x66, 0xd3, 0x59, 0xe3, 0x0d, 0x8f, 0xbb, 0x3c, 0xfb, 0x1e, 0x03, 0x7b, 0x7f, 0x7b, 0xcd, 0x05,
	0xcc, 0x1e, 0xb5, 0x75, 0x95, 0x69, 0x45, 0xbc, 0x8a, 0xd6, 0xa9, 0x1c, 0x53, 0xfe, 0x1c, 0xd2,
	0x52, 0xb7, 0xda, 0x2a, 0x6f, 0xac, 0x60, 0x54, 0x3b, 0x02, 0xac, 0x3e, 0x98, 0xae, 0xb7, 0x55,
	0xb9, 0xf3, 0x62, 0x32, 0x54, 0x0f, 0x80, 0xaf, 0x60, 0xa1, 0xbc, 0xb7, 0x55, 0xbe, 0xf7, 0x78,
	0xf2, 0x94, 0xea, 0xbf, 0x22, 0x9c, 0x5b, 0x57, 0x0f, 0xba, 0x75, 0x5a, 0x24, 0xc3, 0xdc, 0x90,
	0xf2, 0x0c, 0x92, 0xdc, 0xec, 0xdb, 0xc2, 0x89, 0xd9, 0x2a, 0x5a, 0x2f, 0xb6, 0xb0, 0xe9, 0xf2,
	0xcd, 0x15, 0x11, 0x19, 0x2a, 0xfc, 0x05, 0x4c, 0x5b, 0x53, 0x68, 0x27, 0xe6, 0x2b, 0xb6, 0x5e,
	0x6c, 0xe7, 0xd8, 0x72, 0x63, 0x0a, 0x2d, 0x07, 0xcc, 0x9f, 0xc1, 0xe4, 0x93, 0xea, 0x9d, 0x48,
	0xa9, 0x3c, 0xc3, 0xf2, 0xbd, 0xea, 0x25, 0x41, 0xfe, 0x0a, 0x52, 0xab, 0x6b, 0x85, 0x6b, 0x38,
	0x01, 0xd4, 0xf1, 0x04, 0x3b, 0x64, 0x80, 0xf2, 0x58, 0xce, 0x76, 0x90, 0x0c, 0xa3, 0xf9, 0x25,
	0x24, 0x4d, 0xd5, 0xd6, 0xca, 0x8b, 0x68, 0x15, 0xad, 0x23, 0x19, 0xb2, 0x91, 0x07, 0xfd, 0x02,
	0x37, 0x2d, 0x71, 0xf5, 0x19, 0xfb, 0x59, 0xe0, 0x94, 0x8d, 0xdc, 0xb4, 0xa4, 0x5a, 0xe0, 0xa6,
	0xcd, 0xbe, 0x45, 0xb0, 0xbc, 0xdd, 0x29, 0xab, 0x8b, 0x77, 0x41, 0x26, 0xed, 0xf8, 0x53, 0x88,
	0xab, 0x82, 0x06, 0x32, 0x19, 0x57, 0xc5, 0xa9, 0x5b, 0xd3, 0xdf, 0xdc, 0xf2, 0x55, 0xa3, 0x9d,
	0x57, 0x4d, 0x37, 0xba, 0x75, 0x00, 0x7c, 0x09, 0x6c, 0x5f, 0x15, 0x34, 0x91, 0x49, 0x0c, 0x39,
	0x87, 0xc9, 0xde, 0x69, 0x1b, 0xac, 0xa1, 0x98, 0x3c, 0xdd, 0xa9, 0xb6, 0xd4, 0x4e, 0x7b, 0x72,
	0x85, 0xc9, 0x23, 0xc8, 0xbe, 0x46, 0x30, 0x41, 0x8d, 0xf1, 0xb0, 0xa3, 0x0c, 0x18, 0x12, 0x39,
	0x08, 0x80, 0x21, 0x7f, 0x0b, 0x4b, 0x77, 0x72, 0x19, 0xda, 0x6a, 0xb1, 0xbd, 0x40, 0xa9, 0x4f,
	0x2f, 0x2a, 0xff, 0xe8, 0x46, 0x0b, 0xbd, 0x2a, 0x9d, 0x98, 0x1c, 0x2d, 0xbc, 0x53, 0xa5, 0x24,
	0x98, 0x7d, 0x01, 0x76, 0xaf, 0xfa, 0xb3, 0x53, 0xa2, 0xff, 0x9a, 0x12, 0x9f, 0x99, 0xc2, 0x05,
	0x30, 0x7c, 0x86, 0x8c, 0x6a, 0x09, 0xbd, 0xb1, 0x42, 0x22, 0xca, 0x2e, 0x21, 0xbe, 0x29, 0xf0,
	0xda, 0x56, 0x7f, 0x08, 0xf6, 0x60, 0x98, 0xfd, 0x88, 0x60, 0x3e, 0x3e, 0x23, 0x32, 0xab, 0x72,
	0x55, 0x5e, 0x6b, 0x6a, 0x99, 0xcb, 0x31, 0x3d, 0xbb, 0x77, 0xfc, 0x4f, 0x7b, 0xbf, 0x84, 0x59,
	0xa3, 0x9b, 0x5c, 0xdb, 0x71, 0x3d, 0xfa, 0x97, 0x5c, 0x13, 0x92, 0x63, 0xe9, 0xef, 0x1a, 0xbe,
	0x06, 0x76, 0xa7, 0x4a, 0xbc, 0xc4, 0x47, 0xdd, 0xd3, 0x86, 0xa9, 0xc4, 0x90, 0x5f, 0xc0, 0xf4,
	0x51, 0xd5, 0x7b, 0x1d, 0x3e, 0x08, 0x43, 0x92, 0x5d, 0x41, 0x32, 0x1c, 0x8f, 0x4f, 0xc7, 0xf7,
	0x9d, 0x0e, 0x3f, 0xa1, 0x78, 0x94, 0x22, 0x3e, 0x48, 0x81, 0x5d, 0xd6, 0xd4, 0x3a, 0xbc, 0x45,
	0x8a, 0xf3, 0x84, 0xbe, 0x3f, 0x6f, 0x7e, 0x06, 0x00, 0x00, 0xff, 0xff, 0xe5, 0x71, 0x4a, 0xa6,
	0x8c, 0x04, 0x00, 0x00,
}
