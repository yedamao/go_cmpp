package protocol

import (
	"testing"
)

func TestActiveTest(t *testing.T) {
	op, err := NewActiveTest(1)
	if err != nil {
		t.Error(err)
	}

	parsed, err := ParseOperation(op.Serialize())
	if err != nil {
		t.Fatal(err)
	}

	_, ok := parsed.(*ActiveTest)
	if !ok {
		t.Error("not equal")
	}
}

func TestActiveTestResp(t *testing.T) {
	op, err := NewActiveTestResp(1)
	if err != nil {
		t.Error(err)
	}

	parsed, err := ParseOperation(op.Serialize())
	if err != nil {
		t.Fatal(err)
	}

	resp := parsed.(*ActiveTestResp)

	if resp.Reserved != op.Reserved {
		t.Error("not equal")
	}
}
