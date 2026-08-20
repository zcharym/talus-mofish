package obsidian

import (
	"strings"
	"testing"
)

func TestNormalizeVaultPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in      string
		want    string
		wantErr string
	}{
		{in: "", want: ""},
		{in: "  notes/daily.md  ", want: "notes/daily.md"},
		{in: "notes/", want: "notes"},
		{in: "notes/../secret.md", wantErr: ".."},
		{in: "/absolute.md", wantErr: "absolute"},
		{in: `C:\vault\note.md`, wantErr: "absolute"},
	}

	for _, tt := range tests {
		got, err := NormalizeVaultPath(tt.in)
		if tt.wantErr != "" {
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("NormalizeVaultPath(%q) err = %v, want %q", tt.in, err, tt.wantErr)
			}
			continue
		}
		if err != nil {
			t.Errorf("NormalizeVaultPath(%q) unexpected err = %v", tt.in, err)
			continue
		}
		if got != tt.want {
			t.Errorf("NormalizeVaultPath(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestEncodeVaultPath(t *testing.T) {
	t.Parallel()

	got := EncodeVaultPath("folder/hello world.md")
	want := "folder/hello%20world.md"
	if got != want {
		t.Fatalf("EncodeVaultPath = %q, want %q", got, want)
	}
}

func TestJoinVaultPath(t *testing.T) {
	t.Parallel()

	if got := JoinVaultPath("", "dir/"); got != "dir" {
		t.Fatalf("root dir = %q", got)
	}
	if got := JoinVaultPath("notes", "daily.md"); got != "notes/daily.md" {
		t.Fatalf("nested file = %q", got)
	}
}

func TestIsMarkdown(t *testing.T) {
	t.Parallel()

	if !IsMarkdown("Note.MD") {
		t.Fatal("expected .MD to be markdown")
	}
	if IsMarkdown("image.png") {
		t.Fatal("png should not be markdown")
	}
}
