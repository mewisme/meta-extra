package socket

import (
	"encoding/json"
	"time"

	"go.mau.fi/util/exerrors"

	"go.mewis.me/meta-extra/pkg/messagix/table"
)

type SendReactionTask struct {
	ThreadKey       int64                  `json:"thread_key,omitempty"`
	TimestampMs     int64                  `json:"timestamp_ms"`
	MessageID       string                 `json:"message_id"`
	ActorID         int64                  `json:"actor_id"`
	Reaction        string                 `json:"reaction"` // unicode emoji (empty reaction to remove)
	ReactionStyle   any                    `json:"reaction_style"`
	SyncGroup       int                    `json:"sync_group"`
	SendAttribution table.ThreadSourceType `json:"send_attribution"`
}

func (t *SendReactionTask) GetLabel() string {
	return TaskLabels["SendReactionTask"]
}

func (t *SendReactionTask) Create() (any, string) {
	t.TimestampMs = time.Now().UnixMilli()
	t.SyncGroup = 1
	queueName := exerrors.Must(json.Marshal([]string{"reaction", t.MessageID}))
	return t, string(queueName)
}

type DeleteMessageTask struct {
	MessageId string `json:"message_id"`
}

func (t *DeleteMessageTask) GetLabel() string {
	return TaskLabels["DeleteMessageTask"]
}

func (t *DeleteMessageTask) Create() (any, string) {
	queueName := "unsend_message"
	return t, queueName
}

type DeleteMessageMeOnlyTask struct {
	ThreadKey int64  `json:"thread_key,omitempty"`
	MessageId string `json:"message_id"`
}

func (t *DeleteMessageMeOnlyTask) GetLabel() string {
	return TaskLabels["DeleteMessageMeOnlyTask"]
}

func (t *DeleteMessageMeOnlyTask) Create() (any, string) {
	queueName := "155"
	return t, queueName
}

type FetchReactionsV2UserList struct {
	ThreadID     int64   `json:"thread_id"`
	MessageID    string  `json:"message_id"`
	ReactionFBID *int64  `json:"reaction_fbid"`
	Cursor       *string `json:"cursor"`
	SyncGroup    int64   `json:"sync_group"`
}

func (t *FetchReactionsV2UserList) GetLabel() string {
	return TaskLabels["FetchReactionsV2UserList"]
}

func (t *FetchReactionsV2UserList) Create() (any, string) {
	return t, "fetch_reactions_v2_details_users_list"
}

type SendReactionV2Task struct {
	ThreadID         int64  `json:"thread_id"`
	MessageID        string `json:"message_id"`
	MessageTimestamp int64  `json:"message_timestamp"`
	ActorID          int64  `json:"actor_id"`
	ReactionFBID     int64  `json:"reaction_fbid"`
	ReactionStyle    int    `json:"reaction_style"` // 1
	CurrentCount     int    `json:"current_count"`
	ViewerIsReactor  int    `json:"viewer_is_reactor"` // 1 if adding reaction, 0 if removing reaction
	Operation        int    `json:"operation"`         // 1 for add, 3 for remove
	ReactionLiteral  string `json:"reaction_literal"`  // unicode emoji
	EntryPoint       any    `json:"entry_point"`       // null
	SyncGroup        int    `json:"sync_group"`        // 104
}

func (t *SendReactionV2Task) GetLabel() string {
	return TaskLabels["SendReactionV2"]
}

func (t *SendReactionV2Task) Create() (any, string) {
	return t, string(exerrors.Must(json.Marshal([]string{"reaction_v2", t.MessageID})))
}

type SearchMessagesTask struct {
	Query          string
	ThreadKey      int64
	NextPageCursor *string
}

type searchMessagesPayload struct {
	Query          string                  `json:"query"`
	Type           table.MessageSearchType `json:"type"`
	ThreadKey      int64                   `json:"thread_key"`
	NextPageCursor *string                 `json:"next_page_cursor"`
	ClientCallerID string                  `json:"client_caller_id"`
}

func (t *SearchMessagesTask) GetLabel() string { return TaskLabels["SearchMessagesTask"] }
func (t *SearchMessagesTask) Create() (any, string) {
	callerID := "msgr_search_web_start_conversation_message_search"
	if t.NextPageCursor != nil {
		callerID = "msgr_search_web_in_conversation_message_search"
	}
	return &searchMessagesPayload{Query: t.Query, Type: table.MessageSearchTypeMessage, ThreadKey: t.ThreadKey, NextPageCursor: t.NextPageCursor, ClientCallerID: callerID}, "message_search"
}
