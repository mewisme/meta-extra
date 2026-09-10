package table_test

import (
	"testing"

	"go.mewis.me/meta-extra/pkg/messagix/table"
)

func TestContactRestrictionState(t *testing.T) {
	for _, test := range []struct {
		capabilities2 int64
		restricted    bool
	}{{0, false}, {1 << 2, true}, {1 << 3, false}, {(1 << 2) | (1 << 9), true}} {
		verified := &table.LSVerifyContactRowExists{Capabilities2: test.capabilities2}
		inserted := &table.LSDeleteThenInsertContact{Capabilities2: test.capabilities2}
		if verified.IsMessengerRestricted() != test.restricted || inserted.IsMessengerRestricted() != test.restricted {
			t.Fatalf("capabilities2=%d restriction mismatch", test.capabilities2)
		}
	}
}

func TestBlockedByViewerStatusValues(t *testing.T) {
	if table.BlockedByViewerStatusUnblocked != 0 || table.BlockedByViewerStatusMessageBlocked != 1 || table.BlockedByViewerStatusFullyBlocked != 2 {
		t.Fatal("blocked-by-viewer status values changed")
	}
}
