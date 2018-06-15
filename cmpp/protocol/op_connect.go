package protocol

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

// SP请求连接到ISMG（CMPP_CONNECT）操作
//
// CMPP_CONNECT操作的目的是SP向ISMG注册作为一个合法SP身份，
// 若注册成功后即建立了应用层的连接，此后SP可以通过此ISMG
// 接收和发送短信。
// ISMG以CMPP_CONNECT_RESP消息响应SP的请求。
type Connect struct {
	*Header

	// 源地址，此处为SP_Id，即SP的企业代码
	SourceAddr *OctetString
	// 用于鉴别源地址
	AuthenticatorSource *OctetString
	// 双方协商的版本号(高位4bit表示主版本号,
	// 低位4bit表示次版本号)
	Version uint8
	// 时间戳的明文,由客户端产生,格式为MMDDHHMMSS，
	// 即月日时分秒，10位数字的整型，右对齐 。
	Timestamp uint32
}

func NewConnect(
	sequenceID uint32, sourceAddr, sharedSecret string,
) (*Connect, error) {
	// TODO
	// gen AuthenticatorSource
	auth, err := genAuthenticatorSource(sourceAddr, sharedSecret, genTimestamp())
	if err != nil {
		return nil, err
	}

	op := &Connect{}

	op.Header = &Header{}
	var length uint32 = 4 + 4 + 4 // header length

	op.SourceAddr = &OctetString{Data: []byte(sourceAddr), FixedLen: 6}
	length = length + 6

	op.AuthenticatorSource = &OctetString{Data: []byte(auth), FixedLen: 16}
	length = length + 16

	op.Version = VERSION
	length = length + 1

	op.Timestamp = genTimestamp()
	length = length + 4

	op.Total_Length = length
	op.Command_Id = CMPP_CONNECT
	op.Sequence_Id = sequenceID

	return op, nil
}

func ParseConnect(hdr *Header, data []byte) (*Connect, error) {
	p := 0

	op := &Connect{}
	op.Header = hdr

	op.SourceAddr = &OctetString{Data: data[p : p+6], FixedLen: 6}
	p = p + 6

	op.AuthenticatorSource = &OctetString{Data: data[p : p+16], FixedLen: 6}
	p = p + 16

	op.Version = data[p]
	p = p + 1

	op.Timestamp = unpackUi32(data[p : p+4])
	p = p + 4

	return op, nil
}

func (op *Connect) Serialize() []byte {
	b := op.Header.Serialize()

	b = append(b, op.SourceAddr.Byte()...)
	b = append(b, op.AuthenticatorSource.Byte()...)
	b = append(b, packUi8(op.Version)...)
	b = append(b, packUi32(op.Timestamp)...)

	return b
}

func (op *Connect) String() string {
	var b bytes.Buffer
	b.WriteString(op.Header.String())

	fmt.Fprintln(&b, "--- Login ---")
	fmt.Fprintln(&b, "SourceAddr: ", op.SourceAddr)
	fmt.Fprintln(&b, "AuthenticatorSource: ", op.AuthenticatorSource)
	fmt.Fprintln(&b, "Version: ", op.Version)
	fmt.Fprintln(&b, "Timestamp: ", op.Timestamp)

	return b.String()
}

func (op *Connect) Ok() error {
	// just for interface
	return nil
}

type ConnectResp struct {
	*Header

	// 状态
	Status uint8
	// ISMG认证码，用于鉴别ISMG
	AuthenticatorISMG *OctetString
	// 服务器支持的最高版本号
	Version uint8
}

func NewConnectResp(
	sequenceID uint32, status uint8, authISMG string,
) (*ConnectResp, error) {
	op := &ConnectResp{}

	op.Header = &Header{}
	var length uint32 = 4 + 4 + 4

	op.Status = status
	length = length + 1

	op.AuthenticatorISMG = &OctetString{Data: []byte(authISMG), FixedLen: 16}
	length = length + 16

	op.Version = VERSION
	length = length + 1

	op.Total_Length = length
	op.Command_Id = CMPP_CONNECT_RESP
	op.Sequence_Id = sequenceID

	return op, nil
}

func ParseConnectResp(hdr *Header, data []byte) (*ConnectResp, error) {
	op := &ConnectResp{}
	op.Header = hdr

	p := 0
	op.Status = data[p]
	p = p + 1

	op.AuthenticatorISMG = &OctetString{Data: data[p : p+16], FixedLen: 16}
	p = p + 16

	op.Version = data[p]
	p = p + 1

	return op, nil
}

func (op *ConnectResp) Serialize() []byte {
	b := op.Header.Serialize()

	b = append(b, packUi8(op.Status)...)
	b = append(b, op.AuthenticatorISMG.Byte()...)
	b = append(b, packUi8(op.Version)...)

	return b
}

func (op *ConnectResp) String() string {
	var b bytes.Buffer
	b.WriteString(op.Header.String())

	fmt.Fprintln(&b, "--- ConnectResp ---")
	fmt.Fprintln(&b, "Status: ", op.Status)
	fmt.Fprintln(&b, "AuthenticatorISMG: ", op.AuthenticatorISMG)
	fmt.Fprintln(&b, "Version: ", op.Version)

	return b.String()
}

func (op *ConnectResp) Ok() (err error) {

	switch op.Status {
	case 0: // 正确
		err = nil
	case 1:
		err = errors.New("消息结构错")
	case 2:
		err = errors.New("非法源地址")
	case 3:
		err = errors.New("认证错")
	case 4:
		err = errors.New("版本太高")
	default:
		err = errors.New("其他错误: " + strconv.Itoa(int(op.Status)))
	}

	return err
}
