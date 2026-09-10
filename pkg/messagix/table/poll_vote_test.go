package table_test

import (
	"testing"

	"go.mewis.me/meta-extra/pkg/messagix/lightspeed"
	"go.mewis.me/meta-extra/pkg/messagix/table"
)

func TestRemovePollVoteV2Dependency(t *testing.T) {
	dependencies := table.SPToDepMap([]string{"removePollVoteV2"})
	result := &table.LSTable{}
	decoder := lightspeed.NewLightSpeedDecoder(dependencies, result)
	i64 := func(value string) []interface{} { return []interface{}{float64(lightspeed.I64_FROM_STRING), value} }
	decoder.Decode([]interface{}{float64(lightspeed.CALL_STORED_PROCEDURE), "removePollVoteV2", i64("11"), i64("7"), i64("9"), i64("30"), i64("0"), i64("3"), "m1", i64("4")})
	if len(result.LSRemovePollVoteV2) != 1 {
		t.Fatalf("unexpected remove vote decode: %#v", result.LSRemovePollVoteV2)
	}
	vote := result.LSRemovePollVoteV2[0]
	if vote.OptionID != 11 || vote.PollID != 7 || vote.ContactID != 9 || vote.ThreadKey != 3 || vote.MessageID != "m1" {
		t.Fatalf("unexpected remove vote: %#v", vote)
	}
}
