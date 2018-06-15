package protocol

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

// ISMG向SP送交短信（CMPP_DELIVER）操作
//
// CMPP_DELIVER操作的目的是ISMG把从短信中心
// 或其它ISMG转发来的短信送交SP，SP以CMPP_DELIVER_RESP消息回应
type Deliver struct {
	*Header

	// Body
	MsgId         uint64       // 信息标识
	DestId        *OctetString // 目的号码
	ServiceId     *OctetString // 业务类型，是数字、字母和符号的组合
	TP_pid        uint8
	TP_udhi       uint8
	MsgFmt        uint8
	SrcTerminalId *OctetString

	// 是否为状态报告
	// 0：非状态报告
	// 1：状态报告
	RegisteredDelivery uint8
	MsgLength          uint8
	MsgContent         []byte
	Reserved           *OctetString
}

func NewDeliver(
	sequenceID uint32, msgId uint64, destId, serviceId string,
	TP_pid, TP_udhi, msgFmt uint8, srcTerminalId string,
	isReport uint8, content []byte,
) (*Deliver, error) {
	op := &Deliver{}

	op.Header = &Header{}
	var length uint32 = 4 + 4 + 4 // header length

	op.MsgId = msgId
	length = length + 8

	op.DestId = &OctetString{Data: []byte(destId), FixedLen: 21}
	length = length + 21

	op.ServiceId = &OctetString{Data: []byte(serviceId), FixedLen: 10}
	length = length + 10

	op.TP_pid = TP_pid
	length = length + 1

	op.TP_udhi = TP_udhi
	length = length + 1

	op.MsgFmt = msgFmt
	length = length + 1

	op.SrcTerminalId = &OctetString{Data: []byte(srcTerminalId), FixedLen: 21}
	length = length + 21

	op.RegisteredDelivery = isReport
	length = length + 1

	msgLen := len(content)
	op.MsgLength = uint8(msgLen)
	length = length + 1

	op.MsgContent = content
	length = length + uint32(msgLen)

	op.Reserved = &OctetString{FixedLen: 8}
	length = length + 8

	op.Total_Length = length
	op.Command_Id = CMPP_DELIVER
	op.Sequence_Id = sequenceID

	return op, nil
}

func ParseDeliver(hdr *Header, data []byte) (*Deliver, error) {
	p := 0
	op := &Deliver{}
	op.Header = hdr

	op.MsgId = unpackUi64(data[p : p+8])
	p = p + 8

	op.DestId = &OctetString{Data: data[p : p+21], FixedLen: 21}
	p = p + 21

	op.ServiceId = &OctetString{Data: data[p : p+10], FixedLen: 10}
	p = p + 10

	op.TP_pid = data[p]
	p = p + 1
	op.TP_udhi = data[p]
	p = p + 1
	op.MsgFmt = data[p]
	p = p + 1

	op.SrcTerminalId = &OctetString{Data: data[p : p+21], FixedLen: 21}
	p = p + 21

	op.RegisteredDelivery = data[p]
	p = p + 1
	op.MsgLength = data[p]
	p = p + 1

	op.MsgContent = data[p : p+int(op.MsgLength)]
	p = p + int(op.MsgLength)

	op.Reserved = &OctetString{Data: data[p : p+8], FixedLen: 8}
	p = p + 8

	return op, nil
}

func (op *Deliver) Serialize() []byte {
	b := op.Header.Serialize()

	b = append(b, packUi64(op.MsgId)...)
	b = append(b, op.DestId.Byte()...)
	b = append(b, op.ServiceId.Byte()...)

	b = append(b, packUi8(op.TP_pid)...)
	b = append(b, packUi8(op.TP_udhi)...)
	b = append(b, packUi8(op.MsgFmt)...)

	b = append(b, op.SrcTerminalId.Byte()...)
	b = append(b, packUi8(op.RegisteredDelivery)...)
	b = append(b, packUi8(op.MsgLength)...)
	b = append(b, op.MsgContent...)
	b = append(b, op.Reserved.Byte()...)

	return b
}

func (op *Deliver) String() string {
	var b bytes.Buffer
	b.WriteString(op.Header.String())

	fmt.Fprintln(&b, "--- Deliver ---")
	fmt.Fprintln(&b, "MsgID: ", op.MsgId)
	fmt.Fprintln(&b, "DestId: ", op.DestId)
	fmt.Fprintln(&b, "ServiceId: ", op.ServiceId)

	fmt.Fprintln(&b, "TP_pid: ", op.TP_pid)
	fmt.Fprintln(&b, "TP_udhi: ", op.TP_udhi)
	fmt.Fprintln(&b, "MsgFmt: ", op.MsgFmt)

	fmt.Fprintln(&b, "SrcTerminalId: ", op.SrcTerminalId)
	fmt.Fprintln(&b, "RegisteredDelivery: ", op.RegisteredDelivery)

	fmt.Fprintln(&b, "MsgLength: ", op.MsgLength)
	fmt.Fprintln(&b, "MsgContent: ", string(op.MsgContent))

	return b.String()
}

func (op *Deliver) Ok() error {
	return nil
}

type DeliverResp struct {
	*Header

	// Body
	MsgId  uint64 // 信息标识
	Result uint8  // 结果
}

func NewDeliverResp(
	sequenceID uint32, msgId uint64, result uint8) (*DeliverResp, error) {

	op := &DeliverResp{}

	op.Header = &Header{}
	var length uint32 = 4 + 4 + 4

	op.MsgId = msgId
	length = length + 8

	op.Result = result
	length = length + 1

	op.Total_Length = length
	op.Command_Id = CMPP_DELIVER_RESP
	op.Sequence_Id = sequenceID

	return op, nil
}

func ParseDeliverResp(hdr *Header, data []byte) (*DeliverResp, error) {
	op := &DeliverResp{}
	op.Header = hdr

	p := 0
	op.MsgId = unpackUi64(data[p : p+8])
	p = p + 8

	op.Result = data[p]
	p = p + 1

	return op, nil
}

func (op *DeliverResp) Serialize() []byte {
	b := op.Header.Serialize()

	b = append(b, packUi64(op.MsgId)...)
	b = append(b, packUi8(op.Result)...)

	return b
}

func (op *DeliverResp) String() string {
	var b bytes.Buffer
	b.WriteString(op.Header.String())

	fmt.Fprintln(&b, "--- DeliverResp ---")
	fmt.Fprintln(&b, "MsgID: ", op.MsgId)
	fmt.Fprintln(&b, "Result: ", op.Result)

	return b.String()
}

func (op *DeliverResp) Ok() (err error) {

	switch op.Result {
	case 0: // 正确
		err = nil
	case 1:
		err = errors.New("消息结构错")
	case 2:
		err = errors.New("命令字错")
	case 3:
		err = errors.New("消息序号重复")
	case 4:
		err = errors.New("消息长度错")
	case 5:
		err = errors.New("资费代码错")
	case 6:
		err = errors.New("超过最大信息长")
	case 7:
		err = errors.New("业务代码错")
	case 8:
		err = errors.New("流量控制错")
	default:
		err = errors.New("其他错误: " + strconv.Itoa(int(op.Result)))
	}

	return err
}
