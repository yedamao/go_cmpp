package protocol

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"time"
)

func genTimestamp() uint32 {
	t := time.Now()
	return uint32(int(t.Month())*100000000 + t.Day()*1000000 +
		t.Hour()*10000 + t.Minute()*100 + t.Second())
}

// AuthenticatorSource =
// MD5（Source_Addr+9 字节的0 +shared secret+timestamp）
// Shared secret 由中国移动与源地址实体事先商定，
// timestamp格式为：MMDDHHMMSS，即月日时分秒，10位。
func genAuthenticatorSource(sourceAddr, secret string, timestamp uint32) ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.WriteString(sourceAddr)
	buf.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0})
	buf.WriteString(secret)
	buf.WriteString(fmt.Sprintf("%010d", timestamp))

	h := md5.New()
	_, err := h.Write(buf.Bytes())
	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
