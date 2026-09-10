package socket_test

import (
	"encoding/json"
	"testing"

	"go.mewis.me/meta-extra/pkg/messagix/socket"
)

func TestThreadCustomizationTasks(t *testing.T) {
	tests := []struct {
		name, label, queue, payload string
		task                        socket.Task
	}{
		{"nickname", "44", "thread_participant_nickname", `{"thread_key":"123","contact_id":456,"nickname":"mew","sync_group":1}`, &socket.SetThreadNicknameTask{ThreadKey: "123", ContactID: 456, Nickname: "mew", SyncGroup: 1}},
		{"emoji", "53", "thread_custom_emoji", `{"thread_key":"123","custom_emoji":"😀","sync_group":1}`, &socket.SetThreadEmojiTask{ThreadKey: "123", CustomEmoji: "😀", SyncGroup: 1}},
		{"approval", "28", "set_needs_admin_approval_for_new_participant", `{"thread_key":"123","enabled":1,"sync_group":1}`, &socket.SetThreadApprovalModeTask{ThreadKey: "123", Enabled: 1, SyncGroup: 1}},
		{"archive", "146", "123", `{"thread_key":123,"remove_type":1,"sync_group":1}`, &socket.ArchiveThreadTask{ThreadKey: 123, SyncGroup: 1}},
		{"unarchive", "242", "unarchive_thread", `{"thread_id":"123","sync_group":1}`, &socket.UnarchiveThreadTask{ThreadID: "123", SyncGroup: 1}},
	}
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
