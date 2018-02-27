package tests

import (
	"testing"

	"golang.org/x/net/websocket"

	"github.com/d15johro/examensarbete-vt18/websocket/codec"
)

var testy = Test{
	Text:   "some test",
	Number: 42,
}

func TestPBCodecMarshal(t *testing.T) {
	_, payloadType, err := codec.PB.Marshal(&testy)
	if err != nil {
		t.Error("expected no error:", err)
	}
	if payloadType != websocket.BinaryFrame {
		t.Errorf("expected payload to be %d, got %d", payloadType, websocket.BinaryFrame)
	}
	_, _, err = codec.PB.Marshal(struct{ s string }{"wrong type"})
	if err == nil {
		t.Error("expected error")
	}
}

func TestPBCodecUnmarshal(t *testing.T) {
	data, _, _ := codec.PB.Marshal(&testy)
	var x Test
	err := codec.PB.Unmarshal(data, websocket.BinaryFrame, &x)
	if err != nil {
		t.Error("expected no error:", err)
	}
	if x.Text != "some test" {
		t.Errorf("expected %s, got %s", "some test", x.Text)
	}
	if x.Number != 42 {
		t.Errorf("expected %d, got %d", 42, x.Number)
	}
	err = codec.PB.Unmarshal(data, websocket.TextFrame, &x)
	if err == nil {
		t.Error("expected error")
	}
}
