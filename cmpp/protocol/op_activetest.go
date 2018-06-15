package protocol

import (
	"bytes"
	"fmt"
)

//	链路检测（CMPP_ACTIVE_TEST）操作
type ActiveTest struct {
	*Header

	// Body
}

func NewActiveTest(sequenceID uint32) (*ActiveTest, error) {
	op := &ActiveTest{}

	op.Header = &Header{}

	op.Total_Length = 4 + 4 + 4
	op.Command_Id = CMPP_ACTIVE_TEST
	op.Sequence_Id = sequenceID

	return op, nil
}

func ParseActiveTest(hdr *Header, data []byte) (*ActiveTest, error) {
	op := &ActiveTest{}
	op.Header = hdr

	return op, nil
}

func (op *ActiveTest) Serialize() []byte {

	return op.Header.Serialize()
}

func (op *ActiveTest) String() string {

	return op.Header.String()
}

func (op *ActiveTest) Ok() error {
	return nil
}

type ActiveTestResp struct {
	*Header

	// Body
	Reserved uint8
}

func NewActiveTestResp(sequenceID uint32) (*ActiveTestResp, error) {
	op := &ActiveTestResp{}

	op.Header = &Header{}
	var length uint32 = 4 + 4 + 4

	op.Reserved = 0
	length = length + 1

	op.Total_Length = length
	op.Command_Id = CMPP_ACTIVE_TEST_RESP
	op.Sequence_Id = sequenceID

	return op, nil
}

func ParseActiveTestResp(hdr *Header, data []byte) (*ActiveTestResp, error) {
	op := &ActiveTestResp{}
	op.Header = hdr

	p := 0
	op.Reserved = data[p]
	p = p + 1

	return op, nil
}

func (op *ActiveTestResp) Serialize() []byte {
	b := op.Header.Serialize()

	b = append(b, packUi8(op.Reserved)...)

	return b
}

func (op *ActiveTestResp) String() string {
	var b bytes.Buffer
	b.WriteString(op.Header.String())

	fmt.Fprintln(&b, "--- ActiveTestResp ---")
	fmt.Fprintln(&b, "Reserved: ", op.Reserved)

	return b.String()
}

func (op *ActiveTestResp) Ok() error {
	return nil
}
