package vdiupload

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ErrorCode classifies job outcomes for exit codes and retries.
type ErrorCode string

const (
	CodeOK                 ErrorCode = "ok"
	CodeConfigInvalid      ErrorCode = "config_invalid"
	CodeWindowNotFound     ErrorCode = "window_not_found"
	CodeFocusFailed        ErrorCode = "focus_failed"
	CodeClipboardFailed    ErrorCode = "clipboard_failed"
	CodeInputFailed        ErrorCode = "input_failed"
	CodeCompletionTimeout  ErrorCode = "completion_timeout"
	CodeCompletionNegative ErrorCode = "completion_negative"
	CodeCancelled          ErrorCode = "cancelled"
	CodeUnsupportedOS      ErrorCode = "unsupported_os"
	CodeBusy               ErrorCode = "busy"
)

// ExitCode maps domain errors to process exit codes.
func (c ErrorCode) ExitCode() int {
	switch c {
	case CodeOK:
		return 0
	case CodeCancelled:
		return 130
	case CodeUnsupportedOS, CodeConfigInvalid:
		return 3
	default:
		return 2
	}
}

// Retryable reports whether the orchestrator should retry.
func (c ErrorCode) Retryable() bool {
	switch c {
	case CodeWindowNotFound, CodeFocusFailed, CodeClipboardFailed,
		CodeInputFailed, CodeCompletionTimeout, CodeCompletionNegative:
		return true
	default:
		return false
	}
}

// Job is one upload request.
type Job struct {
	Files      []string
	Profile    string
	DestHint   string
	Timeout    time.Duration
	DryRun     bool
	WaitWindow time.Duration
}

// Result is the outcome of Run.
type Result struct {
	Code     ErrorCode
	Attempts int
	Duration time.Duration
	Message  string
	Window   string
}

func (r Result) Err() error {
	if r.Code == CodeOK {
		return nil
	}
	if r.Message == "" {
		return fmt.Errorf("%s", r.Code)
	}
	return fmt.Errorf("%s: %s", r.Code, r.Message)
}

// ValidateFiles checks existence and profile clipboard limits.
func ValidateFiles(paths []string, clip ClipboardConfig) (abs []string, total int64, err error) {
	if len(paths) == 0 {
		return nil, 0, fmt.Errorf("%s: at least one file required", CodeConfigInvalid)
	}
	if clip.MaxFiles > 0 && len(paths) > clip.MaxFiles {
		return nil, 0, fmt.Errorf("%s: %d files exceeds max_files=%d", CodeConfigInvalid, len(paths), clip.MaxFiles)
	}
	abs = make([]string, 0, len(paths))
	for _, p := range paths {
		a, err := filepath.Abs(p)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: %w", CodeConfigInvalid, err)
		}
		st, err := os.Stat(a)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: %w", CodeConfigInvalid, err)
		}
		if st.IsDir() {
			return nil, 0, fmt.Errorf("%s: directories not supported yet: %s", CodeConfigInvalid, a)
		}
		total += st.Size()
		abs = append(abs, a)
	}
	if clip.MaxTotalBytes > 0 && total > clip.MaxTotalBytes {
		return nil, total, fmt.Errorf("%s: total size %d exceeds max_total_bytes=%d", CodeConfigInvalid, total, clip.MaxTotalBytes)
	}
	return abs, total, nil
}
