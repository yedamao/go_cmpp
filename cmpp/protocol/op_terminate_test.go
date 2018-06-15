package protocol

import (
	"testing"
)

func TestTerminate(t *testing.T) {
	op, err := NewTerminate(1)
	if err != nil {
		t.Error(err)
	}

	_, err = ParseOperation(op.Serialize())
	if err != nil {
		t.Fatal(err)
	}

}

func TestTerminateResp(t *testing.T) {
	op, err := NewTerminateResp(1)
	if err != nil {
		t.Error(err)
	}

	_, err = ParseOperation(op.Serialize())
	if err != nil {
		t.Fatal(err)
	}

}
