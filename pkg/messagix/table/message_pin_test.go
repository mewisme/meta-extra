package table_test

import (
	"testing"

	"go.mewis.me/meta-extra/pkg/messagix/table"
)

func TestSetPinnedMessageState(t *testing.T) {
	if !(&table.LSSetPinnedMessage{PinnedTimestampMs: 1}).IsPinned() {
		t.Fatal("positive pinned timestamp must be pinned")
	}
	if (&table.LSSetPinnedMessage{PinnedTimestampMs: 0}).IsPinned() {
		t.Fatal("zero pinned timestamp must be unpinned")
	}
}
