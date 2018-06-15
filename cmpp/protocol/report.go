package protocol

import (
	"errors"
)

// 状态报告
type Report struct {
	// SP提交短信（CMPP_SUBMIT）操作时，
	// 与SP相连的ISMG产生的Msg_Id
	MsgId          uint64
	Stat           *OctetString // 发送短信的应答结果
	SubmitTime     *OctetString
	DoneTime       *OctetString
	DestTerminalId *OctetString
	SMSCSequence   uint32
}

func ParseReport(rawData []byte) (*Report, error) {
	if len(rawData) != 60 {
		return nil, errors.New("ParseReport: len error")
	}

	p := 0
	rpt := &Report{}

	rpt.MsgId = unpackUi64(rawData[p : p+8])
	p = p + 8
	rpt.Stat = &OctetString{Data: rawData[p : p+7]}
	p = p + 7
	rpt.SubmitTime = &OctetString{Data: rawData[p : p+10]}
	p = p + 10
	rpt.DoneTime = &OctetString{Data: rawData[p : p+10]}
	p = p + 10
	rpt.DestTerminalId = &OctetString{Data: rawData[p : p+21]}
	p = p + 21
	rpt.SMSCSequence = unpackUi32(rawData[p : p+4])
	p = p + 4

	return rpt, nil
}

func NewReport(
	msgId uint64,
	stat, submitTime, doneTime, destTerminalId string,
	smscSequence uint32,
) (*Report, error) {

	rpt := &Report{}
	rpt.MsgId = msgId
	rpt.Stat = &OctetString{Data: []byte(stat), FixedLen: 7}
	rpt.SubmitTime = &OctetString{Data: []byte(submitTime), FixedLen: 10}
	rpt.DoneTime = &OctetString{Data: []byte(doneTime), FixedLen: 10}
	rpt.DestTerminalId = &OctetString{Data: []byte(destTerminalId), FixedLen: 21}
	rpt.SMSCSequence = smscSequence

	return rpt, nil
}

func (r *Report) Serialize() []byte {

	b := packUi64(r.MsgId)
	b = append(b, r.Stat.Byte()...)
	b = append(b, r.SubmitTime.Byte()...)
	b = append(b, r.DoneTime.Byte()...)
	b = append(b, r.DestTerminalId.Byte()...)
	b = append(b, packUi32(r.SMSCSequence)...)

	return b
}
