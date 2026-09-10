package messagix

import (
	"context"
	"errors"
	"testing"
)

func TestMessengerExtensionValidation(t *testing.T) {
	client := newStateTestClient()
	client.mailboxStateLoaded.Store(true)
	emptyCursor := ""
	tests := []struct {
		name string
		run  func() error
	}{
		{"nickname thread", func() error { _, err := client.SetThreadNickname(t.Context(), 0, 1, "mew"); return err }},
		{"nickname contact", func() error { _, err := client.SetThreadNickname(t.Context(), 1, 0, "mew"); return err }},
		{"emoji thread", func() error { _, err := client.SetThreadEmoji(t.Context(), 0, "x"); return err }},
		{"theme thread", func() error { _, err := client.SetThreadTheme(t.Context(), 0, 1); return err }},
		{"theme id", func() error { _, err := client.SetThreadTheme(t.Context(), 1, 0); return err }},
		{"call mute thread", func() error { _, err := client.SetThreadCallsMute(t.Context(), 0, 0); return err }},
		{"call mute expiry", func() error { _, err := client.SetThreadCallsMute(t.Context(), 1, -2); return err }},
		{"approval thread", func() error { _, err := client.SetThreadApprovalMode(t.Context(), 0, true); return err }},
		{"archive thread", func() error { _, err := client.SetThreadArchived(t.Context(), 0, true); return err }},
		{"pin thread", func() error { _, err := client.SetMessagePinned(t.Context(), 0, "m1", true); return err }},
		{"pin message", func() error { _, err := client.SetMessagePinned(t.Context(), 1, "", true); return err }},
		{"pinned query", func() error { _, err := client.ListPinnedMessages(0); return err }},
		{"poll query", func() error { _, err := client.GetPollDetails(0); return err }},
		{"poll fetch", func() error { _, err := client.FetchPollDetails(t.Context(), 0); return err }},
		{"search thread", func() error {
			_, err := client.SearchMessages(t.Context(), MessageSearchRequest{ThreadKey: 0, Query: "needle"})
			return err
		}},
		{"search query", func() error {
			_, err := client.SearchMessages(t.Context(), MessageSearchRequest{ThreadKey: 1, Query: "   "})
			return err
		}},
		{"search cursor", func() error {
			_, err := client.SearchMessages(t.Context(), MessageSearchRequest{ThreadKey: 1, Query: "needle", Cursor: &emptyCursor})
			return err
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.run(); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestMessengerExtensionCancellation(t *testing.T) {
	client := newStateTestClient()
	client.mailboxStateLoaded.Store(true)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	tests := []struct {
		name string
		run  func() error
	}{
		{"nickname", func() error { _, err := client.SetThreadNickname(ctx, 1, 2, "mew"); return err }},
		{"emoji", func() error { _, err := client.SetThreadEmoji(ctx, 1, "x"); return err }},
		{"theme", func() error { _, err := client.SetThreadTheme(ctx, 1, 2); return err }},
		{"call mute", func() error { _, err := client.SetThreadCallsMute(ctx, 1, 0); return err }},
		{"approval", func() error { _, err := client.SetThreadApprovalMode(ctx, 1, true); return err }},
		{"archive", func() error { _, err := client.SetThreadArchived(ctx, 1, true); return err }},
		{"pin", func() error { _, err := client.SetMessagePinned(ctx, 1, "m1", true); return err }},
		{"poll fetch", func() error { _, err := client.FetchPollDetails(ctx, 7); return err }},
		{"search", func() error {
			_, err := client.SearchMessages(ctx, MessageSearchRequest{ThreadKey: 1, Query: "needle"})
			return err
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.run(); !errors.Is(err, context.Canceled) {
				t.Fatalf("expected context cancellation, got %v", err)
			}
		})
	}
}

func TestMessengerExtensionNilClient(t *testing.T) {
	var client *Client
	if _, err := client.SetThreadEmoji(t.Context(), 1, "x"); !errors.Is(err, ErrClientIsNil) {
		t.Fatalf("expected nil-client error, got %v", err)
	}
	if _, err := client.FetchPollDetails(t.Context(), 1); !errors.Is(err, ErrClientIsNil) {
		t.Fatalf("expected nil-client error, got %v", err)
	}
	if _, err := client.ListPinnedMessages(1); !errors.Is(err, ErrClientIsNil) {
		t.Fatalf("expected nil-client error, got %v", err)
	}
	if _, err := client.SearchMessages(t.Context(), MessageSearchRequest{ThreadKey: 1, Query: "needle"}); !errors.Is(err, ErrClientIsNil) {
		t.Fatalf("expected nil-client error, got %v", err)
	}
}
