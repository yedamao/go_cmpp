package protocol

import (
	"testing"
)

func TestDeliver(t *testing.T) {
	op, err := NewDeliver(
		1, 12345, "1069000000", "", 0, 0, 0, "16600000000",
		0, []byte("hello haha, test msg"),
	)
	if err != nil {
		t.Error(err)
	}

	parsed, err := ParseOperation(op.Serialize())
	if err != nil {
		t.Fatal(err)
	}

	deliver := parsed.(*Deliver)

	if deliver.MsgId != op.MsgId ||
		deliver.DestId.String() != op.DestId.String() ||
		deliver.ServiceId.String() != op.ServiceId.String() ||
		deliver.TP_pid != op.TP_pid ||
		deliver.TP_udhi != op.TP_udhi ||
		deliver.MsgFmt != op.MsgFmt ||
		deliver.SrcTerminalId.String() != op.SrcTerminalId.String() ||
		deliver.RegisteredDelivery != op.RegisteredDelivery {

		t.Error("parsedLogin not equal")
	}
}

func TestDeliverResp(t *testing.T) {
	op, err := NewDeliverResp(1, 12345, 0)
	if err != nil {
		t.Error(err)
	}

	parsed, err := ParseOperation(op.Serialize())
	if err != nil {
		t.Fatal(err)
	}

	resp := parsed.(*DeliverResp)

	if resp.MsgId != op.MsgId ||
		resp.Result != op.Result {
		t.Error("parsedLoginResp not equal")
	}
}
