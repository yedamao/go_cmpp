package cmpptest

import (
	"flag"
	"log"
	"net"
	"time"

	connp "github.com/yedamao/go_cmpp/cmpp/conn"
	"github.com/yedamao/go_cmpp/cmpp/protocol"
)

var (
	mo         = flag.Bool("mo", false, "是否模拟上行短信")
	rpt        = flag.Bool("rpt", false, "是否模拟上行状态报告")
	activeTest = flag.Bool("activeTest", false, "是否activeTest")
)

func newSession(rawConn net.Conn) {
	s := &Session{connp.Conn{Conn: rawConn}}

	go s.start()

	if *mo || *rpt {
		go s.deliverWorker()
	}
	if *activeTest {
		go s.activeTestWorker()
	}
}

// 代表sp->运营商的一条连接
type Session struct {
	connp.Conn
	// TODO newSeqFunc
}

func (s *Session) ConnectResp(
	sequenceID uint32, status uint8, auth string) error {

	op, err := protocol.NewConnectResp(sequenceID, status, auth)
	if err != nil {
		return err
	}

	return s.Write(op)
}

func (s *Session) TerminateResp(sequenceID uint32) error {

	op, err := protocol.NewTerminateResp(sequenceID)
	if err != nil {
		return err
	}

	return s.Write(op)
}

func (s *Session) ActiveTest() error {
	var mockId uint32 = 1

	op, err := protocol.NewActiveTest(mockId)
	if err != nil {
		return err
	}

	return s.Write(op)
}

func (s *Session) ActiveTestResp(sequenceID uint32) error {

	op, err := protocol.NewActiveTest(sequenceID)
	if err != nil {
		return err
	}

	return s.Write(op)
}

func (s *Session) SubmitResp(sequenceID uint32, msgId uint64, result uint8) error {

	op, err := protocol.NewSubmitResp(sequenceID, msgId, protocol.OK)
	if err != nil {
		return err
	}

	return s.Write(op)
}

func (s *Session) QueryResp(
	sequenceID uint32,
	time string, queryTye uint8, queryCode string,
	MT_TLMsg, MT_Tlusr uint32,
	MT_Scs, MT_WT, MT_FL, MO_Scs, MO_WT, MO_FL uint32,
) error {

	op, err := protocol.NewQueryResp(
		sequenceID,
		time, queryTye, queryCode,
		MT_TLMsg, MT_Tlusr,
		MT_Scs, MT_WT, MT_FL, MO_Scs, MO_WT, MO_FL,
	)
	if err != nil {
		return err
	}

	return s.Write(op)
}

func (s *Session) Deliver(isReport uint8, content []byte) error {
	var mockId uint32 = 1

	op, err := protocol.NewDeliver(
		mockId, 12345, "1069000000", "", 0, 0, 0, "16611111111",
		isReport, content,
	)

	if err != nil {
		return err
	}

	return s.Write(op)
}

// 模拟状态包
func (s *Session) mockReport() error {
	rpt, err := protocol.NewReport(1234, "DELIVRD", "", "", "17600000000", 0)
	if err != nil {
		return err
	}

	return s.Deliver(protocol.IS_REPORT, rpt.Serialize())
}

// 模拟上行短信
func (s *Session) mockMo() error {

	return s.Deliver(protocol.NOT_REPORT, []byte("hello test msg"))
}

func (s *Session) start() {
	defer s.Close()

	for {
		op, err := s.Read()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Println(err)
			return
		}

		log.Println(op)

		switch op.GetHeader().Command_Id {
		case protocol.CMPP_CONNECT:
			s.ConnectResp(op.GetHeader().Sequence_Id, 0, "mockauth")

		case protocol.CMPP_SUBMIT:
			submit, ok := op.(*protocol.Submit)
			if !ok {
				log.Println("Type assert error: ", op)
			}
			s.SubmitResp(submit.Header.Sequence_Id, submit.MsgId, protocol.OK)

		case protocol.CMPP_QUERY:
			query, ok := op.(*protocol.Query)
			if !ok {
				log.Println("Type assert error: ", op)
			}
			s.QueryResp(
				query.Header.Sequence_Id, query.Time.String(),
				query.Query_Type, query.Query_Code.String(),
				0, 0, 0, 0, 0, 0, 0, 0,
			)

		case protocol.CMPP_DELIVER_RESP:
			log.Println("Deliver Response")

		case protocol.CMPP_ACTIVE_TEST:
			s.ActiveTestResp(op.GetHeader().Sequence_Id)
			log.Println("ActiveTest ... ")
		case protocol.CMPP_ACTIVE_TEST_RESP:
			log.Println("ActiveTest Response")

		case protocol.CMPP_TERMINATE:
			s.TerminateResp(op.GetHeader().Sequence_Id)
			log.Println("Terminate. Close Session")
			return
		case protocol.CMPP_TERMINATE_RESP:
			log.Println("Terminate response. Close Session")
			return

		default:
			log.Println("not support CmdId. close session.")
			return
		}
	}
}

func (s *Session) deliverWorker() {

	doFunc := s.mockMo
	if *rpt {
		log.Println("deliver (report) worker running")
		doFunc = s.mockReport
	} else {
		log.Println("deliver (mo) worker running")
	}

	for {
		tick := time.NewTicker(5 * time.Second)
		select {
		case <-tick.C:
			if err := doFunc(); err != nil {
				log.Println("Deliver error: ", err)
				return
			}
		}
	}
}

func (s *Session) activeTestWorker() {

	log.Println("activeTest worker running")

	for {
		tick := time.NewTicker(30 * time.Second)
		select {
		case <-tick.C:
			if err := s.ActiveTest(); err != nil {
				log.Println("ActiveTest error: ", err)
				return
			}
			log.Println("send ActiveTest ...")
		}
	}
}
