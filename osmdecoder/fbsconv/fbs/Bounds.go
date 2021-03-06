// automatically generated by the FlatBuffers compiler, do not modify

package fbs

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Bounds struct {
	_tab flatbuffers.Table
}

func GetRootAsBounds(buf []byte, offset flatbuffers.UOffsetT) *Bounds {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Bounds{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Bounds) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Bounds) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Bounds) Minlat() float64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetFloat64(o + rcv._tab.Pos)
	}
	return 0.0
}

func (rcv *Bounds) MutateMinlat(n float64) bool {
	return rcv._tab.MutateFloat64Slot(4, n)
}

func (rcv *Bounds) Minlon() float64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetFloat64(o + rcv._tab.Pos)
	}
	return 0.0
}

func (rcv *Bounds) MutateMinlon(n float64) bool {
	return rcv._tab.MutateFloat64Slot(6, n)
}

func (rcv *Bounds) Maxlat() float64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetFloat64(o + rcv._tab.Pos)
	}
	return 0.0
}

func (rcv *Bounds) MutateMaxlat(n float64) bool {
	return rcv._tab.MutateFloat64Slot(8, n)
}

func (rcv *Bounds) Maxlon() float64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetFloat64(o + rcv._tab.Pos)
	}
	return 0.0
}

func (rcv *Bounds) MutateMaxlon(n float64) bool {
	return rcv._tab.MutateFloat64Slot(10, n)
}

func BoundsStart(builder *flatbuffers.Builder) {
	builder.StartObject(4)
}
func BoundsAddMinlat(builder *flatbuffers.Builder, minlat float64) {
	builder.PrependFloat64Slot(0, minlat, 0.0)
}
func BoundsAddMinlon(builder *flatbuffers.Builder, minlon float64) {
	builder.PrependFloat64Slot(1, minlon, 0.0)
}
func BoundsAddMaxlat(builder *flatbuffers.Builder, maxlat float64) {
	builder.PrependFloat64Slot(2, maxlat, 0.0)
}
func BoundsAddMaxlon(builder *flatbuffers.Builder, maxlon float64) {
	builder.PrependFloat64Slot(3, maxlon, 0.0)
}
func BoundsEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
