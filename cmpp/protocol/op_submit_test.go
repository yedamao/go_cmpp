package protocol

import (
	"testing"
)

func TestSubmit(t *testing.T) {
	op, err := NewSubmit(
		1,
		1, 1, 1, 1,
		"", 2, "",
		0, 0, 0,
		"", "01", "", "", "",
		"1069000000", []string{"17600000000", "16600000000"},
		[]byte("hello haha, test msg"),
	)
	if err != nil {
		t.Error(err)
	}

	parsed, err := ParseOperation(op.Serialize())
	if err != nil {
		t.Fatal(err)
	}

	submit := parsed.(*Submit)

	if submit.MsgId != op.MsgId ||
		submit.PkTotal != op.PkTotal ||
		submit.PkNumber != op.PkNumber ||
		submit.MsgFmt != op.MsgFmt ||
		submit.TP_pid != op.TP_pid ||
		submit.TP_udhi != op.TP_udhi ||
		submit.SrcId.String() != op.SrcId.String() ||
		submit.ServiceId.String() != op.ServiceId.String() ||
		submit.DestUsrTl != op.DestUsrTl ||
		submit.MsgLength != op.MsgLength {

		t.Error("not equal")
	}
}

func TestSubmitResp(t *testing.T) {
	op, err := NewSubmitResp(1, 12345, 0)
	if err != nil {
		t.Error(err)
	}

	parsed, err := ParseOperation(op.Serialize())
	if err != nil {
		t.Fatal(err)
	}

	resp := parsed.(*SubmitResp)

	if resp.MsgId != op.MsgId ||
		resp.Result != op.Result {
		t.Error("parsedLoginResp not equal")
	}
}
