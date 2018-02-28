package codec

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
)

// PB is a custom Protocol Buffers codec for the "golang.org/x/net/websocket" package.
var PB = websocket.Codec{
	Marshal: func(v interface{}) (data []byte, payloadType byte, err error) {
		msg, ok := v.(proto.Message)
		if !ok {
			err = fmt.Errorf("expected v to of type proto.Message, got %v", reflect.TypeOf(v))
			return
		}
		data, err = proto.Marshal(msg)
		payloadType = websocket.BinaryFrame
		return
	},
	Unmarshal: func(data []byte, payloadType byte, v interface{}) (err error) {
		if payloadType != websocket.BinaryFrame {
			err = fmt.Errorf("expected payloadType to be %d, got %d", websocket.BinaryFrame, payloadType)
			return
		}
		msg, ok := v.(proto.Message)
		if !ok {
			err = fmt.Errorf("expected v to of type proto.Message, got %v", reflect.TypeOf(v))
			return
		}
		err = proto.Unmarshal(data, msg)
		return
	},
}
