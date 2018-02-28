// automatically generated by the FlatBuffers compiler, do not modify

package fbs

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Test struct {
	_tab flatbuffers.Table
}

func GetRootAsTest(buf []byte, offset flatbuffers.UOffsetT) *Test {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Test{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Test) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Test) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Test) Text() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Test) Number() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Test) MutateNumber(n int32) bool {
	return rcv._tab.MutateInt32Slot(6, n)
}

func TestStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func TestAddText(builder *flatbuffers.Builder, text flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(text), 0)
}
func TestAddNumber(builder *flatbuffers.Builder, number int32) {
	builder.PrependInt32Slot(1, number, 0)
}
func TestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
