package messagix

import (
	"errors"
	"testing"

	"go.mewis.me/meta-extra/pkg/messagix/table"
)

func TestMessageSearchValidation(t *testing.T) {
	client := newStateTestClient()
	client.mailboxStateLoaded.Store(true)
	if _, err := client.SearchMessages(t.Context(), MessageSearchRequest{ThreadKey: 1, Query: "   "}); err == nil {
		t.Fatal("expected empty-query validation error")
	}
	empty := ""
	if _, err := client.SearchMessages(t.Context(), MessageSearchRequest{ThreadKey: 1, Query: "needle", Cursor: &empty}); err == nil {
		t.Fatal("expected empty-cursor validation error")
	}
	if _, err := client.SearchMessages(t.Context(), MessageSearchRequest{ThreadKey: 0, Query: "needle"}); err == nil {
		t.Fatal("expected thread-key validation error")
	}
}

func TestMessageSearchStatePagination(t *testing.T) {
	client := newStateTestClient()
	key := messageSearchKey{Type: table.MessageSearchTypeMessage, Query: "needle", ThreadKey: 7}
	baselineRevision, baselineResults := client.prepareMessageSearch(key, true)
	client.applyMailboxState(&table.LSTable{
		LSUpdateMessageSearchQueryStatus: []*table.LSUpdateMessageSearchQueryStatus{{Query: "needle", Type_: table.MessageSearchTypeMessage, Status: table.MessageSearchStatusComplete, HasNextPage: true, NextPageCursor: "cursor-1", ThreadKeyV2: 7, ResultCount: 2}},
		LSInsertMessageSearchResult:      []*table.LSInsertMessageSearchResult{{Query: "needle", GlobalIndex: 0, ThreadKey: 7, Type_: table.MessageSearchTypeMessage, ThreadType: table.GROUP_THREAD, DisplayName: "Alice", MessageId: "m1", MessageTimestampMs: 100, ContextLine: "first needle", ProfilePicUrl: "avatar", MatchOffsets: "6", MatchLengths: "6"}},
	})
	page, status, revision, err := client.messageSearchSnapshot(key, baselineResults)
	if err != nil {
		t.Fatal(err)
	}
	if revision <= baselineRevision || status != table.MessageSearchStatusComplete || len(page.Results) != 1 || page.Results[0].MessageID != "m1" || page.NextCursor == nil || *page.NextCursor != "cursor-1" || !page.HasNextPage || page.ResultCount != 2 {
		t.Fatalf("unexpected first page: page=%#v status=%d revision=%d baseline=%d", page, status, revision, baselineRevision)
	}
	if len(page.Results[0].Highlights) != 1 || page.Results[0].Highlights[0].Offset != 6 || page.Results[0].Highlights[0].Length != 6 {
		t.Fatalf("unexpected highlights: %#v", page.Results[0].Highlights)
	}
	baselineRevision, baselineResults = client.prepareMessageSearch(key, false)
	client.applyMailboxState(&table.LSTable{
		LSUpdateMessageSearchQueryStatus: []*table.LSUpdateMessageSearchQueryStatus{{Query: "needle", Type_: table.MessageSearchTypeMessage, Status: table.MessageSearchStatusComplete, ThreadKeyV2: 7, ResultCount: 2}},
		LSInsertMessageSearchResult:      []*table.LSInsertMessageSearchResult{{Query: "needle", GlobalIndex: 1, ThreadKey: 7, Type_: table.MessageSearchTypeMessage, ThreadType: table.GROUP_THREAD, DisplayName: "Bob", MessageId: "m2", MessageTimestampMs: 200, ContextLine: "second needle"}},
	})
	page, status, revision, err = client.messageSearchSnapshot(key, baselineResults)
	if err != nil {
		t.Fatal(err)
	}
	if revision <= baselineRevision || status != table.MessageSearchStatusComplete || len(page.Results) != 1 || page.Results[0].MessageID != "m2" || page.HasNextPage || page.NextCursor != nil {
		t.Fatalf("unexpected pagination page: %#v", page)
	}
	_, baselineResults = client.prepareMessageSearch(key, true)
	if len(baselineResults) != 0 {
		t.Fatalf("initial search did not clear previous results: %#v", baselineResults)
	}
}

func TestMessageSearchStatusErrors(t *testing.T) {
	if err := messageSearchStatusError(table.MessageSearchStatusComplete, false); err != nil {
		t.Fatalf("complete status failed: %v", err)
	}
	if err := messageSearchStatusError(table.MessageSearchStatusPending, false); err != nil {
		t.Fatalf("pending status failed: %v", err)
	}
	if err := messageSearchStatusError(table.MessageSearchStatusFailed, false); !errors.Is(err, ErrMessageSearchFailed) {
		t.Fatalf("expected search failure, got %v", err)
	}
	if err := messageSearchStatusError(table.MessageSearchStatusFailed, true); !errors.Is(err, ErrMessageSearchCursorRejected) {
		t.Fatalf("expected cursor rejection, got %v", err)
	}
	if err := messageSearchStatusError(table.MessageSearchStatus(99), false); !errors.Is(err, ErrMessageSearchProtocol) {
		t.Fatalf("expected protocol mismatch, got %v", err)
	}
}

func TestParseMessageSearchHighlights(t *testing.T) {
	highlights := parseMessageSearchHighlights("19,5,0,11", "19,5,4,7")
	if len(highlights) != 4 || highlights[0].Offset != 19 || highlights[2].Length != 4 {
		t.Fatalf("unexpected highlights: %#v", highlights)
	}
	if parseMessageSearchHighlights("1,2", "3") != nil || parseMessageSearchHighlights("bad", "1") != nil {
		t.Fatal("invalid highlight data should not be exposed")
	}
}
