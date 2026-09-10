package socket_test

import (
	"encoding/json"
	"testing"

	"go.mewis.me/meta-extra/pkg/messagix/socket"
	"go.mewis.me/meta-extra/pkg/messagix/table"
)

func TestCoreTaskCompatibility(t *testing.T) {
	tests := []struct {
		name, label, queue, payload string
		task                        socket.Task
	}{
		{"reply", "46", "123", `{"thread_id":123,"otid":"456","source":0,"send_type":1,"sync_group":1,"reply_metadata":{"reply_source_id":"mid.reply","reply_source_type":1,"reply_type":0},"text":"hello","skip_url_preview_gen":0,"text_has_links":0,"multitab_env":0}`, &socket.SendMessageTask{ThreadId: 123, Otid: 456, SendType: table.TEXT, SyncGroup: 1, ReplyMetaData: &socket.ReplyMetaData{ReplyMessageId: "mid.reply", ReplySourceType: 1}, Text: "hello"}},
		{"edit", "742", "edit_message", `{"message_id":"mid.edit","text":"updated"}`, &socket.EditMessageTask{MessageID: "mid.edit", Text: "updated"}},
		{"unsend", "33", "unsend_message", `{"message_id":"mid.unsend"}`, &socket.DeleteMessageTask{MessageId: "mid.unsend"}},
		{"read", "21", "123", `{"thread_id":123,"last_read_watermark_ts":456,"sync_group":1}`, &socket.ThreadMarkReadTask{ThreadId: 123, LastReadWatermarkTs: 456, SyncGroup: 1}},
		{"rename", "32", "123", `{"thread_key":123,"thread_name":"group","sync_group":1}`, &socket.RenameThreadTask{ThreadKey: 123, ThreadName: "group", SyncGroup: 1}},
		{"typing", "3", "", `{"thread_key":123,"is_group_thread":1,"is_typing":1,"attribution":0,"sync_group":1,"thread_type":2}`, &socket.UpdatePresenceTask{ThreadKey: 123, IsGroupThread: 1, IsTyping: 1, SyncGroup: 1, ThreadType: 2}},
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

func TestCoreReactionCompatibility(t *testing.T) {
	task := &socket.SendReactionTask{ThreadKey: 123, MessageID: "mid.react", ActorID: 456, Reaction: "x"}
	if task.GetLabel() != "29" {
		t.Fatalf("unexpected label: %s", task.GetLabel())
	}
	payload, queue := task.Create()
	if queue != `["reaction","mid.react"]` {
		t.Fatalf("unexpected queue: %s", queue)
	}
	created := payload.(*socket.SendReactionTask)
	if created.TimestampMs <= 0 || created.SyncGroup != 1 || created.ThreadKey != 123 || created.MessageID != "mid.react" || created.ActorID != 456 || created.Reaction != "x" {
		t.Fatalf("unexpected reaction payload: %#v", created)
	}
}
