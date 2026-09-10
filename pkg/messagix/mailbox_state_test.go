package messagix

import (
	"context"
	"errors"
	"testing"

	"github.com/rs/zerolog"

	"go.mewis.me/meta-extra/pkg/messagix/cookies"
	"go.mewis.me/meta-extra/pkg/messagix/table"
	"go.mewis.me/meta-extra/pkg/messagix/types"
)

func newStateTestClient() *Client {
	return NewClient(&cookies.Cookies{Platform: types.Facebook}, zerolog.Nop(), &Config{})
}

func TestPinnedMessageState(t *testing.T) {
	client := newStateTestClient()
	if _, err := client.ListPinnedMessages(1); !errors.Is(err, ErrMailboxStateNotLoaded) {
		t.Fatalf("expected not loaded error, got %v", err)
	}
	client.mailboxStateLoaded.Store(true)
	client.applyMailboxState(&table.LSTable{LSSetPinnedMessage: []*table.LSSetPinnedMessage{
		{ThreadKey: 1, MessageId: "older", PinnedTimestampMs: 10, AuthorityLevel: 1},
		{ThreadKey: 1, MessageId: "newer", PinnedTimestampMs: 20, AuthorityLevel: 2},
	}})
	pins, err := client.ListPinnedMessages(1)
	if err != nil {
		t.Fatal(err)
	} else if len(pins) != 2 || pins[0].MessageID != "newer" || pins[1].MessageID != "older" {
		t.Fatalf("unexpected pins: %#v", pins)
	}
	client.applyMailboxState(&table.LSTable{LSSetPinnedMessage: []*table.LSSetPinnedMessage{{ThreadKey: 1, MessageId: "newer", PinnedTimestampMs: 0}}})
	pins, err = client.ListPinnedMessages(1)
	if err != nil || len(pins) != 1 || pins[0].MessageID != "older" {
		t.Fatalf("unexpected pins after unpin: %#v, %v", pins, err)
	}
	client.applyMailboxState(&table.LSTable{LSClearPinnedMessages: []*table.LSClearPinnedMessages{{ThreadKey: 1}}})
	pins, err = client.ListPinnedMessages(1)
	if err != nil || len(pins) != 0 {
		t.Fatalf("unexpected pins after clear: %#v, %v", pins, err)
	}
}

func TestPollState(t *testing.T) {
	client := newStateTestClient()
	client.mailboxStateLoaded.Store(true)
	client.applyMailboxState(&table.LSTable{
		LSAddPollForThread: []*table.LSAddPollForThread{{PollID: 7, ThreadKey: 3, LastUpdateMessageID: "m1", LastUpdateMessageTimestampMS: 100, LastUpdateMessageEventType: 4}},
		LSAddPollOption:    []*table.LSAddPollOption{{OptionID: 12, PollID: 7, OptionText: "B", SortKeyCreationTimestamp: 20}, {OptionID: 11, PollID: 7, OptionText: "A", SortKeyCreationTimestamp: 10}},
		LSAddPollVoteV2:    []*table.LSAddPollVote{{OptionID: 11, PollID: 7, ContactID: 9, TimestampMS: 30, VoteCount: 1, ThreadKey: 3, MessageID: "m1"}},
	})
	details, err := client.GetPollDetails(7)
	if err != nil {
		t.Fatal(err)
	} else if details.ThreadKey != 3 || details.LastUpdateMessageID != "m1" || len(details.Options) != 2 || details.Options[0].ID != 11 || len(details.Votes) != 1 || details.Votes[0].ContactID != 9 {
		t.Fatalf("unexpected poll details: %#v", details)
	}
	client.applyMailboxState(&table.LSTable{LSRemovePollVoteV2: []*table.LSAddPollVote{{OptionID: 11, PollID: 7, ContactID: 9, ThreadKey: 3, MessageID: "m1"}}})
	details, err = client.GetPollDetails(7)
	if err != nil || len(details.Votes) != 0 {
		t.Fatalf("unexpected poll votes after removal: %#v, %v", details, err)
	}
}

func TestPollStateFromAdminCTA(t *testing.T) {
	client := newStateTestClient()
	client.mailboxStateLoaded.Store(true)
	client.applyMailboxState(&table.LSTable{LSApplyAdminMessageCTAV2: []*table.LSApplyAdminMessageCTAV2{{ThreadKey: 3, TimestampMS: 100, MessageID: "m1", CTAType: "admin_msg_poll_details", PollID: 7, PollTitle: "Question"}}})
	details, err := client.GetPollDetails(7)
	if err != nil {
		t.Fatal(err)
	} else if details.ThreadKey != 3 || details.Title != "Question" || details.LastUpdateMessageID != "m1" {
		t.Fatalf("unexpected poll details: %#v", details)
	}
}

func TestFetchPollDetailsUsesLoadedV1State(t *testing.T) {
	client := newStateTestClient()
	client.mailboxStateLoaded.Store(true)
	client.applyMailboxState(&table.LSTable{
		LSWriteThreadCapabilities: []*table.LSWriteThreadCapabilities{{ThreadKey: 3, Capabilities3: 0}},
		LSAddPollForThread:        []*table.LSAddPollForThread{{PollID: 7, ThreadKey: 3}},
		LSAddPollOption:           []*table.LSAddPollOption{{OptionID: 11, PollID: 7, OptionText: "A"}},
	})
	details, err := client.FetchPollDetails(context.Background(), 7)
	if err != nil {
		t.Fatal(err)
	} else if len(details.Options) != 1 || details.Options[0].ID != 11 {
		t.Fatalf("unexpected poll details: %#v", details)
	}
}

func TestFetchPollDetailsRejectsIncompleteV1State(t *testing.T) {
	client := newStateTestClient()
	client.mailboxStateLoaded.Store(true)
	client.applyMailboxState(&table.LSTable{
		LSWriteThreadCapabilities: []*table.LSWriteThreadCapabilities{{ThreadKey: 3, Capabilities3: 0}},
		LSAddPollForThread:        []*table.LSAddPollForThread{{PollID: 7, ThreadKey: 3}},
	})
	if _, err := client.FetchPollDetails(context.Background(), 7); !errors.Is(err, ErrPollDetailsUnavailable) {
		t.Fatalf("expected unavailable error, got %v", err)
	}
}

func TestPollV2CapabilityBit(t *testing.T) {
	client := newStateTestClient()
	client.applyMailboxState(&table.LSTable{LSWriteThreadCapabilities: []*table.LSWriteThreadCapabilities{{ThreadKey: 3, Capabilities3: int64(1) << 31}}})
	if enabled, known := client.pollV2Enabled(3); !known || !enabled {
		t.Fatalf("unexpected capability: enabled=%t known=%t", enabled, known)
	}
}
