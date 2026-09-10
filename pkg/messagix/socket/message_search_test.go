package socket_test

import (
	"encoding/json"
	"testing"

	"go.mewis.me/meta-extra/pkg/messagix/socket"
)

func TestSearchMessagesTask(t *testing.T) {
	tests := []struct {
		name   string
		cursor *string
		want   string
	}{
		{name: "initial", want: `{"query":"needle","type":2,"thread_key":456,"next_page_cursor":null,"client_caller_id":"msgr_search_web_start_conversation_message_search"}`},
		{name: "pagination", cursor: ptr("cursor-1"), want: `{"query":"needle","type":2,"thread_key":456,"next_page_cursor":"cursor-1","client_caller_id":"msgr_search_web_in_conversation_message_search"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			task := &socket.SearchMessagesTask{Query: "needle", ThreadKey: 456, NextPageCursor: test.cursor}
			if task.GetLabel() != "107" {
				t.Fatalf("unexpected label: %s", task.GetLabel())
			}
			payload, queue := task.Create()
			if queue != "message_search" {
				t.Fatalf("unexpected queue: %s", queue)
			}
			data, err := json.Marshal(payload)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != test.want {
				t.Fatalf("unexpected payload: %s", data)
			}
		})
	}
}

func ptr(value string) *string { return &value }
