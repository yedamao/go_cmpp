package protocol

import (
	"testing"
)

func TestParseHeader(t *testing.T) {

	h := &Header{
		Total_Length: 99,
		Command_Id:   CMPP_CONNECT,
		Sequence_Id:  1,
	}

	parsed, _ := ParseHeader(h.Serialize())

	if h.Total_Length != parsed.Total_Length ||
		h.Command_Id != parsed.Command_Id ||
		h.Sequence_Id != parsed.Sequence_Id {
		t.Error("header parsed not equal")
	}
}
