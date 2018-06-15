package protocol

import (
	"bytes"
	"errors"
	"fmt"
)

//	SP向ISMG发起删除短信（CMPP_CANCEL）操作
//
// CMPP_CANCEL操作的目的是SP通过此操作可以将
// 已经提交给ISMG的短信删除，ISMG将以CMPP_CANCEL_RESP
// 回应删除操作的结果。
type Cancel struct {
	*Header

	// Body
	MsgId uint64 // 信息标识（SP想要删除的信息标识）
}

func NewCancel(sequenceID uint32, msgId uint64) (*Cancel, error) {
	op := &Cancel{}

	op.Header = &Header{}
	var length uint32 = 4 + 4 + 4 // header length

	op.MsgId = msgId
	length = length + 8

	op.Total_Length = length
	op.Command_Id = CMPP_CANCEL
	op.Sequence_Id = sequenceID

	return op, nil
}

func ParseCancel(hdr *Header, data []byte) (*Cancel, error) {
	p := 0

	op := &Cancel{}
	op.Header = hdr

	op.MsgId = unpackUi64(data[p : p+8])
	p = p + 8

	return op, nil
}

func (op *Cancel) Serialize() []byte {
	b := op.Header.Serialize()

	b = append(b, packUi64(op.MsgId)...)

	return b
}

func (op *Cancel) String() string {
	var b bytes.Buffer
	b.WriteString(op.Header.String())

	fmt.Fprintln(&b, "--- Cancel ---")
	fmt.Fprintln(&b, "Msg_Id: ", op.MsgId)

	return b.String()
}

func (op *Cancel) Ok() error {
	return nil
}

type CancelResp struct {
	*Header

	// Body

	// 成功标识
	// 0：成功
	// 1：失败
	SuccessId uint8
}

func NewCancelResp(sequenceID uint32, successId uint8) (*CancelResp, error) {
	op := &CancelResp{}

	op.Header = &Header{}
	var length uint32 = 4 + 4 + 4

	op.SuccessId = successId
	length = length + 1

	op.Total_Length = length
	op.Command_Id = CMPP_CANCEL_RESP
	op.Sequence_Id = sequenceID

	return op, nil
}

func ParseCancelResp(hdr *Header, data []byte) (*CancelResp, error) {
	op := &CancelResp{}
	op.Header = hdr

	p := 0
	op.SuccessId = data[p]
	p = p + 1

	return op, nil
}

func (op *CancelResp) Serialize() []byte {
	b := op.Header.Serialize()

	b = append(b, packUi8(op.SuccessId)...)

	return b
}

func (op *CancelResp) String() string {
	var b bytes.Buffer
	b.WriteString(op.Header.String())

	fmt.Fprintln(&b, "--- CancelResp ---")
	fmt.Fprintln(&b, "SuccessId: ", op.SuccessId)

	return b.String()
}

func (op *CancelResp) Ok() (err error) {

	if 0 == op.SuccessId {
		err = nil
	} else {
		err = errors.New("失败")
	}

	return err
}
