package table_test

import (
	"testing"

	"go.mewis.me/meta-extra/pkg/messagix/lightspeed"
	"go.mewis.me/meta-extra/pkg/messagix/table"
)

func TestInsertXMAAttachmentLoggingTypeString(t *testing.T) {
	dependencies := table.SPToDepMap([]string{"insertXmaAttachment"})
	result := &table.LSTable{}
	decoder := lightspeed.NewLightSpeedDecoder(dependencies, result)
	args := make([]interface{}, 2+113)
	args[0] = float64(lightspeed.CALL_STORED_PROCEDURE)
	args[1] = "insertXmaAttachment"
	args[2+112] = "xma_open_contact_share"
	decoder.Decode(args)
	if len(result.LSInsertXmaAttachment) != 1 {
		t.Fatalf("unexpected XMA decode count: %d", len(result.LSInsertXmaAttachment))
	}
	if result.LSInsertXmaAttachment[0].AttachmentLoggingType != "xma_open_contact_share" {
		t.Fatalf("unexpected attachment logging type: %q", result.LSInsertXmaAttachment[0].AttachmentLoggingType)
	}
}
