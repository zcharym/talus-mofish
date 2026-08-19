package aiclient

import "errors"

var (
	ErrModelRequired  = errors.New("ai model is required")
	ErrAPIKeyRequired = errors.New("ai api key is required")
	ErrNotConfigured  = errors.New("ai client is not configured")
)
