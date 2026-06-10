package apkg

import "testing"

func TestRenderCardBasic(t *testing.T) {
	model := &NoteModel{
		Type: 0,
		Fields: []ModelField{
			{Name: "Front", Ord: 0},
			{Name: "Back", Ord: 1},
		},
		Tmpls: []CardTemplate{{
			Ord:  0,
			QFmt: "{{Front}}",
			AFmt: "{{FrontSide}}<hr>{{Back}}",
		}},
	}
	front, back := RenderCard(model, []string{"hello", "world"}, 0)
	if front != "hello" {
		t.Fatalf("front = %q", front)
	}
	if back != "hello<hr>world" {
		t.Fatalf("back = %q", back)
	}
}

func TestSplitFields(t *testing.T) {
	got := splitFields("a\x1fb\x1fc")
	if len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Fatalf("splitFields = %v", got)
	}
}
