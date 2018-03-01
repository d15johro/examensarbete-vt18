// automatically generated by the FlatBuffers compiler, do not modify

package fbs

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Node struct {
	_tab flatbuffers.Table
}

func GetRootAsNode(buf []byte, offset flatbuffers.UOffsetT) *Node {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Node{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Node) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Node) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Node) Lat() float64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetFloat64(o + rcv._tab.Pos)
	}
	return 0.0
}

func (rcv *Node) MutateLat(n float64) bool {
	return rcv._tab.MutateFloat64Slot(4, n)
}

func (rcv *Node) Lon() float64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetFloat64(o + rcv._tab.Pos)
	}
	return 0.0
}

func (rcv *Node) MutateLon(n float64) bool {
	return rcv._tab.MutateFloat64Slot(6, n)
}

func (rcv *Node) SharedAttributes(obj *SharedAttributes) *SharedAttributes {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(SharedAttributes)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *Node) Tags(obj *Tag, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Node) TagsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func NodeStart(builder *flatbuffers.Builder) {
	builder.StartObject(4)
}
func NodeAddLat(builder *flatbuffers.Builder, lat float64) {
	builder.PrependFloat64Slot(0, lat, 0.0)
}
func NodeAddLon(builder *flatbuffers.Builder, lon float64) {
	builder.PrependFloat64Slot(1, lon, 0.0)
}
func NodeAddSharedAttributes(builder *flatbuffers.Builder, sharedAttributes flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(sharedAttributes), 0)
}
func NodeAddTags(builder *flatbuffers.Builder, tags flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(tags), 0)
}
func NodeStartTagsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func NodeEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
