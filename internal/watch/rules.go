package watch

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Match struct {
	Target Target
	Rule   Rule
	Text   string
}

type RuleEngine struct {
	cooldown time.Duration
	lastSent map[string]time.Time
}

func NewRuleEngine(cooldown time.Duration) *RuleEngine {
	return &RuleEngine{
		cooldown: cooldown,
		lastSent: make(map[string]time.Time),
	}
}

func (e *RuleEngine) Evaluate(target Target, text string) (*Match, bool) {
	normalized := strings.ToLower(text)
	for _, rule := range target.Rules {
		if e.matchesRule(rule, normalized, text) {
			key := target.Name + ":" + rule.ID
			if last, ok := e.lastSent[key]; ok && time.Since(last) < e.cooldown {
				continue
			}
			e.lastSent[key] = time.Now()
			return &Match{Target: target, Rule: rule, Text: text}, true
		}
	}
	return nil, false
}

func (e *RuleEngine) matchesRule(rule Rule, normalizedLower, original string) bool {
	for _, phrase := range rule.AnyText {
		if strings.Contains(normalizedLower, strings.ToLower(phrase)) {
			return true
		}
	}
	if rule.Regex != "" {
		re := regexp.MustCompile(rule.Regex)
		if re.MatchString(original) {
			return true
		}
	}
	return false
}

func (m Match) Alert() Alert {
	title := fmt.Sprintf("%s: %s", m.Target.Name, m.Rule.ID)
	body := fmt.Sprintf("Matched %q", m.Rule.ID)
	if len(m.Rule.AnyText) > 0 {
		body = fmt.Sprintf("Detected %q", m.Rule.AnyText[0])
	}
	severity := m.Rule.Severity
	if severity != "" {
		title = fmt.Sprintf("[%s] %s", strings.ToUpper(severity), title)
	}
	return Alert{
		RuleID: m.Rule.ID,
		Title:  title,
		Body:   body,
		Target: m.Target.Name,
	}
}
