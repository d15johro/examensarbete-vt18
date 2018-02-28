package tests

import (
	"testing"

	"github.com/d15johro/examensarbete-vt18/websocket/codec"
	"github.com/d15johro/examensarbete-vt18/websocket/codec/fbs_tests/fbs"
	"github.com/google/flatbuffers/go"
	"golang.org/x/net/websocket"
)

func buildTesty(builder *flatbuffers.Builder) {
	builder.Reset()
	text := builder.CreateString("some test")
	fbs.TestStart(builder)
	fbs.TestAddText(builder, text)
	fbs.TestAddNumber(builder, 42)
	testy := fbs.TestEnd(builder)
	builder.Finish(testy)
}

func TestFBSCodecMarshal(t *testing.T) {
	builder := flatbuffers.NewBuilder(0)
	buildTesty(builder)
	_, payloadType, err := codec.FBS.Marshal(builder)
	if err != nil {
		t.Errorf("expected no error:", err)
	}
	if payloadType != websocket.BinaryFrame {
		t.Errorf("expected payload to be %d, got %d", payloadType, websocket.BinaryFrame)
	}
	_, _, err = codec.FBS.Marshal("wrong type")
	if err == nil {
		t.Error("expected error")
	}
}

func TestFBSCodecUnmarshal(t *testing.T) {
	builder := flatbuffers.NewBuilder(0)
	buildTesty(builder)
	data, _, _ := codec.FBS.Marshal(builder)

	var testy fbs.Test
	err := codec.FBS.Unmarshal(data, websocket.BinaryFrame, &testy)
	if err != nil {
		t.Error("expected no error:", err)
	}
	if string(testy.Text()) != "some test" {
		t.Errorf("expected %s, got %s", "some test", testy.Text())
	}
	if testy.Number() != 42 {
		t.Errorf("expected %d, got %d", 42, testy.Number())
	}
	err = codec.FBS.Unmarshal(data, websocket.BinaryFrame, "wong input")
	if err == nil {
		t.Error("expected error")
	}
}
