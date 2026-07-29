//go:build !windows

package main

import "github.com/songwei.ma/talus-mofish/internal/watch"

func newCapturer() (watch.Capturer, error) {
	return watch.NewStubCapturer(), nil
}
