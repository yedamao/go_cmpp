package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/yedamao/encoding"
	"github.com/yedamao/go_cmpp/utils"

	"github.com/yedamao/go_cmpp/cmpp"
	"github.com/yedamao/go_cmpp/cmpp/protocol"
)

var (
	addr         = flag.String("addr", ":7890", "smgw addr(运营商地址)")
	sourceAddr   = flag.String("sourceAddr", "", "源地址，即SP的企业代码")
	sharedSecret = flag.String("secret", "", "登陆密码")

	serviceId = flag.String("serviceId", "", "业务类型，是数字、字母和符号的组合")

	srcId      = flag.String("srcId", "", "SP的接入号码")
	destNumber = flag.String("dest-number", "", "接收手机号码, 86..., 多个使用，分割")
	msg        = flag.String("msg", "", "短信内容")
)

func init() {
	flag.Parse()
}

var sequenceID uint32 = 0

func newSeqNum() uint32 {
	sequenceID++

	return sequenceID
}

func main() {
	if "" == *sourceAddr || "" == *sharedSecret {
		log.Println("Arg error: sourceAddr or sharedSecret must not be empty .")
		flag.Usage()
		os.Exit(-1)
	}

	destNumbers := strings.Split(*destNumber, ",")
	log.Println("destNumbers: ", destNumbers)

	ts, err := cmpp.NewCmpp(*addr, *sourceAddr, *sharedSecret, newSeqNum)
	if err != nil {
		log.Println("Connection Err", err)
		os.Exit(-1)
	}
	log.Println("connect succ")
	// encoding msg
	content := encoding.UTF82GBK([]byte(*msg))

	if len(content) > 140 {
		log.Println("msg Err: not suport long sms")
	}

	_, err = ts.Submit(
		1, 1, 1, 0, *serviceId, 0, "", protocol.GB18030,
		"02", "", *srcId, destNumbers, content,
	)
	if err != nil {
		log.Println("Submit err ", err)
		os.Exit(-1)
	}

	for {
		op, err := ts.Read() // This is blocking
		if err != nil {
			log.Println("Read Err:", err)
			break
		}

		switch op.GetHeader().Command_Id {
		case protocol.CMPP_SUBMIT_RESP:
			ts.Terminate()
			if err := op.Ok(); err != nil {
				log.Println(err)
			} else {
				submitResp, ok := op.(*protocol.SubmitResp)
				if ok {
					println("Submit Ok：", utils.ToJsonString(submitResp))
				} else {
					log.Println("Submit fail")
				}
			}

		case protocol.CMPP_TERMINATE_RESP:
			log.Println("Terminate response")
			ts.Close()
			return

		default:
			log.Printf("Unexpect CmdId: %0x\n", op.GetHeader().Command_Id)
			ts.Close()
			return
		}
	}
}
