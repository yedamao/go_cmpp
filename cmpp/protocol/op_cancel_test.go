package protocol

import (
	"testing"
)

func TestCancel(t *testing.T) {
	op, err := NewCancel(1, 123456789)
	if err != nil {
		t.Error(err)
	}

	parsed, err := ParseOperation(op.Serialize())
	if err != nil {
		t.Fatal(err)
	}

	cancel := parsed.(*Cancel)

	if cancel.MsgId != op.MsgId {
		t.Error("not equal")
	}
}

func TestCancelResp(t *testing.T) {
	op, err := NewCancelResp(1, 0)
	if err != nil {
		t.Error(err)
	}

	parsed, err := ParseOperation(op.Serialize())
	if err != nil {
		t.Fatal(err)
	}

	cancelResp := parsed.(*CancelResp)

	if cancelResp.SuccessId != op.SuccessId {
		t.Error("not equal")
	}
}
