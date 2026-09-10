package socket_test

import (
	"encoding/json"
	"testing"

	"go.mewis.me/meta-extra/pkg/messagix/socket"
)

func TestFetchPollDetailsTask(t *testing.T) {
	task := &socket.FetchPollDetailsTask{PollID: 7}
	if task.GetLabel() != "545" {
		t.Fatalf("unexpected label: %s", task.GetLabel())
	}
	payload, queue := task.Create()
	if queue != "fetch_community_chat_poll_detail" {
		t.Fatalf("unexpected queue: %s", queue)
	}
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"poll_id":7}` {
		t.Fatalf("unexpected payload: %s", data)
	}
}
