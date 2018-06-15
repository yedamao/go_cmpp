package protocol

// SP或ISMG请求拆除连接（CMPP_TERMINATE）操作
type Terminate struct {
	*Header
}

func NewTerminate(sequenceID uint32) (*Terminate, error) {
	op := &Terminate{}

	op.Header = &Header{}

	op.Total_Length = 4 + 4 + 4
	op.Command_Id = CMPP_TERMINATE
	op.Sequence_Id = sequenceID

	return op, nil
}

func ParseTermanite(hdr *Header, data []byte) (*Terminate, error) {
	op := &Terminate{}
	op.Header = hdr

	return op, nil
}

func (op *Terminate) Serialize() []byte {

	return op.Header.Serialize()
}

func (op *Terminate) String() string {

	return op.Header.String()
}

func (op *Terminate) Ok() error {
	return nil
}

type TerminateResp struct {
	*Header
}

func NewTerminateResp(sequenceID uint32) (*TerminateResp, error) {
	op := &TerminateResp{}

	op.Header = &Header{}

	op.Total_Length = 4 + 4 + 4
	op.Command_Id = CMPP_TERMINATE_RESP
	op.Sequence_Id = sequenceID

	return op, nil
}

func ParseTermaniteResp(hdr *Header, data []byte) (*TerminateResp, error) {
	op := &TerminateResp{}
	op.Header = hdr

	return op, nil
}

func (op *TerminateResp) Serialize() []byte {

	return op.Header.Serialize()
}

func (op *TerminateResp) String() string {

	return op.Header.String()
}

func (op *TerminateResp) Ok() error {
	return nil
}
