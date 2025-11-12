package main

import (
	"flag"
	"log"
	"os"

	"github.com/yedamao/go_cmpp/cmpp"
	"github.com/yedamao/go_cmpp/cmpp/protocol"
	"github.com/yedamao/go_cmpp/utils"
)

var (
	addr         = flag.String("addr", ":7890", "smgw addr(运营商地址)")
	sourceAddr   = flag.String("sourceAddr", "", "源地址，即SP的企业代码")
	sharedSecret = flag.String("secret", "", "登陆密码")
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

	ts, err := cmpp.NewCmpp(*addr, *sourceAddr, *sharedSecret, newSeqNum)
	if err != nil {
		log.Println("Connection Err", err)
		os.Exit(-1)
	}
	log.Println("connect succ")

	for {
		op, err := ts.Read() // This is blocking
		if err != nil {
			log.Println("Read Err:", err)
			break
		}

		switch op.GetHeader().Command_Id {
		case protocol.CMPP_DELIVER:
			dlv, ok := op.(*protocol.Deliver)
			if !ok {
				log.Println("Type assert error: ", op)
			}
			if dlv.RegisteredDelivery == protocol.IS_REPORT {
				rpt, err := protocol.ParseReport(dlv.MsgContent)
				if err != nil {
					log.Println(err)
				}
				log.Printf("recv DeliverReport：%s\nReport：%s\n\n", utils.ToJsonString(dlv), utils.ToJsonString(rpt))
			} else {
				log.Printf("recv Deliver：%s\n", utils.ToJsonString(dlv))
			}
			ts.DeliverResp(dlv.Header.Sequence_Id, dlv.MsgId, protocol.OK)

		case protocol.CMPP_ACTIVE_TEST:
			log.Println("recv ActiveTest")
			ts.ActiveTestResp(op.GetHeader().Sequence_Id)

		case protocol.CMPP_TERMINATE:
			log.Println("recv Terminate")
			ts.TerminateResp(op.GetHeader().Sequence_Id)
			ts.Close()
			return

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
