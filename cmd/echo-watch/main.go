package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/songwei.ma/talus-mofish/internal/watch"
)

func main() {
	configPath := flag.String("config", "", "path to echo-watch.yaml (default: %AppData%/TalusEcho/echo-watch.yaml)")
	triggerTest := flag.Bool("trigger-test", false, "send a test alert without screen capture")
	flag.Parse()

	path := *configPath
	if path == "" {
		var err error
		path, err = watch.DefaultConfigPath()
		if err != nil {
			log.Fatalf("resolve config path: %v", err)
		}
	}

	cfg, err := watch.LoadConfig(path)
	if err != nil {
		log.Fatalf("load config %s: %v", path, err)
	}

	notifier := watch.NewNotifier(cfg.WorkerURL, cfg.Secret)

	if *triggerTest {
		if err := notifier.SendTest(context.Background()); err != nil {
			log.Fatalf("trigger test failed: %v", err)
		}
		fmt.Println("test alert sent")
		return
	}

	capture, err := newCapturer()
	if err != nil {
		log.Fatalf("init capturer: %v", err)
	}

	w := watch.NewWatcher(cfg, notifier, capture)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := w.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("watch failed: %v", err)
	}
}
