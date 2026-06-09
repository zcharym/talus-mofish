//go:build windows

package autostart

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const runKey = `Software\Microsoft\Windows\CurrentVersion\Run`

func (m *Manager) enable(identifier, exe string, args []string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("open run key: %w", err)
	}
	defer key.Close()

	value := quoteWindowsArg(exe)
	if len(args) > 0 {
		quoted := make([]string, len(args))
		for i, arg := range args {
			quoted[i] = quoteWindowsArg(arg)
		}
		value = value + " " + strings.Join(quoted, " ")
	}

	if err := key.SetStringValue(identifier, value); err != nil {
		return fmt.Errorf("set run value: %w", err)
	}
	return nil
}

func (m *Manager) disable(identifier string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("open run key: %w", err)
	}
	defer key.Close()

	if err := key.DeleteValue(identifier); err != nil {
		if err == registry.ErrNotExist {
			return nil
		}
		return fmt.Errorf("delete run value: %w", err)
	}
	return nil
}

func (m *Manager) status(identifier string) (Status, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.QUERY_VALUE)
	if err != nil {
		return Status{}, fmt.Errorf("open run key: %w", err)
	}
	defer key.Close()

	value, _, err := key.GetStringValue(identifier)
	if err != nil {
		if err == registry.ErrNotExist {
			return Status{Enabled: false, Strategy: "registry"}, nil
		}
		return Status{}, fmt.Errorf("query run value: %w", err)
	}

	return Status{
		Enabled:  true,
		Path:     value,
		Strategy: "registry",
	}, nil
}

func quoteWindowsArg(arg string) string {
	if arg == "" {
		return `""`
	}
	if !strings.ContainsAny(arg, " \t\"") {
		return arg
	}
	return `"` + strings.ReplaceAll(arg, `"`, `\"`) + `"`
}
