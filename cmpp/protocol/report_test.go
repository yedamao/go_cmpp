package protocol

import (
	"testing"
)

func TestParseReport(t *testing.T) {
	rawData := []byte{102, 66, 144, 128, 9, 217, 89, 254, 77, 75, 58, 48, 48, 49, 50, 49, 56, 48, 54, 49, 50, 49, 54, 52, 48, 49, 56, 48, 54, 49, 50, 49, 54, 52, 48, 56, 54, 49, 52, 55, 49, 52, 55, 54, 50, 56, 57, 52, 0, 0, 0, 0, 0, 0, 0, 0, 63, 138, 46, 192}

	// test parse raw
	_, err := ParseReport(rawData)
	if err != nil {
		t.Error(err)
	}

	// test parse new
	rpt, err := NewReport(1234, "DELIVRD", "", "", "17600000000", 0)
	if err != nil {
		t.Error(err)
	}

	parsed, err := ParseReport(rpt.Serialize())
	if err != nil {
		t.Error(err)
	}

	if rpt.MsgId != parsed.MsgId ||
		rpt.Stat.String() != parsed.Stat.String() ||
		rpt.DestTerminalId.String() != parsed.DestTerminalId.String() {
		t.Error("parsed report not match")
	}
}
