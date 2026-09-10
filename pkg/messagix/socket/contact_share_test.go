package socket_test

import (
	"encoding/json"
	"testing"

	"go.mewis.me/meta-extra/pkg/messagix/socket"
)

func TestShareContactTask(t *testing.T) {
	task := &socket.ShareContactTask{ContactID: 456, SyncGroup: 1, ThreadID: 123}
	if task.GetLabel() != "359" {
		t.Fatalf("unexpected label: %s", task.GetLabel())
	}
	payload, queue := task.Create()
	if queue != "messenger_contact_sharing" {
		t.Fatalf("unexpected queue: %s", queue)
	}
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"contact_id":456,"sync_group":1,"text":null,"thread_id":"123"}` {
		t.Fatalf("unexpected payload: %s", data)
	}

	text := "hello"
	task.Text = &text
	data, err = json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"contact_id":456,"sync_group":1,"text":"hello","thread_id":"123"}` {
		t.Fatalf("unexpected payload with text: %s", data)
	}
}
