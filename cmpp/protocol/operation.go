package protocol

import (
	"errors"
	"fmt"
)

type Operation interface {
	//Header returns the Operation header, decoded. Header fields
	// can be updated before reserialzing .
	GetHeader() *Header

	// SerializeTo encodes Operation to it's binary form,
	// include the header and body
	Serialize() []byte

	// String
	String() string

	// if status/result not ok, return error
	Ok() error
}

func ParseOperation(data []byte) (Operation, error) {
	if len(data) < 12 {
		return nil, errors.New("Invalide data length")
	}

	header, err := ParseHeader(data)
	if err != nil {
		return nil, err
	}

	if int(header.Total_Length) != len(data) {
		return nil, errors.New("Invalide data length")
	}

	var op Operation

	switch header.Command_Id {
	case CMPP_CONNECT:
		op, err = ParseConnect(header, data[12:])
	case CMPP_CONNECT_RESP:
		op, err = ParseConnectResp(header, data[12:])

	case CMPP_SUBMIT:
		op, err = ParseSubmit(header, data[12:])
	case CMPP_SUBMIT_RESP:
		op, err = ParseSubmitResp(header, data[12:])

	case CMPP_DELIVER:
		op, err = ParseDeliver(header, data[12:])
	case CMPP_DELIVER_RESP:
		op, err = ParseDeliverResp(header, data[12:])

	case CMPP_ACTIVE_TEST:
		op, err = ParseActiveTest(header, data[12:])
	case CMPP_ACTIVE_TEST_RESP:
		op, err = ParseActiveTestResp(header, data[12:])

	case CMPP_CANCEL:
		op, err = ParseCancel(header, data[12:])
	case CMPP_CANCEL_RESP:
		op, err = ParseCancelResp(header, data[12:])

	case CMPP_TERMINATE:
		op, err = ParseTermanite(header, data[12:])
	case CMPP_TERMINATE_RESP:
		op, err = ParseTermaniteResp(header, data[12:])

	case CMPP_QUERY:
		op, err = ParseQuery(header, data[12:])
	case CMPP_QUERY_RESP:
		op, err = ParseQueryResp(header, data[12:])

	default:
		err = fmt.Errorf("Unknow Operation CmdId: 0x%x", header.Command_Id)
	}

	return op, err
}
