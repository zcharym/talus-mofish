package agent

import (
	"strings"
	"testing"

	"github.com/songwei.ma/talus-mofish/backend/storage/store"
	"github.com/songwei.ma/talus-mofish/backend/utils/aiclient"
)

func TestBuildMessagesIncludesDomainContext(t *testing.T) {
	msgs := BuildMessages(nil, "Which worker is noisy?", "Cloudflare snapshot: errors=12")
	if len(msgs) != 2 {
		t.Fatalf("len(msgs) = %d, want 2", len(msgs))
	}
	if msgs[0].Role != aiclient.RoleSystem {
		t.Fatalf("first role = %q, want system", msgs[0].Role)
	}
	if !strings.Contains(msgs[0].Content, "Cloudflare snapshot: errors=12") {
		t.Fatalf("system prompt missing domain context: %q", msgs[0].Content)
	}
	if msgs[1].Role != aiclient.RoleUser || msgs[1].Content != "Which worker is noisy?" {
		t.Fatalf("user message = %+v", msgs[1])
	}
}

func TestBuildMessagesOmitsEmptyDomainContext(t *testing.T) {
	msgs := BuildMessages(nil, "hello", "  ")
	if !strings.Contains(msgs[0].Content, "You are Talus Agent") {
		t.Fatalf("expected default system prompt, got %q", msgs[0].Content)
	}
	if strings.Contains(msgs[0].Content, "\n\n") {
		t.Fatalf("empty domain context should not add a blank system addendum: %q", msgs[0].Content)
	}
}

func TestBuildMessagesSkipsEmptyHistoryContent(t *testing.T) {
	history := []store.ChatMessage{
		{Role: "user", Content: "   "},
		{Role: "assistant", Content: "Earlier answer"},
	}
	msgs := BuildMessages(history, "next", "")
	if len(msgs) != 3 {
		t.Fatalf("len(msgs) = %d, want 3 (system + assistant + user)", len(msgs))
	}
	if msgs[1].Role != aiclient.RoleAssistant || msgs[1].Content != "Earlier answer" {
		t.Fatalf("history message = %+v", msgs[1])
	}
}
