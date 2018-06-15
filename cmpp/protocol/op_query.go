package protocol

import (
	"bytes"
	"fmt"
)

// SP向ISMG查询发送短信状态（CMPP_QUERY）操作
//
// CMPP_QUERY操作的目的是SP向ISMG查询某时间的
// 业务统计情况，可以按总数或按业务代码查询。
// ISMG以CMPP_QUERY_RESP应答。
type Query struct {
	*Header

	// Body
	// 时间YYYYMMDD(精确至日)
	Time *OctetString
	// 查询类别
	// 0：总数查询
	// 1：按业务类型查询
	Query_Type uint8
	// 查询码
	// 当Query_Type为0时，此项无效；
	// 当Query_Type为1时，此项填写业务类型Service_Id.
	Query_Code *OctetString
	Reserve    *OctetString
}

func NewQuery(
	sequenceID uint32,
	time string, queryTye uint8, queryCode string,
) (*Query, error) {

	op := &Query{}

	op.Header = &Header{}
	var length uint32 = 4 + 4 + 4 // header length

	op.Time = &OctetString{Data: []byte(time), FixedLen: 8}
	length = length + 8

	op.Query_Type = queryTye
	length = length + 1

	op.Query_Code = &OctetString{Data: []byte(queryCode), FixedLen: 10}
	length = length + 10

	op.Reserve = &OctetString{FixedLen: 8}
	length = length + 8

	op.Total_Length = length
	op.Command_Id = CMPP_QUERY
	op.Sequence_Id = sequenceID

	return op, nil
}

func ParseQuery(hdr *Header, data []byte) (*Query, error) {
	p := 0

	op := &Query{}
	op.Header = hdr

	op.Time = &OctetString{Data: data[p : p+8]}
	p = p + 8

	op.Query_Type = data[p]
	p = p + 1

	op.Query_Code = &OctetString{Data: data[p : p+10]}
	p = p + 10

	op.Reserve = &OctetString{Data: data[p : p+8]}
	p = p + 8

	return op, nil
}

func (op *Query) Serialize() []byte {
	b := op.Header.Serialize()

	b = append(b, op.Time.Byte()...)
	b = append(b, packUi8(op.Query_Type)...)
	b = append(b, op.Query_Code.Byte()...)
	b = append(b, op.Reserve.Byte()...)

	return b
}

func (op *Query) String() string {
	var b bytes.Buffer
	b.WriteString(op.Header.String())

	fmt.Fprintln(&b, "--- Query ---")
	fmt.Fprintln(&b, "Time: ", op.Time)
	fmt.Fprintln(&b, "Query_Type: ", op.Query_Type)
	fmt.Fprintln(&b, "Query_Code: ", op.Query_Code)
	fmt.Fprintln(&b, "Reserve: ", op.Reserve)

	return b.String()
}

func (op *Query) Ok() error {
	return nil
}

type QueryResp struct {
	*Header

	// Body
	// 时间YYYYMMDD(精确至日)
	Time *OctetString
	// 查询类别
	// 0：总数查询
	// 1：按业务类型查询
	Query_Type uint8
	// 查询码
	// 当Query_Type为0时，此项无效；
	// 当Query_Type为1时，此项填写业务类型Service_Id.
	Query_Code *OctetString

	MT_TLMsg uint32 // 从SP接收信息总数
	MT_Tlusr uint32 // 从SP接收用户总数

	MT_Scs uint32 // 成功转发数量
	MT_WT  uint32 // 待转发数量
	MT_FL  uint32 // 转发失败数量
	MO_Scs uint32 // 向SP成功送达数量
	MO_WT  uint32 // 向SP待送达数量
	MO_FL  uint32 // 向SP送达失败数量
}

func NewQueryResp(
	sequenceID uint32,
	time string, queryTye uint8, queryCode string,
	MT_TLMsg, MT_Tlusr uint32,
	MT_Scs, MT_WT, MT_FL, MO_Scs, MO_WT, MO_FL uint32,
) (*QueryResp, error) {

	op := &QueryResp{}

	op.Header = &Header{}
	var length uint32 = 4 + 4 + 4 // header length

	op.Time = &OctetString{Data: []byte(time), FixedLen: 8}
	length = length + 8

	op.Query_Type = queryTye
	length = length + 1

	op.Query_Code = &OctetString{Data: []byte(queryCode), FixedLen: 10}
	length = length + 10

	op.MT_TLMsg = MT_TLMsg
	length = length + 4
	op.MT_Tlusr = MT_Tlusr
	length = length + 4

	op.MT_Scs = MT_Scs
	length = length + 4
	op.MT_WT = MT_WT
	length = length + 4
	op.MT_FL = MT_FL
	length = length + 4

	op.MO_Scs = MO_Scs
	length = length + 4
	op.MO_WT = MO_WT
	length = length + 4
	op.MO_FL = MO_FL
	length = length + 4

	op.Total_Length = length
	op.Command_Id = CMPP_QUERY_RESP
	op.Sequence_Id = sequenceID

	return op, nil
}

func ParseQueryResp(hdr *Header, data []byte) (*QueryResp, error) {
	p := 0

	op := &QueryResp{}
	op.Header = hdr

	op.Time = &OctetString{Data: data[p : p+8]}
	p = p + 8

	op.Query_Type = data[p]
	p = p + 1

	op.Query_Code = &OctetString{Data: data[p : p+10]}
	p = p + 10

	op.MT_TLMsg = unpackUi32(data[p : p+4])
	p = p + 4
	op.MT_Tlusr = unpackUi32(data[p : p+4])
	p = p + 4

	op.MT_Scs = unpackUi32(data[p : p+4])
	p = p + 4
	op.MT_WT = unpackUi32(data[p : p+4])
	p = p + 4
	op.MT_FL = unpackUi32(data[p : p+4])
	p = p + 4
	op.MO_Scs = unpackUi32(data[p : p+4])
	p = p + 4
	op.MO_WT = unpackUi32(data[p : p+4])
	p = p + 4
	op.MO_FL = unpackUi32(data[p : p+4])
	p = p + 4

	return op, nil
}

func (op *QueryResp) Serialize() []byte {
	b := op.Header.Serialize()

	b = append(b, op.Time.Byte()...)
	b = append(b, packUi8(op.Query_Type)...)
	b = append(b, op.Query_Code.Byte()...)

	b = append(b, packUi32(op.MT_TLMsg)...)
	b = append(b, packUi32(op.MT_Tlusr)...)

	b = append(b, packUi32(op.MT_Scs)...)
	b = append(b, packUi32(op.MT_WT)...)
	b = append(b, packUi32(op.MT_FL)...)
	b = append(b, packUi32(op.MO_Scs)...)
	b = append(b, packUi32(op.MO_WT)...)
	b = append(b, packUi32(op.MO_FL)...)

	return b
}

func (op *QueryResp) String() string {
	var b bytes.Buffer
	b.WriteString(op.Header.String())

	fmt.Fprintln(&b, "--- Query ---")
	fmt.Fprintln(&b, "Time: ", op.Time)
	fmt.Fprintln(&b, "Query_Type: ", op.Query_Type)
	fmt.Fprintln(&b, "Query_Code: ", op.Query_Code)

	fmt.Fprintln(&b, "从SP接收信息总数: ", op.MT_TLMsg)
	fmt.Fprintln(&b, "从SP接收用户总数: ", op.MT_Tlusr)

	fmt.Fprintln(&b, "成功转发数量: ", op.MT_Scs)
	fmt.Fprintln(&b, "待转发数量: ", op.MT_WT)
	fmt.Fprintln(&b, "转发失败数量: ", op.MT_FL)
	fmt.Fprintln(&b, "向SP成功送达数量: ", op.MO_Scs)
	fmt.Fprintln(&b, "向SP待送达数量: ", op.MO_WT)
	fmt.Fprintln(&b, "向SP送达失败数量: ", op.MO_FL)

	return b.String()
}

func (op *QueryResp) Ok() error {
	return nil
}
