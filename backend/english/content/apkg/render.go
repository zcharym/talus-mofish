package apkg

import (
	"regexp"
	"strings"
)

var (
	reCloze      = regexp.MustCompile(`\{\{c(\d+)::([^:}]+)(?:::([^}]*))?\}\}`)
	reTypeField  = regexp.MustCompile(`\{\{type:([^}]+)\}\}`)
	reFieldRef   = regexp.MustCompile(`\{\{([^#/][^}:]*)\}\}`)
	reConditional = regexp.MustCompile(`\{\{#([^}]+)\}\}([\s\S]*?)\{\{/([^}]+)\}\}`)
	reSound        = regexp.MustCompile(`\[sound:[^\]]+\]`)
)

func fieldMap(model *NoteModel, fields []string) map[string]string {
	m := make(map[string]string, len(model.Fields))
	for _, f := range model.Fields {
		idx := f.Ord
		if idx < 0 || idx >= len(fields) {
			m[f.Name] = ""
			continue
		}
		m[f.Name] = fields[idx]
	}
	return m
}

func expandCloze(text string, show bool) string {
	return reCloze.ReplaceAllStringFunc(text, func(match string) string {
		parts := reCloze.FindStringSubmatch(match)
		if len(parts) < 3 {
			return match
		}
		inner := parts[2]
		if show {
			return inner
		}
		return "[...]"
	})
}

func substituteFields(tmpl string, fields map[string]string, clozeShow bool) string {
	out := tmpl
	out = reTypeField.ReplaceAllStringFunc(out, func(match string) string {
		parts := reTypeField.FindStringSubmatch(match)
		if len(parts) < 2 {
			return match
		}
		return fields[parts[1]]
	})
	out = reConditional.ReplaceAllStringFunc(out, func(match string) string {
		parts := reConditional.FindStringSubmatch(match)
		if len(parts) < 4 || parts[1] != parts[3] {
			return match
		}
		name := parts[1]
		body := parts[2]
		if strings.TrimSpace(fields[name]) != "" {
			return body
		}
		return ""
	})
	if strings.Contains(out, "{{cloze:") {
		for name, val := range fields {
			out = strings.ReplaceAll(out, "{{cloze:"+name+"}}", expandCloze(val, clozeShow))
		}
	}
	for name, val := range fields {
		val = expandCloze(val, clozeShow)
		out = strings.ReplaceAll(out, "{{"+name+"}}", val)
	}
	out = reFieldRef.ReplaceAllString(out, "")
	out = reSound.ReplaceAllString(out, "")
	return out
}

const frontSidePlaceholder = "\x00FRONTSIDE\x00"

// RenderCard produces front and back HTML for a card template.
func RenderCard(model *NoteModel, fields []string, templateOrd int) (front, back string) {
	fm := fieldMap(model, fields)
	tmpl := templateByOrd(model, templateOrd)
	if tmpl == nil {
		return "", ""
	}
	front = substituteFields(tmpl.QFmt, fm, false)
	afmt := strings.ReplaceAll(tmpl.AFmt, "{{FrontSide}}", frontSidePlaceholder)
	back = substituteFields(afmt, fm, true)
	back = strings.ReplaceAll(back, frontSidePlaceholder, front)
	return front, back
}

func templateByOrd(model *NoteModel, ord int) *CardTemplate {
	for i := range model.Tmpls {
		if model.Tmpls[i].Ord == ord {
			return &model.Tmpls[i]
		}
	}
	if ord >= 0 && ord < len(model.Tmpls) {
		return &model.Tmpls[ord]
	}
	return nil
}

func FieldValue(fields []string, fieldIndex int) string {
	if fieldIndex < 0 || fieldIndex >= len(fields) {
		return ""
	}
	return fields[fieldIndex]
}

func FieldValueByName(fields []string, mapping map[string]int, key string) string {
	idx, ok := mapping[key]
	if !ok {
		return ""
	}
	return FieldValue(fields, idx)
}

func WordCountHTML(html string) int {
	text := stripHTML(html)
	if text == "" {
		return 0
	}
	return len(strings.Fields(text))
}

func stripHTML(s string) string {
	s = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(s, " ")
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	return strings.TrimSpace(s)
}
