package messagix

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"go.mewis.me/meta-extra/pkg/messagix/socket"
	"go.mewis.me/meta-extra/pkg/messagix/table"
)

var (
	ErrMailboxStateNotLoaded  = errors.New("messenger mailbox state is not loaded")
	ErrPollDetailsUnavailable = errors.New("messenger poll details are not available in mailbox state")
	ErrPollNotFound           = errors.New("messenger poll not found")
)

type PinnedMessage struct {
	ThreadKey         int64  `json:"thread_key"`
	MessageID         string `json:"message_id"`
	PinnedTimestampMS int64  `json:"pinned_timestamp_ms"`
	AuthorityLevel    int64  `json:"authority_level"`
}

type PollDetails struct {
	ID                           int64        `json:"id"`
	ThreadKey                    int64        `json:"thread_key"`
	Title                        string       `json:"title,omitempty"`
	LastUpdateMessageID          string       `json:"last_update_message_id"`
	LastUpdateMessageTimestampMS int64        `json:"last_update_message_timestamp_ms"`
	LastUpdateMessageEventType   int64        `json:"last_update_message_event_type"`
	Options                      []PollOption `json:"options"`
	Votes                        []PollVote   `json:"votes"`
}

type PollOption struct {
	ID                       int64  `json:"id"`
	Text                     string `json:"text"`
	SortKeyVotingTimestamp   int64  `json:"sort_key_voting_timestamp"`
	SortKeyCreationTimestamp int64  `json:"sort_key_creation_timestamp"`
}

type PollVote struct {
	OptionID    int64  `json:"option_id"`
	ContactID   int64  `json:"contact_id"`
	TimestampMS int64  `json:"timestamp_ms"`
	VoteCount   int64  `json:"vote_count"`
	ThreadKey   int64  `json:"thread_key"`
	MessageID   string `json:"message_id"`
}

type pollVoteKey struct {
	OptionID  int64
	ContactID int64
	MessageID string
}

type pollState struct {
	details  PollDetails
	options  map[int64]PollOption
	votes    map[pollVoteKey]PollVote
	revision uint64
}

func (c *Client) resetMailboxState() {
	if c == nil {
		return
	}
	c.stateMu.Lock()
	clear(c.pinnedMessages)
	clear(c.polls)
	clear(c.pollsV2Threads)
	clear(c.messageSearches)
	c.stateMu.Unlock()
	c.mailboxStateLoaded.Store(false)
}

func (c *Client) applyMailboxState(tbl *table.LSTable) {
	if c == nil || tbl == nil {
		return
	}
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.applyMessageSearchStateLocked(tbl)
	for _, item := range tbl.LSClearPinnedMessages {
		if item != nil {
			delete(c.pinnedMessages, item.ThreadKey)
		}
	}
	for _, item := range tbl.LSSetPinnedMessage {
		if item == nil || item.ThreadKey <= 0 || item.MessageId == "" {
			continue
		}
		messages := c.pinnedMessages[item.ThreadKey]
		if item.IsPinned() {
			if messages == nil {
				messages = make(map[string]PinnedMessage)
				c.pinnedMessages[item.ThreadKey] = messages
			}
			messages[item.MessageId] = PinnedMessage{ThreadKey: item.ThreadKey, MessageID: item.MessageId, PinnedTimestampMS: item.PinnedTimestampMs, AuthorityLevel: item.AuthorityLevel}
		} else if messages != nil {
			delete(messages, item.MessageId)
			if len(messages) == 0 {
				delete(c.pinnedMessages, item.ThreadKey)
			}
		}
	}
	for _, item := range tbl.LSWriteThreadCapabilities {
		if item != nil && item.ThreadKey > 0 {
			c.pollsV2Threads[item.ThreadKey] = item.Capabilities3&(int64(1)<<31) != 0
		}
	}
	for _, item := range tbl.LSApplyAdminMessageCTAV2 {
		if item == nil || item.CTAType != "admin_msg_poll_details" || item.PollID <= 0 {
			continue
		}
		state := c.ensurePollState(item.PollID)
		state.details.ID = item.PollID
		state.details.ThreadKey = item.ThreadKey
		state.details.Title = item.PollTitle
		state.details.LastUpdateMessageID = item.MessageID
		state.details.LastUpdateMessageTimestampMS = item.TimestampMS
		state.revision++
	}
	for _, item := range tbl.LSAddPollForThread {
		if item == nil || item.PollID <= 0 {
			continue
		}
		state := c.ensurePollState(item.PollID)
		state.details.ID = item.PollID
		state.details.ThreadKey = item.ThreadKey
		state.details.LastUpdateMessageID = item.LastUpdateMessageID
		state.details.LastUpdateMessageTimestampMS = item.LastUpdateMessageTimestampMS
		state.details.LastUpdateMessageEventType = item.LastUpdateMessageEventType
		state.revision++
	}
	for _, items := range [][]*table.LSAddPollOption{tbl.LSAddPollOption, tbl.LSAddPollOptionV2} {
		for _, item := range items {
			if item == nil || item.PollID <= 0 || item.OptionID <= 0 {
				continue
			}
			state := c.ensurePollState(item.PollID)
			state.options[item.OptionID] = PollOption{ID: item.OptionID, Text: item.OptionText, SortKeyVotingTimestamp: item.SortKeyVotingTimestamp, SortKeyCreationTimestamp: item.SortKeyCreationTimestamp}
			state.revision++
		}
	}
	for _, items := range [][]*table.LSAddPollVote{tbl.LSAddPollVote, tbl.LSAddPollVoteV2} {
		for _, item := range items {
			if item == nil || item.PollID <= 0 || item.OptionID <= 0 {
				continue
			}
			state := c.ensurePollState(item.PollID)
			vote := PollVote{OptionID: item.OptionID, ContactID: item.ContactID, TimestampMS: item.TimestampMS, VoteCount: item.VoteCount, ThreadKey: item.ThreadKey, MessageID: item.MessageID}
			key := pollVoteKey{OptionID: item.OptionID, ContactID: item.ContactID, MessageID: item.MessageID}
			if item.TimestampMS <= 0 && item.VoteCount <= 0 {
				delete(state.votes, key)
			} else {
				state.votes[key] = vote
			}
			state.revision++
		}
	}
	for _, item := range tbl.LSRemovePollVoteV2 {
		if item == nil || item.PollID <= 0 || item.OptionID <= 0 || item.ContactID <= 0 {
			continue
		}
		state := c.ensurePollState(item.PollID)
		for key := range state.votes {
			if key.OptionID == item.OptionID && key.ContactID == item.ContactID {
				delete(state.votes, key)
			}
		}
		state.revision++
	}
}

func (c *Client) ensurePollState(pollID int64) *pollState {
	state := c.polls[pollID]
	if state == nil {
		state = &pollState{details: PollDetails{ID: pollID}, options: make(map[int64]PollOption), votes: make(map[pollVoteKey]PollVote)}
		c.polls[pollID] = state
	}
	return state
}

func (c *Client) ListPinnedMessages(threadKey int64) ([]PinnedMessage, error) {
	if c == nil {
		return nil, ErrClientIsNil
	} else if threadKey <= 0 {
		return nil, fmt.Errorf("thread key must be positive")
	} else if !c.mailboxStateLoaded.Load() {
		return nil, ErrMailboxStateNotLoaded
	}
	c.stateMu.RLock()
	messages := make([]PinnedMessage, 0, len(c.pinnedMessages[threadKey]))
	for _, item := range c.pinnedMessages[threadKey] {
		messages = append(messages, item)
	}
	c.stateMu.RUnlock()
	sort.Slice(messages, func(i, j int) bool {
		if messages[i].PinnedTimestampMS == messages[j].PinnedTimestampMS {
			return messages[i].MessageID < messages[j].MessageID
		}
		return messages[i].PinnedTimestampMS > messages[j].PinnedTimestampMS
	})
	return messages, nil
}

func (c *Client) GetPollDetails(pollID int64) (*PollDetails, error) {
	if c == nil {
		return nil, ErrClientIsNil
	} else if pollID <= 0 {
		return nil, fmt.Errorf("poll ID must be positive")
	} else if !c.mailboxStateLoaded.Load() {
		return nil, ErrMailboxStateNotLoaded
	}
	details, _, err := c.pollDetailsSnapshot(pollID)
	return details, err
}

func (c *Client) FetchPollDetails(ctx context.Context, pollID int64) (*PollDetails, error) {
	if c == nil {
		return nil, ErrClientIsNil
	} else if pollID <= 0 {
		return nil, fmt.Errorf("poll ID must be positive")
	} else if !c.mailboxStateLoaded.Load() {
		return nil, ErrMailboxStateNotLoaded
	}
	details, baseline, stateErr := c.pollDetailsSnapshot(pollID)
	if stateErr == nil {
		enabled, known := c.pollV2Enabled(details.ThreadKey)
		if known && !enabled {
			if len(details.Options) == 0 {
				return nil, ErrPollDetailsUnavailable
			}
			return details, nil
		}
	}
	if _, err := c.ExecuteTasks(ctx, &socket.FetchPollDetailsTask{PollID: pollID}); err != nil {
		return nil, fmt.Errorf("fetch poll details: %w", err)
	}
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for {
		details, revision, err := c.pollDetailsSnapshot(pollID)
		if err == nil && revision > baseline && len(details.Options) > 0 {
			return details, nil
		} else if err != nil && !errors.Is(err, ErrPollNotFound) {
			return nil, err
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
	}
}

func (c *Client) pollV2Enabled(threadKey int64) (bool, bool) {
	c.stateMu.RLock()
	enabled, ok := c.pollsV2Threads[threadKey]
	c.stateMu.RUnlock()
	return enabled, ok
}

func (c *Client) pollDetailsSnapshot(pollID int64) (*PollDetails, uint64, error) {
	c.stateMu.RLock()
	state := c.polls[pollID]
	if state == nil {
		c.stateMu.RUnlock()
		return nil, 0, ErrPollNotFound
	}
	details := state.details
	revision := state.revision
	details.Options = make([]PollOption, 0, len(state.options))
	for _, item := range state.options {
		details.Options = append(details.Options, item)
	}
	details.Votes = make([]PollVote, 0, len(state.votes))
	for _, item := range state.votes {
		details.Votes = append(details.Votes, item)
	}
	c.stateMu.RUnlock()
	sort.Slice(details.Options, func(i, j int) bool {
		if details.Options[i].SortKeyCreationTimestamp == details.Options[j].SortKeyCreationTimestamp {
			return details.Options[i].ID < details.Options[j].ID
		}
		return details.Options[i].SortKeyCreationTimestamp < details.Options[j].SortKeyCreationTimestamp
	})
	sort.Slice(details.Votes, func(i, j int) bool {
		if details.Votes[i].TimestampMS == details.Votes[j].TimestampMS {
			if details.Votes[i].OptionID == details.Votes[j].OptionID {
				return details.Votes[i].ContactID < details.Votes[j].ContactID
			}
			return details.Votes[i].OptionID < details.Votes[j].OptionID
		}
		return details.Votes[i].TimestampMS < details.Votes[j].TimestampMS
	})
	return &details, revision, nil
}
