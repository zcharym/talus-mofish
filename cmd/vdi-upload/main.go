package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/songwei.ma/talus-mofish/backend/vdiupload"
)

func main() {
	configPath := flag.String("config", "", "path to vdi-upload.yaml")
	profile := flag.String("profile", "", "profile name (default from config)")
	dryRun := flag.Bool("dry-run", false, "resolve config/steps without clipboard or input")
	waitWindow := flag.Duration("wait-window", 0, "wait for VDI client window before failing")
	logFormat := flag.String("log-format", "text", "text|json")
	var files multiFlag
	flag.Var(&files, "file", "file to upload (repeatable)")
	flag.Parse()

	logger := newLogger(*logFormat)
	slog.SetDefault(logger)

	path := *configPath
	if path == "" {
		var err error
		path, err = vdiupload.DefaultConfigPath()
		if err != nil {
			fatal(vdiupload.CodeConfigInvalid, "resolve config path: %v", err)
		}
	}

	cfg, err := vdiupload.LoadConfig(path)
	if err != nil {
		fatal(vdiupload.CodeConfigInvalid, "load config %s: %v", path, err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	orch := vdiupload.NewOrchestrator(cfg, vdiupload.DefaultDependencies())
	result := orch.Run(ctx, vdiupload.Job{
		Files:      files,
		Profile:    *profile,
		DryRun:     *dryRun,
		WaitWindow: *waitWindow,
	})

	if result.Code == vdiupload.CodeOK {
		logger.Info("upload finished", "code", result.Code, "attempts", result.Attempts, "window", result.Window, "duration", result.Duration.String())
		os.Exit(0)
	}

	logger.Error("upload failed", "code", result.Code, "attempts", result.Attempts, "err", result.Message, "duration", result.Duration.String())
	os.Exit(result.Code.ExitCode())
}

func fatal(code vdiupload.ErrorCode, format string, args ...any) {
	slog.Error(fmt.Sprintf(format, args...), "code", code)
	os.Exit(code.ExitCode())
}

func newLogger(format string) *slog.Logger {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	switch strings.ToLower(format) {
	case "json":
		return slog.New(slog.NewJSONHandler(os.Stderr, opts))
	default:
		return slog.New(slog.NewTextHandler(os.Stderr, opts))
	}
}

type multiFlag []string

func (m *multiFlag) String() string { return strings.Join(*m, ",") }

func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}
