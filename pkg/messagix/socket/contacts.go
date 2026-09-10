package socket

import (
	"encoding/json"
	"fmt"
	"time"

	"go.mau.fi/util/exerrors"

	"go.mewis.me/meta-extra/pkg/messagix/methods"
	"go.mewis.me/meta-extra/pkg/messagix/table"
)

type GetContactsTask struct {
	Limit int64 `json:"limit,omitempty"`
}

func (t *GetContactsTask) GetLabel() string {
	return TaskLabels["GetContactsTask"]
}

func (t *GetContactsTask) Create() (any, string) {
	queueName := exerrors.Must(json.Marshal([]string{"search_contacts", methods.GenerateTimestampString()}))
	return t, string(queueName)
}

type GetContactsFullTask struct {
	ContactID int64 `json:"contact_id"`
}

func (t *GetContactsFullTask) GetLabel() string {
	return TaskLabels["GetContactsFullTask"]
}

func (t *GetContactsFullTask) Create() (any, string) {
	queueName := "cpq_v2"
	return t, queueName
}

type SearchUserTask struct {
	Query                string             `json:"query"`
	SupportedTypes       []table.SearchType `json:"supported_types"`       // [1,3,4,2,6,7,8,9] on instagram, also 13 on messenger
	SessionID            any                `json:"session_id"`            // null
	SurfaceType          int64              `json:"surface_type"`          // 15
	SelectedParticipants []int64            `json:"selected_participants"` // null
	GroupID              *int64             `json:"group_id"`              // null
	CommunityID          *int64             `json:"community_id"`          // null
	QueryID              *int64             `json:"query_id"`              // null

	Secondary bool `json:"-"`
}

func (t *SearchUserTask) GetLabel() string {
	if t.Secondary {
		return TaskLabels["SearchUserSecondaryTask"]
	}
	return TaskLabels["SearchUserTask"]
}

func (t *SearchUserTask) Create() (any, string) {
	if t.Secondary {
		return t, "search_secondary"
	}
	return t, fmt.Sprintf(`["search_primary",%d]`, time.Now().UnixMilli())
}

type MessengerRestrictAction int64

const (
	MessengerRestrictActionRestrict MessengerRestrictAction = iota
	MessengerRestrictActionUnrestrict
)

type SetMessengerRestrictTask struct {
	RestricteeID int64                   `json:"restrictee_id"`
	RequestID    *string                 `json:"request_id"`
	Action       MessengerRestrictAction `json:"messenger_restrict_action"`
}

func (t *SetMessengerRestrictTask) GetLabel() string      { return TaskLabels["SetMessengerRestrictTask"] }
func (t *SetMessengerRestrictTask) Create() (any, string) { return t, "messenger_restrict" }

type MessengerBlockStatus int64

const (
	MessengerBlockStatusUnblocked MessengerBlockStatus = iota
	MessengerBlockStatusMessageBlocked
)

type SetMessengerBlockStatusTask struct {
	BlockeeID                int64                `json:"blockee_id"`
	BlockedByViewerStatus    MessengerBlockStatus `json:"blocked_by_viewer_status"`
	RequestID                *string              `json:"request_id"`
	UseOptimisticBlockStatus bool                 `json:"use_optimistic_block_status"`
}

func (t *SetMessengerBlockStatusTask) GetLabel() string {
	return TaskLabels["SetMessengerBlockStatusTask"]
}
func (t *SetMessengerBlockStatusTask) Create() (any, string) { return t, "block_status" }
