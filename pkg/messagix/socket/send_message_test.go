package socket_test

import (
	"encoding/json"
	"testing"

	"go.mewis.me/meta-extra/pkg/messagix/socket"
	"go.mewis.me/meta-extra/pkg/messagix/table"
)

func TestSendMessageTaskPayloadCompatibility(t *testing.T) {
	base := socket.SendMessageTask{ThreadId: 123, Otid: 456, SendType: table.TEXT, SyncGroup: 1, Text: "hello"}
	if base.GetLabel() != "46" {
		t.Fatalf("unexpected label: %s", base.GetLabel())
	}
	payload, queue := base.Create()
	if queue != "123" {
		t.Fatalf("unexpected queue: %s", queue)
	}
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	const baseline = `{"thread_id":123,"otid":"456","source":0,"send_type":1,"sync_group":1,"text":"hello","skip_url_preview_gen":0,"text_has_links":0,"multitab_env":0}`
	if string(data) != baseline {
		t.Fatalf("unexpected baseline payload: %s", data)
	}
}
