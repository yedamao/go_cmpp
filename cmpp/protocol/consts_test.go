package protocol

import (
	"testing"
)

func TestConsts(t *testing.T) {
	if CMPP_SUBMIT != 0x00000004 ||
		CMPP_SUBMIT_RESP != 0x80000004 {

		t.Error("const Command_Id error")
	}
}
