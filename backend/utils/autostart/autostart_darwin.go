//go:build darwin

package autostart

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (m *Manager) enable(identifier, exe string, args []string) error {
	plistPath, err := launchAgentPath(identifier)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(plistPath), 0o755); err != nil {
		return fmt.Errorf("create launch agents directory: %w", err)
	}

	programArgs := append([]string{exe}, args...)
	content := buildLaunchAgentPlist(identifier, programArgs)
	if err := os.WriteFile(plistPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write launch agent: %w", err)
	}
	return nil
}

func (m *Manager) disable(identifier string) error {
	plistPath, err := launchAgentPath(identifier)
	if err != nil {
		return err
	}
	if err := os.Remove(plistPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove launch agent: %w", err)
	}
	return nil
}

func (m *Manager) status(identifier string) (Status, error) {
	plistPath, err := launchAgentPath(identifier)
	if err != nil {
		return Status{}, err
	}
	if _, err := os.Stat(plistPath); err != nil {
		if os.IsNotExist(err) {
			return Status{Enabled: false, Strategy: "launch-agent"}, nil
		}
		return Status{}, fmt.Errorf("stat launch agent: %w", err)
	}
	return Status{
		Enabled:  true,
		Path:     plistPath,
		Strategy: "launch-agent",
	}, nil
}

func launchAgentPath(identifier string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}
	return filepath.Join(home, "Library", "LaunchAgents", identifier+".plist"), nil
}

func buildLaunchAgentPlist(identifier string, programArgs []string) string {
	argsXML := ""
	for _, arg := range programArgs {
		argsXML += fmt.Sprintf("        <string>%s</string>\n", xmlEscape(arg))
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>%s</string>
    <key>ProgramArguments</key>
    <array>
%s    </array>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>
`, xmlEscape(identifier), argsXML)
}

func xmlEscape(value string) string {
	replacer := strings.NewReplacer(
		`&`, "&amp;",
		`<`, "&lt;",
		`>`, "&gt;",
		`"`, "&quot;",
		`'`, "&apos;",
	)
	return replacer.Replace(value)
}
