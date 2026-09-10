package socket_test

import (
	"encoding/json"
	"testing"

	"go.mewis.me/meta-extra/pkg/messagix/socket"
)

func TestMessengerExtensionTasksAreStateful(t *testing.T) {
	tasks := []socket.Task{
		&socket.SetThreadNicknameTask{ThreadKey: "1", ContactID: 2, SyncGroup: 1},
		&socket.SetThreadEmojiTask{ThreadKey: "1", SyncGroup: 1},
		&socket.SetThreadApprovalModeTask{ThreadKey: "1", SyncGroup: 1},
		&socket.ArchiveThreadTask{ThreadKey: 1, SyncGroup: 1},
		&socket.UnarchiveThreadTask{ThreadID: "1", SyncGroup: 1},
		&socket.MuteThreadCallsTask{ThreadKey: "1", SyncGroup: 1},
		&socket.MessagePinTask{MessageID: "m1", ThreadKey: "1", Pinned: true},
		&socket.MessagePinTask{MessageID: "m1", ThreadKey: "1", Pinned: false},
		&socket.FetchPollDetailsTask{PollID: 1},
		&socket.ShareContactTask{ContactID: 2, ThreadID: 1, SyncGroup: 1},
		&socket.SearchMessagesTask{Query: "needle", ThreadKey: 1},
		&socket.SetMessengerRestrictTask{RestricteeID: 2},
		&socket.SetMessengerBlockStatusTask{BlockeeID: 2},
	}
	tasks = append(tasks, socket.NewSetThreadThemeTasks(1, 2)...)
	for _, task := range tasks {
		if task.GetLabel() == "" {
			t.Fatalf("%T has an empty label", task)
		}
		payload, queue := task.Create()
		if queue == "" {
			t.Fatalf("%T has an empty queue", task)
		}
		if _, err := json.Marshal(payload); err != nil {
			t.Fatalf("%T payload is not JSON serializable: %v", task, err)
		}
	}
}
