// automatically generated by the FlatBuffers compiler, do not modify

package fbs

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Way struct {
	_tab flatbuffers.Table
}

func GetRootAsWay(buf []byte, offset flatbuffers.UOffsetT) *Way {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Way{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Way) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Way) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Way) SharedAttributes(obj *SharedAttributes) *SharedAttributes {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
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

func (rcv *Way) Nds(obj *Nd, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Way) NdsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Way) Tags(obj *Tag, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Way) TagsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func WayStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func WayAddSharedAttributes(builder *flatbuffers.Builder, sharedAttributes flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(sharedAttributes), 0)
}
func WayAddNds(builder *flatbuffers.Builder, nds flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(nds), 0)
}
func WayStartNdsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func WayAddTags(builder *flatbuffers.Builder, tags flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(tags), 0)
}
func WayStartTagsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func WayEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
