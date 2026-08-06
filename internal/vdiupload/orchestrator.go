package vdiupload

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// Dependencies are ports required by the orchestrator.
type Dependencies struct {
	Windows   WindowLocator
	Clipboard Clipboard
	Input     InputAutomator
	Detect    CompletionDetector
	Log       *slog.Logger
}

// Orchestrator runs upload jobs against a profile.
type Orchestrator struct {
	cfg  *Config
	deps Dependencies
}

// NewOrchestrator wires config and ports.
func NewOrchestrator(cfg *Config, deps Dependencies) *Orchestrator {
	if deps.Log == nil {
		deps.Log = slog.Default()
	}
	return &Orchestrator{cfg: cfg, deps: deps}
}

// Run executes a job with retries from the selected profile.
func (o *Orchestrator) Run(ctx context.Context, job Job) Result {
	start := time.Now()
	profile, err := o.cfg.ProfileOrDefault(job.Profile)
	if err != nil {
		return Result{Code: CodeConfigInvalid, Message: err.Error(), Duration: time.Since(start)}
	}

	files, _, err := ValidateFiles(job.Files, profile.Clipboard)
	if err != nil {
		return Result{Code: CodeConfigInvalid, Message: err.Error(), Duration: time.Since(start)}
	}

	attempts := profile.Retry.Attempts
	if attempts < 1 {
		attempts = 1
	}

	var last Result
	for attempt := 1; attempt <= attempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return Result{Code: CodeCancelled, Attempts: attempt - 1, Duration: time.Since(start), Message: err.Error()}
		}

		o.deps.Log.Info("vdiupload attempt", "attempt", attempt, "files", len(files), "dry_run", job.DryRun)
		last = o.runOnce(ctx, job, profile, files, attempt)
		last.Duration = time.Since(start)
		if last.Code == CodeOK {
			return last
		}
		if !last.Code.Retryable() || attempt == attempts {
			return last
		}

		delay := backoff(profile.Retry, attempt)
		o.deps.Log.Warn("vdiupload retrying", "code", last.Code, "after", delay.String())
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return Result{Code: CodeCancelled, Attempts: attempt, Duration: time.Since(start), Message: ctx.Err().Error()}
		case <-timer.C:
		}
	}
	return last
}

func (o *Orchestrator) runOnce(ctx context.Context, job Job, profile Profile, files []string, attempt int) Result {
	if o.deps.Windows == nil || o.deps.Clipboard == nil || o.deps.Input == nil {
		return Result{Code: CodeUnsupportedOS, Attempts: attempt, Message: "Win32 backends not wired"}
	}

	win, err := o.findWindow(ctx, profile.Window, job.WaitWindow)
	if err != nil {
		return Result{Code: CodeWindowNotFound, Attempts: attempt, Message: err.Error()}
	}

	if job.DryRun {
		o.deps.Log.Info("dry-run: would upload", "window", win.Title, "files", files, "steps", len(profile.Steps))
		return Result{Code: CodeOK, Attempts: attempt, Window: win.Title, Message: "dry-run"}
	}

	restore, err := o.deps.Clipboard.SetFiles(ctx, files, profile.Clipboard)
	if err != nil {
		return Result{Code: CodeClipboardFailed, Attempts: attempt, Window: win.Title, Message: err.Error()}
	}
	if restore != nil {
		defer restore()
	}

	for i, step := range profile.Steps {
		if err := ctx.Err(); err != nil {
			return Result{Code: CodeCancelled, Attempts: attempt, Window: win.Title, Message: err.Error()}
		}
		if err := o.runStep(ctx, win, step); err != nil {
			code := classifyStepErr(err)
			return Result{Code: code, Attempts: attempt, Window: win.Title, Message: fmt.Sprintf("step %d: %v", i, err)}
		}
	}

	return Result{Code: CodeOK, Attempts: attempt, Window: win.Title}
}

func (o *Orchestrator) findWindow(ctx context.Context, match WindowMatch, wait time.Duration) (Window, error) {
	if wait <= 0 {
		return o.deps.Windows.Find(ctx, match)
	}

	findCtx, cancel := context.WithTimeout(ctx, wait)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	var lastErr error
	for {
		win, err := o.deps.Windows.Find(findCtx, match)
		if err == nil {
			return win, nil
		}
		lastErr = err
		select {
		case <-findCtx.Done():
			if lastErr != nil {
				return Window{}, fmt.Errorf("%w: %v", lastErr, findCtx.Err())
			}
			return Window{}, findCtx.Err()
		case <-ticker.C:
		}
	}
}

func (o *Orchestrator) runStep(ctx context.Context, win Window, step Step) error {
	if step.Focus {
		if err := o.deps.Input.Focus(ctx, win); err != nil {
			return fmt.Errorf("%s: %w", CodeFocusFailed, err)
		}
	}
	if len(step.Keys) > 0 {
		if err := o.deps.Input.RunKeys(ctx, step.Keys); err != nil {
			return fmt.Errorf("%s: %w", CodeInputFailed, err)
		}
	}
	if step.WaitMs > 0 {
		timer := time.NewTimer(time.Duration(step.WaitMs) * time.Millisecond)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
	if step.WaitOCR != nil {
		if o.deps.Detect == nil {
			return fmt.Errorf("%s: completion detector not configured", CodeCompletionTimeout)
		}
		waitCtx := ctx
		var cancel context.CancelFunc
		if step.WaitOCR.Timeout > 0 {
			waitCtx, cancel = context.WithTimeout(ctx, step.WaitOCR.Timeout)
			defer cancel()
		}
		if err := o.deps.Detect.WaitOCR(waitCtx, win, *step.WaitOCR); err != nil {
			return err
		}
	}
	return nil
}

func classifyStepErr(err error) ErrorCode {
	msg := err.Error()
	for _, code := range []ErrorCode{
		CodeFocusFailed, CodeInputFailed, CodeCompletionTimeout, CodeCompletionNegative, CodeCancelled,
	} {
		if containsCode(msg, string(code)) {
			return code
		}
	}
	return CodeInputFailed
}

func containsCode(msg, code string) bool {
	return len(msg) >= len(code) && (msg == code || len(msg) > len(code) && (msg[:len(code)+1] == code+":" || msg[:len(code)+1] == code+" "))
}

func backoff(r RetryConfig, attempt int) time.Duration {
	d := r.BaseDelay
	for i := 1; i < attempt; i++ {
		d *= 2
		if d > r.MaxDelay {
			return r.MaxDelay
		}
	}
	if d > r.MaxDelay {
		return r.MaxDelay
	}
	return d
}
