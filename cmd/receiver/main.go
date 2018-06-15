package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yedamao/go_cmpp/cmpp"
	"github.com/yedamao/go_cmpp/cmpp/protocol"
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
		fmt.Println("Arg error: sourceAddr or sharedSecret must not be empty .")
		flag.Usage()
		os.Exit(-1)
	}

	ts, err := cmpp.NewCmpp(*addr, *sourceAddr, *sharedSecret, newSeqNum)
	if err != nil {
		fmt.Println("Connection Err", err)
		os.Exit(-1)
	}
	fmt.Println("connect succ")

	for {
		op, err := ts.Read() // This is blocking
		if err != nil {
			fmt.Println("Read Err:", err)
			break
		}

		switch op.GetHeader().Command_Id {
		case protocol.CMPP_DELIVER:
			dlv, ok := op.(*protocol.Deliver)
			if !ok {
				log.Println("Type assert error: ", op)
			}

			if dlv.RegisteredDelivery == protocol.IS_REPORT {
				// 状态报告
				rpt, err := protocol.ParseReport(dlv.MsgContent)
				if err != nil {
					log.Println(err)
				}
				fmt.Println(rpt)
			} else {
				// 上行短信
				fmt.Println(dlv)
			}
			ts.DeliverResp(dlv.Header.Sequence_Id, dlv.MsgId, protocol.OK)

		case protocol.CMPP_ACTIVE_TEST:
			fmt.Println("recv ActiveTest")
			ts.ActiveTestResp(op.GetHeader().Sequence_Id)

		case protocol.CMPP_TERMINATE:
			fmt.Println("recv Terminate")
			ts.TerminateResp(op.GetHeader().Sequence_Id)
			ts.Close()
			return

		case protocol.CMPP_TERMINATE_RESP:
			fmt.Println("Terminate response")
			ts.Close()
			return

		default:
			fmt.Printf("Unexpect CmdId: %0x\n", op.GetHeader().Command_Id)
			ts.Close()
			return
		}
	}
}
