package table_test

import (
	"testing"

	"go.mewis.me/meta-extra/pkg/messagix/lightspeed"
	"go.mewis.me/meta-extra/pkg/messagix/table"
)

func TestMessageSearchDependencies(t *testing.T) {
	dependencies := table.SPToDepMap([]string{"updateMessageSearchQueryStatus", "insertMessageSearchResult"})
	if dependencies["updateMessageSearchQueryStatus"] != "LSUpdateMessageSearchQueryStatus" || dependencies["insertMessageSearchResult"] != "LSInsertMessageSearchResult" {
		t.Fatalf("unexpected dependency mappings: %#v", dependencies)
	}
	result := &table.LSTable{}
	decoder := lightspeed.NewLightSpeedDecoder(dependencies, result)
	decoder.Decode([]interface{}{float64(lightspeed.CALL_STORED_PROCEDURE), "updateMessageSearchQueryStatus", "needle", []interface{}{float64(lightspeed.I64_FROM_STRING), "2"}, []interface{}{float64(lightspeed.I64_FROM_STRING), "2"}, false, []interface{}{float64(lightspeed.UNDEFINED)}, []interface{}{float64(lightspeed.I64_FROM_STRING), "456"}, []interface{}{float64(lightspeed.I64_FROM_STRING), "1"}})
	decoder.Decode([]interface{}{float64(lightspeed.CALL_STORED_PROCEDURE), "insertMessageSearchResult", "needle", []interface{}{float64(lightspeed.I64_FROM_STRING), "0"}, []interface{}{float64(lightspeed.I64_FROM_STRING), "456"}, []interface{}{float64(lightspeed.I64_FROM_STRING), "2"}, []interface{}{float64(lightspeed.I64_FROM_STRING), "2"}, "Sender", "mid.1", []interface{}{float64(lightspeed.I64_FROM_STRING), "123456789"}, "hello needle", "https://example.invalid/avatar.jpg", []interface{}{float64(lightspeed.UNDEFINED)}, "6", "6", []interface{}{float64(lightspeed.I64_FROM_STRING), "4"}})
	if len(result.LSUpdateMessageSearchQueryStatus) != 1 || len(result.LSInsertMessageSearchResult) != 1 {
		t.Fatalf("unexpected decode counts: status=%d result=%d", len(result.LSUpdateMessageSearchQueryStatus), len(result.LSInsertMessageSearchResult))
	}
	status := result.LSUpdateMessageSearchQueryStatus[0]
	if status.Query != "needle" || status.Type_ != table.MessageSearchTypeMessage || status.Status != table.MessageSearchStatusComplete || status.HasNextPage || status.NextPageCursor != "" || status.ThreadKeyV2 != 456 || status.ResultCount != 1 {
		t.Fatalf("unexpected search status: %#v", status)
	}
	item := result.LSInsertMessageSearchResult[0]
	if item.Query != "needle" || item.GlobalIndex != 0 || item.ThreadKey != 456 || item.Type_ != table.MessageSearchTypeMessage || item.ThreadType != table.GROUP_THREAD || item.DisplayName != "Sender" || item.MessageId != "mid.1" || item.MessageTimestampMs != 123456789 || item.ContextLine != "hello needle" || item.MatchOffsets != "6" || item.MatchLengths != "6" {
		t.Fatalf("unexpected search result: %#v", item)
	}
	if value, ok := item.Unrecognized[13]; !ok || value != int64(4) {
		t.Fatalf("unknown current slot 13 was not preserved: %#v", item.Unrecognized)
	}
	if _, ok := item.Unrecognized[10]; ok {
		t.Fatalf("undefined slot 10 should not be preserved: %#v", item.Unrecognized)
	}
}

func TestMessageSearchEnumValues(t *testing.T) {
	if table.MessageSearchTypeThread != 1 || table.MessageSearchTypeMessage != 2 {
		t.Fatal("message search type values changed")
	}
	if table.MessageSearchStatusPending != 1 || table.MessageSearchStatusComplete != 2 || table.MessageSearchStatusFailed != 3 {
		t.Fatal("message search status values changed")
	}
}
