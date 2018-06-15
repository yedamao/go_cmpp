package cmpp

import (
	"errors"
	"net"

	"github.com/yedamao/go_cmpp/cmpp/conn"
	"github.com/yedamao/go_cmpp/cmpp/protocol"
)

type SequenceFunc func() uint32

type Cmpp struct {
	conn.Conn
	newSeqNum SequenceFunc

	spId string // SP的企业代码
}

func NewCmpp(
	addr string, sourceAddr, sharedSecret string,
	newSeqNum SequenceFunc,
) (*Cmpp, error) {
	if nil == newSeqNum {
		return nil, errors.New("newSeqNum must not be nil")
	}

	s := &Cmpp{
		newSeqNum: newSeqNum,
		spId:      sourceAddr,
	}

	if err := s.connect(addr); err != nil {
		return nil, err
	}

	// 登陆
	if err := s.Connect(sourceAddr, sharedSecret); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Cmpp) connect(addr string) error {
	connection, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	s.Conn = conn.Conn{Conn: connection}

	return nil
}

func (s *Cmpp) Connect(sourceAddr, sharedSecret string) error {
	op, err := protocol.NewConnect(s.newSeqNum(), sourceAddr, sharedSecret)
	if err != nil {
		return err
	}
	if err = s.Write(op); err != nil {
		return err
	}

	// Read block
	var resp protocol.Operation
	if resp, err = s.Read(); err != nil {
		return err
	}

	if resp.GetHeader().Command_Id != protocol.CMPP_CONNECT_RESP {
		return errors.New("Connect Resp Wrong RequestID")
	}

	return resp.Ok()
}

func (s *Cmpp) Terminate() error {

	op, err := protocol.NewTerminate(s.newSeqNum())
	if err != nil {
		return err
	}

	return s.Write(op)
}

func (s *Cmpp) TerminateResp(sequenceID uint32) error {

	op, err := protocol.NewTerminateResp(sequenceID)
	if err != nil {
		return err
	}

	return s.Write(op)
}

func (s *Cmpp) ActiveTest() error {

	op, err := protocol.NewActiveTest(s.newSeqNum())
	if err != nil {
		return err
	}

	return s.Write(op)
}

func (s *Cmpp) ActiveTestResp(sequenceID uint32) error {

	op, err := protocol.NewActiveTestResp(sequenceID)
	if err != nil {
		return err
	}

	return s.Write(op)
}

func (s *Cmpp) Cancel(msgId uint64) error {

	op, err := protocol.NewCancel(s.newSeqNum(), msgId)
	if err != nil {
		return err
	}

	return s.Write(op)
}

func (s *Cmpp) DeliverResp(sequenceID uint32, msgId uint64, result uint8) error {

	op, err := protocol.NewDeliverResp(sequenceID, msgId, protocol.OK)
	if err != nil {
		return err
	}

	return s.Write(op)
}

func (s *Cmpp) Submit(
	pkTotal, pkNumber, needReport, msgLevel uint8,
	serviceId string, feeUserType uint8, feeTerminalId string,
	msgFmt uint8, feeType, feeCode, srcId string,
	destTermId []string, content []byte,
) (uint32, error) {

	var (
		TP_udhi    uint8  = 0
		sequenceID uint32 = s.newSeqNum()
	)

	if pkTotal > 1 {
		TP_udhi = 1
	}

	op, err := protocol.NewSubmit(
		sequenceID,
		pkTotal, pkNumber, needReport, msgLevel,
		serviceId, feeUserType, feeTerminalId,
		0, TP_udhi, msgFmt,
		s.spId, feeType, feeCode, "", "", srcId,
		destTermId, content,
	)

	if err != nil {
		return sequenceID, err
	}

	return sequenceID, s.Write(op)
}

func (s *Cmpp) Query(time, serviceId string) error {

	var (
		queryTye  uint8
		queryCode string
	)
	if "" != serviceId {
		queryTye = 1
		queryCode = serviceId
	}

	op, err := protocol.NewQuery(
		s.newSeqNum(), time, queryTye, queryCode)
	if err != nil {
		return err
	}

	return s.Write(op)
}
