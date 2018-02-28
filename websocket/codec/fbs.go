package codec

import (
	"fmt"
	"reflect"

	flatbuffers "github.com/google/flatbuffers/go"
	"golang.org/x/net/websocket"
)

// FBS is a custom FlatBuffers codec for the "golang.org/x/net/websocket" package.
var FBS = websocket.Codec{
	Marshal: func(v interface{}) (data []byte, payloadType byte, err error) {
		b, ok := v.(*flatbuffers.Builder)
		if !ok {
			err = fmt.Errorf("expected v to of type *flatbuffers.Builder, got %v", reflect.TypeOf(v))
			return
		}
		data, err = flatbuffers.FlatbuffersCodec{}.Marshal(b)
		payloadType = websocket.BinaryFrame
		return
	},
	Unmarshal: func(data []byte, payloadType byte, v interface{}) (err error) {
		if payloadType != websocket.BinaryFrame {
			err = fmt.Errorf("expected payloadType to be %d, got %d", websocket.BinaryFrame, payloadType)
			return
		}
		x, ok := v.(flatbuffers.FlatBuffer)
		if !ok {
			err = fmt.Errorf("expected v to of type flatbuffers.flatbuffersInit, got %v", reflect.TypeOf(v))
			return
		}
		err = flatbuffers.FlatbuffersCodec{}.Unmarshal(data, x)
		return
	},
}
