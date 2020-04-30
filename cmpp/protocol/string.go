package protocol

import (
	"fmt"

	"github.com/yedamao/encoding"
)

func ParseMsgContent(msgFmt int, data []byte) (string, error) {
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
