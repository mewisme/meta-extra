package socket_test

import (
	"encoding/json"
	"testing"

	"go.mewis.me/meta-extra/pkg/messagix/socket"
)

func TestSetMessengerRestrictTask(t *testing.T) {
	tests := []struct {
		action socket.MessengerRestrictAction
		want   string
	}{
		{socket.MessengerRestrictActionRestrict, `{"restrictee_id":456,"request_id":null,"messenger_restrict_action":0}`},
		{socket.MessengerRestrictActionUnrestrict, `{"restrictee_id":456,"request_id":null,"messenger_restrict_action":1}`},
	}
	for _, test := range tests {
		task := &socket.SetMessengerRestrictTask{RestricteeID: 456, Action: test.action}
		if task.GetLabel() != "367" {
			t.Fatalf("unexpected label: %s", task.GetLabel())
		}
		payload, queue := task.Create()
		if queue != "messenger_restrict" {
			t.Fatalf("unexpected queue: %s", queue)
		}
		data, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != test.want {
			t.Fatalf("unexpected payload: %s", data)
		}
	}
}

func TestSetMessengerBlockStatusTask(t *testing.T) {
	tests := []struct {
		status socket.MessengerBlockStatus
		want   string
	}{
		{socket.MessengerBlockStatusUnblocked, `{"blockee_id":456,"blocked_by_viewer_status":0,"request_id":null,"use_optimistic_block_status":false}`},
		{socket.MessengerBlockStatusMessageBlocked, `{"blockee_id":456,"blocked_by_viewer_status":1,"request_id":null,"use_optimistic_block_status":false}`},
	}
	for _, test := range tests {
		task := &socket.SetMessengerBlockStatusTask{BlockeeID: 456, BlockedByViewerStatus: test.status}
		if task.GetLabel() != "151" {
			t.Fatalf("unexpected label: %s", task.GetLabel())
		}
		payload, queue := task.Create()
		if queue != "block_status" {
			t.Fatalf("unexpected queue: %s", queue)
		}
		data, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != test.want {
			t.Fatalf("unexpected payload: %s", data)
		}
	}
}
