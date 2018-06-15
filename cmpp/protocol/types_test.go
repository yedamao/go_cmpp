package protocol

import (
	"testing"
)

func TestOctetString(t *testing.T) {

	// serialize
	str_content := "1069009010053"
	o := &OctetString{
		Data:     []byte(str_content),
		FixedLen: 21,
	}

	if len(o.Byte()) != o.FixedLen {
		t.Error("len error")
	}

	if str_content != o.String() {
		t.Error("OctetString string error")
	}

	// decode
	bytes_content := []byte{49, 48, 54, 57, 48, 48, 57, 48, 49, 48, 48, 53, 51, 0, 0, 0, 0, 0, 0, 0, 0}

	o = &OctetString{
		Data:     bytes_content,
		FixedLen: 21,
	}

	if len(o.Byte()) != o.FixedLen {
		t.Error("len error")
	}

	if str_content != o.String() {
		t.Error("OctetString string error")
	}
}

func TestUi64(t *testing.T) {
	var x uint64 = 123456789

	if x != unpackUi64(packUi64(x)) {
		t.Error("uint64 pack unpack not equal")
	}
}
