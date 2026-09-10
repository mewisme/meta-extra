package table_test

import (
	"testing"

	"go.mewis.me/meta-extra/pkg/messagix/lightspeed"
	"go.mewis.me/meta-extra/pkg/messagix/table"
)

func TestApplyAdminMessageCTAV2Dependency(t *testing.T) {
	dependencies := table.SPToDepMap([]string{"applyAdminMessageCTAV2"})
	if dependencies["applyAdminMessageCTAV2"] != "LSApplyAdminMessageCTAV2" {
		t.Fatalf("unexpected dependency mapping: %q", dependencies["applyAdminMessageCTAV2"])
	}
	result := &table.LSTable{}
	decoder := lightspeed.NewLightSpeedDecoder(dependencies, result)
	decoder.Decode([]interface{}{float64(lightspeed.CALL_STORED_PROCEDURE), "applyAdminMessageCTAV2", []interface{}{float64(lightspeed.I64_FROM_STRING), "3"}, []interface{}{float64(lightspeed.I64_FROM_STRING), "100"}, "m1", "View poll", "admin_msg_poll_details", nil, nil, nil, false, nil, nil, []interface{}{float64(lightspeed.I64_FROM_STRING), "7"}, "Question", []interface{}{float64(lightspeed.I64_FROM_STRING), "0"}, nil, nil})
	if len(result.LSApplyAdminMessageCTAV2) != 1 {
		t.Fatalf("unexpected CTA decode count: %d", len(result.LSApplyAdminMessageCTAV2))
	}
	cta := result.LSApplyAdminMessageCTAV2[0]
	if cta.ThreadKey != 3 || cta.MessageID != "m1" || cta.PollID != 7 || cta.PollTitle != "Question" {
		t.Fatalf("unexpected CTA decode: %#v", cta)
	}
}
