package watch

import (
	"context"
	"log"
	"time"
)

type Watcher struct {
	cfg      *Config
	notifier *Notifier
	rules    *RuleEngine
	capture  Capturer
	gates    map[string]*HashGate
}

func NewWatcher(cfg *Config, notifier *Notifier, capture Capturer) *Watcher {
	return &Watcher{
		cfg:      cfg,
		notifier: notifier,
		rules:    NewRuleEngine(cfg.Cooldown),
		capture:  capture,
		gates:    make(map[string]*HashGate),
	}
}

func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.cfg.PollInterval)
	defer ticker.Stop()

	log.Printf("echo-watch started: %d target(s), poll=%s", len(w.cfg.Targets), w.cfg.PollInterval)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			w.pollOnce(ctx)
		}
	}
}

func (w *Watcher) pollOnce(ctx context.Context) {
	for _, target := range w.cfg.Targets {
		w.pollTarget(ctx, target)
	}
}

func (w *Watcher) pollTarget(ctx context.Context, target Target) {
	state, err := w.capture.Resolve(target)
	if err != nil {
		log.Printf("[%s] resolve: %v", target.Name, err)
		return
	}
	if !state.Capturable {
		log.Printf("[%s] skip: %s", target.Name, state.Reason)
		return
	}

	pixels, err := w.capture.CaptureRegion(state, target.Region)
	if err != nil {
		log.Printf("[%s] capture: %v", target.Name, err)
		return
	}

	gate := w.gateFor(target.Name)
	if !gate.Changed(pixels) {
		return
	}

	text, err := w.capture.OCR(pixels, target.Region.W, target.Region.H)
	if err != nil {
		log.Printf("[%s] ocr: %v", target.Name, err)
		return
	}
	if text == "" {
		return
	}

	match, ok := w.rules.Evaluate(target, text)
	if !ok {
		return
	}

	alert := match.Alert()
	if err := w.notifier.Send(ctx, alert); err != nil {
		log.Printf("[%s] notify: %v", target.Name, err)
		return
	}
	log.Printf("[%s] alert sent: rule=%s", target.Name, alert.RuleID)
}

func (w *Watcher) gateFor(name string) *HashGate {
	gate, ok := w.gates[name]
	if !ok {
		gate = &HashGate{}
		w.gates[name] = gate
	}
	return gate
}
