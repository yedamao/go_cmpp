package protocol

import (
	"testing"
)

func TestQuery(t *testing.T) {
	op, err := NewQuery(
		1, "20180612", 0, "",
	)
	if err != nil {
		t.Error(err)
	}

	parsed, err := ParseOperation(op.Serialize())
	if err != nil {
		t.Fatal(err)
	}

	query := parsed.(*Query)

	if query.Time.String() != op.Time.String() ||
		query.Query_Type != op.Query_Type ||
		query.Query_Code.String() != op.Query_Code.String() {
		t.Error("not equal")
	}
}

func TestQueryResp(t *testing.T) {
	op, err := NewQueryResp(
		1, "20180612", 0, "",
		499, 499, 9, 9, 9, 9, 9, 9,
	)
	if err != nil {
		t.Error(err)
	}

	parsed, err := ParseOperation(op.Serialize())
	if err != nil {
		t.Fatal(err)
	}

	resp := parsed.(*QueryResp)

	if resp.MT_TLMsg != op.MT_TLMsg ||
		resp.MT_Tlusr != op.MT_Tlusr ||

		resp.MT_Scs != op.MT_Scs ||
		resp.MT_WT != op.MT_WT ||
		resp.MT_FL != op.MT_FL ||
		resp.MO_Scs != op.MO_Scs ||
		resp.MO_WT != op.MO_WT ||
		resp.MO_FL != op.MO_FL {
		t.Error("not equal")
	}
}
