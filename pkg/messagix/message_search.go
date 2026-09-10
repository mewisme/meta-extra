package messagix

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.mewis.me/meta-extra/pkg/messagix/socket"
	"go.mewis.me/meta-extra/pkg/messagix/table"
)

var (
	ErrMessageSearchFailed         = errors.New("messenger message search failed")
	ErrMessageSearchCursorRejected = errors.New("messenger message search cursor was rejected")
	ErrMessageSearchProtocol       = errors.New("messenger message search protocol mismatch")
)

type MessageSearchRequest struct {
	ThreadKey int64
	Query     string
	Cursor    *string
}

type MessageSearchHighlight struct {
	Offset int `json:"offset"`
	Length int `json:"length"`
}

type MessageSearchResult struct {
	MessageID       string                   `json:"message_id"`
	ThreadKey       int64                    `json:"thread_key"`
	ThreadType      table.ThreadType         `json:"thread_type"`
	GlobalIndex     int64                    `json:"global_index"`
	SenderName      string                   `json:"sender_name"`
	SenderAvatarURL string                   `json:"sender_avatar_url,omitempty"`
	TimestampMS     int64                    `json:"timestamp_ms"`
	Text            string                   `json:"text"`
	Highlights      []MessageSearchHighlight `json:"highlights,omitempty"`
}

type MessageSearchPage struct {
	Results     []MessageSearchResult `json:"results"`
	ResultCount int64                 `json:"result_count"`
	HasNextPage bool                  `json:"has_next_page"`
	NextCursor  *string               `json:"next_cursor,omitempty"`
}

type messageSearchKey struct {
	Type      table.MessageSearchType
	Query     string
	ThreadKey int64
}

type messageSearchState struct {
	status         table.LSUpdateMessageSearchQueryStatus
	hasStatus      bool
	statusRevision uint64
	results        map[int64]MessageSearchResult
}

func (c *Client) SearchMessages(ctx context.Context, req MessageSearchRequest) (*MessageSearchPage, error) {
	if c == nil {
		return nil, ErrClientIsNil
	} else if req.ThreadKey <= 0 {
		return nil, fmt.Errorf("thread key must be positive")
	} else if strings.TrimSpace(req.Query) == "" {
		return nil, fmt.Errorf("message search query must not be empty")
	} else if req.Cursor != nil && *req.Cursor == "" {
		return nil, fmt.Errorf("message search cursor must not be empty")
	} else if !c.mailboxStateLoaded.Load() {
		return nil, ErrMailboxStateNotLoaded
	}
	key := messageSearchKey{Type: table.MessageSearchTypeMessage, Query: req.Query, ThreadKey: req.ThreadKey}
	baselineRevision, baselineResults := c.prepareMessageSearch(key, req.Cursor == nil)
	if _, err := c.ExecuteTasks(ctx, &socket.SearchMessagesTask{Query: req.Query, ThreadKey: req.ThreadKey, NextPageCursor: req.Cursor}); err != nil {
		return nil, fmt.Errorf("message search task: %w", err)
	}
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for {
		page, status, revision, err := c.messageSearchSnapshot(key, baselineResults)
		if err != nil {
			return nil, err
		}
		if revision > baselineRevision {
			if err := messageSearchStatusError(status, req.Cursor != nil); err != nil {
				return nil, err
			} else if status == table.MessageSearchStatusComplete {
				if page.HasNextPage && page.NextCursor == nil {
					return nil, fmt.Errorf("%w: response says another page exists but has no cursor", ErrMessageSearchProtocol)
				}
				return page, nil
			}
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
	}
}

func messageSearchStatusError(status table.MessageSearchStatus, hasCursor bool) error {
	switch status {
	case table.MessageSearchStatusPending, table.MessageSearchStatusComplete:
		return nil
	case table.MessageSearchStatusFailed:
		if hasCursor {
			return ErrMessageSearchCursorRejected
		}
		return ErrMessageSearchFailed
	default:
		return fmt.Errorf("%w: unsupported status %d", ErrMessageSearchProtocol, status)
	}
}

func (c *Client) prepareMessageSearch(key messageSearchKey, clearResults bool) (uint64, map[int64]struct{}) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	if clearResults {
		delete(c.messageSearches, key)
	}
	state := c.ensureMessageSearchState(key)
	baseline := make(map[int64]struct{}, len(state.results))
	for index := range state.results {
		baseline[index] = struct{}{}
	}
	return state.statusRevision, baseline
}

func (c *Client) ensureMessageSearchState(key messageSearchKey) *messageSearchState {
	state := c.messageSearches[key]
	if state == nil {
		state = &messageSearchState{results: make(map[int64]MessageSearchResult)}
		c.messageSearches[key] = state
	}
	return state
}

func (c *Client) applyMessageSearchStateLocked(tbl *table.LSTable) {
	for _, item := range tbl.LSUpdateMessageSearchQueryStatus {
		if item == nil || item.Type_ != table.MessageSearchTypeMessage || item.ThreadKeyV2 <= 0 || item.Query == "" {
			continue
		}
		key := messageSearchKey{Type: item.Type_, Query: item.Query, ThreadKey: item.ThreadKeyV2}
		state := c.ensureMessageSearchState(key)
		state.status = *item
		state.hasStatus = true
		state.statusRevision++
	}
	for _, item := range tbl.LSInsertMessageSearchResult {
		if item == nil || item.Type_ != table.MessageSearchTypeMessage || item.ThreadKey <= 0 || item.Query == "" || item.MessageId == "" {
			continue
		}
		key := messageSearchKey{Type: item.Type_, Query: item.Query, ThreadKey: item.ThreadKey}
		state := c.ensureMessageSearchState(key)
		state.results[item.GlobalIndex] = MessageSearchResult{MessageID: item.MessageId, ThreadKey: item.ThreadKey, ThreadType: item.ThreadType, GlobalIndex: item.GlobalIndex, SenderName: item.DisplayName, SenderAvatarURL: item.ProfilePicUrl, TimestampMS: item.MessageTimestampMs, Text: item.ContextLine, Highlights: parseMessageSearchHighlights(item.MatchOffsets, item.MatchLengths)}
	}
}

func (c *Client) messageSearchSnapshot(key messageSearchKey, baseline map[int64]struct{}) (*MessageSearchPage, table.MessageSearchStatus, uint64, error) {
	c.stateMu.RLock()
	state := c.messageSearches[key]
	if state == nil || !state.hasStatus {
		c.stateMu.RUnlock()
		return &MessageSearchPage{}, 0, 0, nil
	}
	status := state.status
	revision := state.statusRevision
	results := make([]MessageSearchResult, 0, len(state.results))
	for index, item := range state.results {
		if _, existed := baseline[index]; !existed {
			results = append(results, item)
		}
	}
	c.stateMu.RUnlock()
	sort.Slice(results, func(i, j int) bool { return results[i].GlobalIndex < results[j].GlobalIndex })
	page := &MessageSearchPage{Results: results, ResultCount: status.ResultCount, HasNextPage: status.HasNextPage}
	if status.NextPageCursor != "" {
		page.NextCursor = &status.NextPageCursor
	}
	return page, status.Status, revision, nil
}

func parseMessageSearchHighlights(offsets, lengths string) []MessageSearchHighlight {
	if offsets == "" || lengths == "" {
		return nil
	}
	offsetParts, lengthParts := strings.Split(offsets, ","), strings.Split(lengths, ",")
	if len(offsetParts) != len(lengthParts) {
		return nil
	}
	result := make([]MessageSearchHighlight, 0, len(offsetParts))
	for i := range offsetParts {
		offset, offsetErr := strconv.Atoi(offsetParts[i])
		length, lengthErr := strconv.Atoi(lengthParts[i])
		if offsetErr != nil || lengthErr != nil {
			return nil
		}
		result = append(result, MessageSearchHighlight{Offset: offset, Length: length})
	}
	return result
}
