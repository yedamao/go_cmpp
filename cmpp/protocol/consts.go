package protocol

const (
	// 高位4bit表示主版本号,低位4bit表示次版本号
	VERSION uint8 = 20
)

// Command_Id定义
const (
	CMPP_CONNECT   = 0x00000001
	CMPP_TERMINATE = 0x00000002
)

const (
	CMPP_SUBMIT = 0x00000004 + iota
	CMPP_DELIVER
	CMPP_QUERY
	CMPP_CANCEL
	CMPP_ACTIVE_TEST
)

const (
	CMPP_CONNECT_RESP   = 0x80000001
	CMPP_TERMINATE_RESP = 0x80000002
)

const (
	CMPP_SUBMIT_RESP = 0x80000004 + iota
	CMPP_DELIVER_RESP
	CMPP_QUERY_RESP
	CMPP_CANCEL_RESP
	CMPP_ACTIVE_TEST_RESP
)

// Status/Result OK
const (
	OK = 0
)

// 信息格式
//  0：ASCII串
//  3：短信写卡操作
//  4：二进制信息
//  8：UCS2编码
//	15：含GB汉字
const (
	ASCII   = 0  // ASCII编码
	BINARY  = 4  // 二进制短消息
	UCS2    = 8  // UCS2编码
	GB18030 = 15 // GB18030编码
)

// 是否要求返回状态报告
const (
	NO_NEED_REPORT = 0
	NEED_REPORT    = 1
)

// 是否为状态报告
const (
	NOT_REPORT = 0
	IS_REPORT  = 1
)
