package messagix

import "go.mewis.me/meta-extra/pkg/messagix/table"

type ThreadSystemEventKind string

const (
	ThreadSystemNicknameUpdated ThreadSystemEventKind = "nickname_updated"
	ThreadSystemEmojiUpdated    ThreadSystemEventKind = "emoji_updated"
	ThreadSystemApprovalUpdated ThreadSystemEventKind = "approval_mode_updated"
	ThreadSystemThemeUpdated    ThreadSystemEventKind = "theme_updated"
	ThreadSystemMemberAdded     ThreadSystemEventKind = "member_added"
	ThreadSystemMemberRemoved   ThreadSystemEventKind = "member_removed"
	ThreadSystemAdminUpdated    ThreadSystemEventKind = "participant_admin_updated"
	ThreadSystemPinUpdated      ThreadSystemEventKind = "message_pin_updated"
	ThreadSystemPollUpdated     ThreadSystemEventKind = "poll_updated"
)

type ThreadSystemEvent struct {
	Kind          ThreadSystemEventKind `json:"kind"`
	ThreadKey     int64                 `json:"thread_key"`
	ParticipantID int64                 `json:"participant_id,omitempty"`
	MessageID     string                `json:"message_id,omitempty"`
	PollID        int64                 `json:"poll_id,omitempty"`
	Nickname      string                `json:"nickname,omitempty"`
	Emoji         string                `json:"emoji,omitempty"`
	Enabled       bool                  `json:"enabled,omitempty"`
	Pinned        bool                  `json:"pinned,omitempty"`
	IsAdmin       bool                  `json:"is_admin,omitempty"`
}

func (c *Client) threadSystemEvents(tbl *table.LSTable) []*ThreadSystemEvent {
	if tbl == nil {
		return nil
	}
	events := make([]*ThreadSystemEvent, 0)
	for _, item := range tbl.LSUpdateThreadParticipantNicknameV2 {
		if item != nil && item.ThreadKey > 0 {
			events = append(events, &ThreadSystemEvent{Kind: ThreadSystemNicknameUpdated, ThreadKey: item.ThreadKey, ParticipantID: item.ParticipantId, Nickname: item.Nickname})
		}
	}
	for _, item := range tbl.LSUpdateThreadCustomEmoji {
		if item != nil && item.ThreadKey > 0 {
			events = append(events, &ThreadSystemEvent{Kind: ThreadSystemEmojiUpdated, ThreadKey: item.ThreadKey, Emoji: item.CustomEmoji})
		}
	}
	for _, item := range tbl.LSUpdateThreadApprovalMode {
		if item != nil && item.ThreadKey > 0 {
			events = append(events, &ThreadSystemEvent{Kind: ThreadSystemApprovalUpdated, ThreadKey: item.ThreadKey, Enabled: item.Value})
		}
	}
	for _, item := range tbl.LSUpdateThreadTheme {
		if item != nil && item.ThreadKey > 0 {
			events = append(events, &ThreadSystemEvent{Kind: ThreadSystemThemeUpdated, ThreadKey: item.ThreadKey})
		}
	}
	for _, item := range tbl.LSAddParticipantIdToGroupThread {
		if item != nil && item.ThreadKey > 0 && item.ContactId > 0 {
			events = append(events, &ThreadSystemEvent{Kind: ThreadSystemMemberAdded, ThreadKey: item.ThreadKey, ParticipantID: item.ContactId})
		}
	}
	for _, item := range tbl.LSRemoveParticipantFromThread {
		if item != nil && item.ThreadKey > 0 && item.ParticipantId > 0 {
			events = append(events, &ThreadSystemEvent{Kind: ThreadSystemMemberRemoved, ThreadKey: item.ThreadKey, ParticipantID: item.ParticipantId})
		}
	}
	for _, item := range tbl.LSUpdateThreadParticipantAdminStatus {
		if item != nil && item.ThreadKey > 0 && item.ContactId > 0 {
			events = append(events, &ThreadSystemEvent{Kind: ThreadSystemAdminUpdated, ThreadKey: item.ThreadKey, ParticipantID: item.ContactId, IsAdmin: item.IsAdmin})
		}
	}
	for _, item := range tbl.LSSetPinnedMessage {
		if item != nil && item.ThreadKey > 0 && item.MessageId != "" {
			events = append(events, &ThreadSystemEvent{Kind: ThreadSystemPinUpdated, ThreadKey: item.ThreadKey, MessageID: item.MessageId, Pinned: item.IsPinned()})
		}
	}

	polls := make(map[int64]struct{})
	for _, item := range tbl.LSAddPollForThread {
		if item != nil && item.PollID > 0 {
			polls[item.PollID] = struct{}{}
		}
	}
	for _, items := range [][]*table.LSAddPollOption{tbl.LSAddPollOption, tbl.LSAddPollOptionV2} {
		for _, item := range items {
			if item != nil && item.PollID > 0 {
				polls[item.PollID] = struct{}{}
			}
		}
	}
	for _, items := range [][]*table.LSAddPollVote{tbl.LSAddPollVote, tbl.LSAddPollVoteV2, tbl.LSRemovePollVoteV2} {
		for _, item := range items {
			if item != nil && item.PollID > 0 {
				polls[item.PollID] = struct{}{}
			}
		}
	}
	for pollID := range polls {
		threadKey := int64(0)
		if details, _, err := c.pollDetailsSnapshot(pollID); err == nil {
			threadKey = details.ThreadKey
		}
		events = append(events, &ThreadSystemEvent{Kind: ThreadSystemPollUpdated, ThreadKey: threadKey, PollID: pollID})
	}
	return events
}
