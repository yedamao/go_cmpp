package protocol

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

// SP向ISMG提交短信（CMPP_SUBMIT）操作
//
// CMPP_SUBMIT操作的目的是SP在与ISMG建立
// 应用层连接后向ISMG提交短信。
// ISMG以CMPP_SUBMIT_RESP消息响应
type Submit struct {
	*Header

	// Body
	// 信息标识，由SP侧短信网关本身产生，本处填空。
	MsgId uint64
	// 相同Msg_Id的信息总条数，从1开始
	PkTotal uint8
	// 相同Msg_Id的信息序号，从1开始
	PkNumber uint8
	// 是否要求返回状态确认报告. 0：不需要 1：需要
	RegisteredDelivery uint8
	// 信息级别
	MsgLevel uint8
	// 业务类型，是数字、字母和符号的组合
	ServiceId *OctetString
	// 计费用户类型字段
	// 0：对目的终端MSISDN计费；
	// 1：对源终端MSISDN计费；
	// 2：对SP计费;
	// 3：表示本字段无效，对谁计费参见Fee_terminal_Id字段
	FeeUserType uint8
	// 被计费用户的号码（如本字节填空，则表示本字段无效，
	// 对谁计费参见Fee_UserType字段，本字段与Fee_UserType字段互斥）
	FeeTerminalId *OctetString
	TP_pid        uint8
	TP_udhi       uint8
	MsgFmt        uint8
	// 信息内容来源(SP_Id)
	MsgSrc *OctetString
	// 资费类别
	// 01：对“计费用户号码”免费
	// 02：对“计费用户号码”按条计信息费
	// 03：对“计费用户号码”按包月收取信息费
	// 04：对“计费用户号码”的信息费封顶
	// 05：对“计费用户号码”的收费是由SP实现
	FeeType *OctetString
	// 资费代码（以分为单位）
	FeeCode *OctetString
	// 存活有效期，格式遵循SMPP3.3协议
	ValidTime *OctetString
	// 定时发送时间，格式遵循SMPP3.3协议
	AtTime *OctetString
	// 源号码
	SrcId *OctetString
	// 接收信息的用户数量(小于100个用户)
	DestUsrTl uint8
	// 接收短信的MSISDN号码
	DestTerminalId []*OctetString
	MsgLength      uint8
	// 信息长度(Msg_Fmt值为0时：<160个字节；其它<=140个字节)
	MsgContent []byte
	Reserve    *OctetString
}

func NewSubmit(
	sequenceID uint32,
	pkTotal, pkNumber, needReport, msgLevel uint8,
	serviceId string, feeUserType uint8, feeTerminalId string,
	TP_pid, TP_udhi, msgFmt uint8,
	msgSrc, feeType, feeCode, validTime, atTime, srcId string,
	destTermId []string, content []byte,
) (*Submit, error) {

	op := &Submit{}

	op.Header = &Header{}
	var length uint32 = 4 + 4 + 4 // header length

	op.MsgId = 0 // 本处填空
	length = length + 8

	op.PkTotal = pkTotal
	length = length + 1
	op.PkNumber = pkNumber
	length = length + 1
	op.RegisteredDelivery = needReport
	length = length + 1
	op.MsgLevel = msgLevel
	length = length + 1

	op.ServiceId = &OctetString{Data: []byte(serviceId), FixedLen: 10}
	length = length + 10

	op.FeeUserType = feeUserType
	length = length + 1

	op.FeeTerminalId = &OctetString{Data: []byte(feeTerminalId), FixedLen: 21}
	length = length + 21

	op.TP_pid = TP_pid
	length = length + 1
	op.TP_udhi = TP_udhi
	length = length + 1
	op.MsgFmt = msgFmt
	length = length + 1

	op.MsgSrc = &OctetString{Data: []byte(msgSrc), FixedLen: 6}
	length = length + 6

	op.FeeType = &OctetString{Data: []byte(feeType), FixedLen: 2}
	length = length + 2
	op.FeeCode = &OctetString{Data: []byte(feeCode), FixedLen: 6}
	length = length + 6
	op.ValidTime = &OctetString{Data: []byte(validTime), FixedLen: 17}
	length = length + 17
	op.AtTime = &OctetString{Data: []byte(atTime), FixedLen: 17}
	length = length + 17

	op.SrcId = &OctetString{Data: []byte(srcId), FixedLen: 21}
	length = length + 21

	// 短消息接收号码总数
	count := len(destTermId)
	if count > 100 {
		return nil, errors.New("too many destTermId")
	}
	op.DestUsrTl = uint8(count)
	length = length + 1

	for _, v := range destTermId {
		op.DestTerminalId = append(
			op.DestTerminalId, &OctetString{Data: []byte(v), FixedLen: 21})
		length = length + 21
	}

	msgLen := len(content)
	op.MsgLength = uint8(msgLen)
	length = length + 1

	op.MsgContent = content
	length = length + uint32(msgLen)

	op.Reserve = &OctetString{FixedLen: 8}
	length = length + 8

	op.Total_Length = length
	op.Command_Id = CMPP_SUBMIT
	op.Sequence_Id = sequenceID

	return op, nil
}

func ParseSubmit(hdr *Header, data []byte) (*Submit, error) {
	p := 0
	op := &Submit{}
	op.Header = hdr

	op.MsgId = unpackUi64(data[p : p+8])
	p = p + 8

	op.PkTotal = data[p]
	p = p + 1
	op.PkNumber = data[p]
	p = p + 1
	op.RegisteredDelivery = data[p]
	p = p + 1
	op.MsgLevel = data[p]
	p = p + 1

	op.ServiceId = &OctetString{Data: data[p : p+10], FixedLen: 10}
	p = p + 10

	op.FeeUserType = data[p]
	p = p + 1

	op.FeeTerminalId = &OctetString{Data: data[p : p+21], FixedLen: 21}
	p = p + 21

	op.TP_pid = data[p]
	p = p + 1
	op.TP_udhi = data[p]
	p = p + 1
	op.MsgFmt = data[p]
	p = p + 1

	op.MsgSrc = &OctetString{Data: data[p : p+6], FixedLen: 6}
	p = p + 6
	op.FeeType = &OctetString{Data: data[p : p+2], FixedLen: 2}
	p = p + 2
	op.FeeCode = &OctetString{Data: data[p : p+6], FixedLen: 6}
	p = p + 6
	op.ValidTime = &OctetString{Data: data[p : p+17], FixedLen: 17}
	p = p + 17
	op.AtTime = &OctetString{Data: data[p : p+17], FixedLen: 17}
	p = p + 17
	op.SrcId = &OctetString{Data: data[p : p+21], FixedLen: 21}
	p = p + 21

	op.DestUsrTl = data[p]
	p = p + 1

	for i := 0; i < int(op.DestUsrTl); i++ {
		op.DestTerminalId = append(
			op.DestTerminalId, &OctetString{Data: data[p : p+21], FixedLen: 21})
		p = p + 21
	}

	op.MsgLength = data[p]
	p = p + 1

	op.MsgContent = data[p : p+int(op.MsgLength)]
	p = p + int(op.MsgLength)

	op.Reserve = &OctetString{Data: data[p : p+8], FixedLen: 8}
	p = p + 8

	return op, nil
}

func (op *Submit) Serialize() []byte {
	b := op.Header.Serialize()

	b = append(b, packUi64(op.MsgId)...)

	b = append(b, packUi8(op.PkTotal)...)
	b = append(b, packUi8(op.PkNumber)...)
	b = append(b, packUi8(op.RegisteredDelivery)...)
	b = append(b, packUi8(op.MsgLevel)...)

	b = append(b, op.ServiceId.Byte()...)
	b = append(b, packUi8(op.FeeUserType)...)
	b = append(b, op.FeeTerminalId.Byte()...)

	b = append(b, packUi8(op.TP_pid)...)
	b = append(b, packUi8(op.TP_udhi)...)
	b = append(b, packUi8(op.MsgFmt)...)

	b = append(b, op.MsgSrc.Byte()...)
	b = append(b, op.FeeType.Byte()...)
	b = append(b, op.FeeCode.Byte()...)
	b = append(b, op.ValidTime.Byte()...)
	b = append(b, op.AtTime.Byte()...)
	b = append(b, op.SrcId.Byte()...)

	b = append(b, packUi8(op.DestUsrTl)...)

	for i := 0; i < int(op.DestUsrTl); i++ {
		b = append(b, op.DestTerminalId[i].Byte()...)
	}

	b = append(b, packUi8(op.MsgLength)...)

	b = append(b, op.MsgContent...)
	b = append(b, op.Reserve.Byte()...)

	return b
}

func (op *Submit) String() string {
	var b bytes.Buffer
	b.WriteString(op.Header.String())

	fmt.Fprintln(&b, "--- Submit ---")
	fmt.Fprintln(&b, "MsgId: ", op.MsgId)
	fmt.Fprintln(&b, "PkTotal: ", op.PkTotal)
	fmt.Fprintln(&b, "PkNumber: ", op.PkNumber)
	fmt.Fprintln(&b, "RegisteredDelivery: ", op.RegisteredDelivery)
	fmt.Fprintln(&b, "MsgLevel: ", op.MsgLevel)

	fmt.Fprintln(&b, "ServiceID: ", op.ServiceId)
	fmt.Fprintln(&b, "FeeUserType: ", op.FeeUserType)
	fmt.Fprintln(&b, "FeeTerminalId: ", op.FeeTerminalId)

	fmt.Fprintln(&b, "TP_pid: ", op.TP_pid)
	fmt.Fprintln(&b, "TP_udhi: ", op.TP_udhi)
	fmt.Fprintln(&b, "MsgFmt: ", op.MsgFmt)

	fmt.Fprintln(&b, "MsgSrc: ", op.MsgSrc)
	fmt.Fprintln(&b, "FeeType: ", op.FeeType)
	fmt.Fprintln(&b, "FeeCode: ", op.FeeCode)
	fmt.Fprintln(&b, "ValidTime: ", op.ValidTime)
	fmt.Fprintln(&b, "AtTime: ", op.AtTime)
	fmt.Fprintln(&b, "SrcId: ", op.SrcId)

	fmt.Fprintln(&b, "DestUsrTl: ", op.DestUsrTl)
	for i := 0; i < int(op.DestUsrTl); i++ {
		fmt.Fprintln(&b, "DestTerminalId: ", op.DestTerminalId[i])
	}

	fmt.Fprintln(&b, "MsgLength: ", op.MsgLength)
	fmt.Fprintln(&b, "MsgContent: ", string(op.MsgContent))

	return b.String()
}

func (op *Submit) Ok() error {
	return nil
}

type SubmitResp struct {
	*Header

	// Body
	MsgId  uint64 // 信息标识
	Result uint8  // 结果
}

func NewSubmitResp(
	sequenceID uint32, msgId uint64, result uint8) (*SubmitResp, error) {

	op := &SubmitResp{}

	op.Header = &Header{}
	var length uint32 = 4 + 4 + 4

	op.MsgId = msgId
	length = length + 8

	op.Result = result
	length = length + 1

	op.Total_Length = length
	op.Command_Id = CMPP_SUBMIT_RESP
	op.Sequence_Id = sequenceID

	return op, nil
}

func ParseSubmitResp(hdr *Header, data []byte) (*SubmitResp, error) {
	op := &SubmitResp{}
	op.Header = hdr

	p := 0
	op.MsgId = unpackUi64(data[p : p+8])
	p = p + 8

	op.Result = data[p]
	p = p + 1

	return op, nil
}

func (op *SubmitResp) Serialize() []byte {
	b := op.Header.Serialize()

	b = append(b, packUi64(op.MsgId)...)
	b = append(b, packUi8(op.Result)...)

	return b
}

func (op *SubmitResp) String() string {
	var b bytes.Buffer
	b.WriteString(op.Header.String())

	fmt.Fprintln(&b, "--- SubmitResp ---")
	fmt.Fprintln(&b, "MsgID: ", op.MsgId)
	fmt.Fprintln(&b, "Result: ", op.Result)

	return b.String()
}

func (op *SubmitResp) Ok() (err error) {

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
