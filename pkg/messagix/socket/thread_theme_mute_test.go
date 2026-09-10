package socket_test

import (
	"encoding/json"
	"testing"

	"go.mewis.me/meta-extra/pkg/messagix/socket"
)

func TestThreadThemeAndCallMuteTasks(t *testing.T) {
	tests := []struct {
		name, label, queue, payload string
		task                        socket.Task
	}{{"call_mute", "229", "123", `{"thread_key":"123","mailbox_type":0,"mute_calls_expire_time_ms":456,"request_id":null,"sync_group":1}`, &socket.MuteThreadCallsTask{ThreadKey: "123", MuteCallsExpireTimeMS: 456, SyncGroup: 1}}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.task.GetLabel() != test.label {
				t.Fatalf("unexpected label: %s", test.task.GetLabel())
			}
			payload, queue := test.task.Create()
			if queue != test.queue {
				t.Fatalf("unexpected queue: %s", queue)
			}
			data, err := json.Marshal(payload)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != test.payload {
				t.Fatalf("unexpected payload: %s", data)
			}
		})
	}
}

func TestThreadThemeTasks(t *testing.T) {
	tasks := socket.NewSetThreadThemeTasks(123, 456)
	expected := []struct{ label, queue, payload string }{
		{"1013", "ai_generated_theme", `{"thread_key":123,"theme_fbid":456,"sync_group":1}`},
		{"1037", "msgr_custom_thread_theme", `{"thread_key":123,"theme_fbid":456,"sync_group":1}`},
		{"1028", "thread_theme_writer", `{"thread_key":123,"theme_fbid":456,"sync_group":1}`},
		{"43", "thread_theme", `{"thread_key":123,"theme_fbid":456,"sync_group":1,"source":null,"payload":null}`},
	}
	if len(tasks) != len(expected) {
		t.Fatalf("unexpected theme task count: %d", len(tasks))
	}
	for i, task := range tasks {
		if task.GetLabel() != expected[i].label {
			t.Fatalf("task %d label = %s", i, task.GetLabel())
		}
		payload, queue := task.Create()
		data, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		if queue != expected[i].queue || string(data) != expected[i].payload {
			t.Fatalf("task %d payload=%s queue=%s", i, data, queue)
		}
	}
}
