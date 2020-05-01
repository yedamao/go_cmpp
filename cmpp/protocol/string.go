package protocol

import (
	"bytes"
	"fmt"

	"github.com/yedamao/encoding"
)

var (
	LongtextPrefixBS = []byte{0x05, 0x00, 0x03}
)

type MsgContentSegment struct {
	Ref     byte
	Total   byte
	No      byte
	Content string
}

func ParseMsgContentSegment(msgFmt uint8, data []byte) (*MsgContentSegment, error) {
	if len(data) > 6 && bytes.Equal(LongtextPrefixBS, data[:3]) {
		content, err := parseMsgContent(msgFmt, data[6:])
		if err != nil {
			return nil, err
		}
		return &MsgContentSegment{
			Ref:     data[3],
			Total:   data[4],
			No:      data[5],
			Content: content,
		}, nil
	}

	content, err := parseMsgContent(msgFmt, data[6:])
	if err != nil {
		return nil, err
	}
	return &MsgContentSegment{
		Ref:     0,
		Total:   1,
		No:      1,
		Content: content,
	}, nil
}

func parseMsgContent(msgFmt uint8, data []byte) (string, error) {
	switch msgFmt {
	case ASCII:
		return string(data), nil
	case BINARY:
		return fmt.Sprintf("%x", data), nil
	case UCS2:
		return string(encoding.UCS22UTF8(data)), nil
	case GB18030:
		return string(encoding.GBK2UTF8(data)), nil
	default:
		return "", fmt.Errorf("不支持消息格式: %v", msgFmt)
	}
}
