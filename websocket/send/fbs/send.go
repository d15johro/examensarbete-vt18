package fbs_send

import (
	"github.com/d15johro/examensarbete-vt18/websocket/send/fbs/fbs"
	flatbuffers "github.com/google/flatbuffers/go"
)

// BuildSend builds the buffer for a Send message.
func BuildSend(builder *flatbuffers.Builder, data string) {
	builder.Reset()
	dataOffset := builder.CreateString(data)
	fbs.SendStart(builder)
	fbs.SendAddData(builder, dataOffset)
	send := fbs.SendEnd(builder)
	builder.Finish(send)
}
