package socket_test

import (
	"encoding/json"
	"testing"

	"go.mewis.me/meta-extra/pkg/messagix/socket"
)

func TestMessagePinTask(t *testing.T) {
	for _, test := range []struct {
		name, label, queue string
		pinned             bool
	}{
		{"pin", "430", "pin_msg_v2_123", true},
		{"unpin", "431", "unpin_msg_v2_123", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			task := &socket.MessagePinTask{MessageID: "mid.$test", ThreadKey: "123", TimestampMS: 456, Pinned: test.pinned}
			if task.GetLabel() != test.label {
				t.Fatalf("unexpected label: %s", task.GetLabel())
			}
			payload, queue := task.Create()
			if queue != test.queue {
				t.Fatalf("unexpected queue: %s", queue)
			}
			data, err := json.Marshal(payload)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != `{"message_id":"mid.$test","thread_key":"123","timestamp_ms":456}` {
				t.Fatalf("unexpected payload: %s", data)
			}
		})
	}
}
