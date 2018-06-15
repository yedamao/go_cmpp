package protocol

import (
	"testing"
)

func TestConnect(t *testing.T) {
	op, err := NewConnect(1, "", "")
	if err != nil {
		t.Error(err)
	}

	parsed, err := ParseOperation(op.Serialize())
	if err != nil {
		t.Fatal(err)
	}

	parsedConnect := parsed.(*Connect)

	if parsedConnect.SourceAddr.String() != op.SourceAddr.String() ||
		parsedConnect.AuthenticatorSource.String() != op.AuthenticatorSource.String() ||
		parsedConnect.Timestamp != op.Timestamp ||
		parsedConnect.Version != op.Version {
		t.Error("parsedLogin not equal")
	}
}

func TestConnectResp(t *testing.T) {
	op, err := NewConnectResp(1, 0, "")
	if err != nil {
		t.Error(err)
	}

	parsed, err := ParseOperation(op.Serialize())
	if err != nil {
		t.Fatal(err)
	}

	parsedConnectResp := parsed.(*ConnectResp)

	if parsedConnectResp.Status != op.Status ||
		parsedConnectResp.AuthenticatorISMG.String() != op.AuthenticatorISMG.String() ||
		parsedConnectResp.Version != op.Version {
		t.Error("parsedLoginResp not equal")
	}
}
