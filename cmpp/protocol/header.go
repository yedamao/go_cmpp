package protocol

import (
	"bytes"
	"fmt"
)

// Message Header
type Header struct {
	// 消息总长度(含消息头及消息体)
	Total_Length uint32
	// 命令或响应类型
	Command_Id uint32
	// 消息流水号,顺序累加,步长为1,循环使用（一对请求和应答消息的流水号必须相同）
	Sequence_Id uint32
}

func (p *Header) GetHeader() *Header {
	return p
}

func (p *Header) Serialize() []byte {
	b := packUi32(p.Total_Length)
	b = append(b, packUi32(p.Command_Id)...)
	b = append(b, packUi32(p.Sequence_Id)...)

	return b
}

func (p *Header) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "--- Header ---")
	fmt.Fprintln(&b, "Total_Length: ", p.Total_Length)
	fmt.Fprintf(&b, "Command_Id: 0x%x\n", p.Command_Id)
	fmt.Fprintln(&b, "Sequence_Id: ", p.Sequence_Id)

	return b.String()
}

func (p *Header) Parse(data []byte) *Header {

	p.Total_Length = unpackUi32(data[:4])
	p.Command_Id = unpackUi32(data[4:8])
	p.Sequence_Id = unpackUi32(data[8:12])

	return p
}

func ParseHeader(data []byte) (*Header, error) {

	h := &Header{}
	h.Parse(data)

	return h, nil
}
