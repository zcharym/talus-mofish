package watch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Alert struct {
	RuleID  string `json:"rule_id"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Target  string `json:"target,omitempty"`
}

type Notifier struct {
	client   *http.Client
	baseURL  string
	secret   string
}

func NewNotifier(workerURL, secret string) *Notifier {
	return &Notifier{
		client: &http.Client{Timeout: 15 * time.Second},
		baseURL: strings.TrimRight(workerURL, "/"),
		secret: secret,
	}
}

func (n *Notifier) Send(ctx context.Context, alert Alert) error {
	body, err := json.Marshal(alert)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.baseURL+"/api/alert", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+n.secret)

	res, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return fmt.Errorf("alert failed: HTTP %d", res.StatusCode)
	}
	return nil
}

func (n *Notifier) SendTest(ctx context.Context) error {
	return n.Send(ctx, Alert{
		RuleID: "test",
		Title:  "Echo Watch",
		Body:   "Test alert from echo-watch agent.",
		Target: "local",
	})
}
