package messagix

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"go.mewis.me/meta-extra/pkg/messagix/socket"
	"go.mewis.me/meta-extra/pkg/messagix/table"
)

func (c *Client) ExecuteTasks(ctx context.Context, tasks ...socket.Task) (*table.LSTable, error) {
	if c == nil {
		return nil, ErrClientIsNil
	} else if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	tskm := c.newTaskManager()
	for _, task := range tasks {
		tskm.AddNewTask(task)
	}
	//tskm.setTraceId(methods.GenerateTraceID())

	payload, err := tskm.FinalizePayload()
	if err != nil {
		return nil, fmt.Errorf("failed to finalize payload: %w", err)
	}

	resp, err := c.makeLSRequest(ctx, payload, 3)
	if err != nil {
		return nil, err
	}

	tbl, err := resp.Parse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return tbl, nil
}

func (c *Client) ExecuteStatelessTask(ctx context.Context, task socket.Task) error {
	if c == nil {
		return ErrClientIsNil
	} else if ctx.Err() != nil {
		return ctx.Err()
	}
	innerPayload, queueName := task.Create()
	label := task.GetLabel()
	if queueName != "" {
		return fmt.Errorf("tried to execute stateful task %s as stateless", label)
	}
	innerPayloadMarshalled, err := json.Marshal(innerPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal inner task %s payload: %w", label, err)
	}
	outerPayload := socket.StatelessTaskData{
		Label:   label,
		Payload: string(innerPayloadMarshalled),
		Version: strconv.FormatInt(c.configs.VersionID, 10),
	}
	outerPayloadMarshalled, err := json.Marshal(outerPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal outer task %s payload: %w", label, err)
	}
	_, err = c.makeLSRequest(ctx, outerPayloadMarshalled, 4)
	return err
}

func (c *Client) SetThreadNickname(ctx context.Context, threadKey, contactID int64, nickname string) (*table.LSTable, error) {
	if c == nil {
		return nil, ErrClientIsNil
	} else if threadKey <= 0 || contactID <= 0 {
		return nil, fmt.Errorf("thread key and contact ID must be positive")
	}
	return c.ExecuteTasks(ctx, &socket.SetThreadNicknameTask{ThreadKey: strconv.FormatInt(threadKey, 10), ContactID: contactID, Nickname: nickname, SyncGroup: 1})
}

func (c *Client) SetThreadEmoji(ctx context.Context, threadKey int64, emoji string) (*table.LSTable, error) {
	if c == nil {
		return nil, ErrClientIsNil
	} else if threadKey <= 0 {
		return nil, fmt.Errorf("thread key must be positive")
	}
	return c.ExecuteTasks(ctx, &socket.SetThreadEmojiTask{ThreadKey: strconv.FormatInt(threadKey, 10), CustomEmoji: emoji, SyncGroup: 1})
}

func (c *Client) SetThreadApprovalMode(ctx context.Context, threadKey int64, enabled bool) (*table.LSTable, error) {
	if c == nil {
		return nil, ErrClientIsNil
	} else if threadKey <= 0 {
		return nil, fmt.Errorf("thread key must be positive")
	}
	value := 0
	if enabled {
		value = 1
	}
	return c.ExecuteTasks(ctx, &socket.SetThreadApprovalModeTask{ThreadKey: strconv.FormatInt(threadKey, 10), Enabled: value, SyncGroup: 1})
}

func (c *Client) SetThreadArchived(ctx context.Context, threadKey int64, archived bool) (*table.LSTable, error) {
	if c == nil {
		return nil, ErrClientIsNil
	} else if threadKey <= 0 {
		return nil, fmt.Errorf("thread key must be positive")
	}
	if archived {
		return c.ExecuteTasks(ctx, &socket.ArchiveThreadTask{ThreadKey: threadKey, SyncGroup: 1})
	}
	return c.ExecuteTasks(ctx, &socket.UnarchiveThreadTask{ThreadID: strconv.FormatInt(threadKey, 10), SyncGroup: 1})
}
