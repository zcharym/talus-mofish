//go:build !windows

package main

import "github.com/songwei.ma/talus-mofish/backend/watch"

func newCapturer() (watch.Capturer, error) {
	return watch.NewStubCapturer(), nil
}
