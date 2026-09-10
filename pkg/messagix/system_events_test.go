package messagix

import (
	"testing"

	"go.mewis.me/meta-extra/pkg/messagix/table"
)

func TestThreadSystemEvents(t *testing.T) {
	client := newStateTestClient()
	client.applyMailboxState(&table.LSTable{LSAddPollForThread: []*table.LSAddPollForThread{{PollID: 90, ThreadKey: 10}}})
	tbl := &table.LSTable{
		LSUpdateThreadParticipantNicknameV2:  []*table.LSUpdateThreadParticipantNicknameV2{{ThreadKey: 10, ParticipantId: 20, Nickname: "nick"}},
		LSUpdateThreadCustomEmoji:            []*table.LSUpdateThreadCustomEmoji{{ThreadKey: 10, CustomEmoji: "x"}},
		LSUpdateThreadApprovalMode:           []*table.LSUpdateThreadApprovalMode{{ThreadKey: 10, Value: true}},
		LSUpdateThreadTheme:                  []*table.LSUpdateThreadTheme{{ThreadKey: 10}},
		LSAddParticipantIdToGroupThread:      []*table.LSAddParticipantIdToGroupThread{{ThreadKey: 10, ContactId: 21}},
		LSRemoveParticipantFromThread:        []*table.LSRemoveParticipantFromThread{{ThreadKey: 10, ParticipantId: 22}},
		LSUpdateThreadParticipantAdminStatus: []*table.LSUpdateThreadParticipantAdminStatus{{ThreadKey: 10, ContactId: 23, IsAdmin: true}},
		LSSetPinnedMessage:                   []*table.LSSetPinnedMessage{{ThreadKey: 10, MessageId: "m1", PinnedTimestampMs: 1}},
		LSAddPollOptionV2:                    []*table.LSAddPollOption{{PollID: 90, OptionID: 1}},
		LSAddPollVoteV2:                      []*table.LSAddPollVote{{PollID: 90, OptionID: 1, ContactID: 20}},
	}
	events := client.threadSystemEvents(tbl)
	if len(events) != 9 {
		t.Fatalf("unexpected event count: %d", len(events))
	}
	byKind := make(map[ThreadSystemEventKind]*ThreadSystemEvent, len(events))
	for _, event := range events {
		byKind[event.Kind] = event
	}
	if event := byKind[ThreadSystemNicknameUpdated]; event == nil || event.ParticipantID != 20 || event.Nickname != "nick" {
		t.Fatalf("unexpected nickname event: %#v", event)
	}
	if event := byKind[ThreadSystemEmojiUpdated]; event == nil || event.Emoji != "x" {
		t.Fatalf("unexpected emoji event: %#v", event)
	}
	if event := byKind[ThreadSystemApprovalUpdated]; event == nil || !event.Enabled {
		t.Fatalf("unexpected approval event: %#v", event)
	}
	if event := byKind[ThreadSystemThemeUpdated]; event == nil || event.ThreadKey != 10 {
		t.Fatalf("unexpected theme event: %#v", event)
	}
	if event := byKind[ThreadSystemMemberAdded]; event == nil || event.ParticipantID != 21 {
		t.Fatalf("unexpected member-added event: %#v", event)
	}
	if event := byKind[ThreadSystemMemberRemoved]; event == nil || event.ParticipantID != 22 {
		t.Fatalf("unexpected member-removed event: %#v", event)
	}
	if event := byKind[ThreadSystemAdminUpdated]; event == nil || event.ParticipantID != 23 || !event.IsAdmin {
		t.Fatalf("unexpected admin event: %#v", event)
	}
	if event := byKind[ThreadSystemPinUpdated]; event == nil || event.MessageID != "m1" || !event.Pinned {
		t.Fatalf("unexpected pin event: %#v", event)
	}
	if event := byKind[ThreadSystemPollUpdated]; event == nil || event.PollID != 90 || event.ThreadKey != 10 {
		t.Fatalf("unexpected poll event: %#v", event)
	}
}

func TestThreadSystemPinUnpinEvent(t *testing.T) {
	client := newStateTestClient()
	events := client.threadSystemEvents(&table.LSTable{LSSetPinnedMessage: []*table.LSSetPinnedMessage{{ThreadKey: 10, MessageId: "m1", PinnedTimestampMs: 0}}})
	if len(events) != 1 || events[0].Kind != ThreadSystemPinUpdated || events[0].Pinned {
		t.Fatalf("unexpected unpin event: %#v", events)
	}
}
