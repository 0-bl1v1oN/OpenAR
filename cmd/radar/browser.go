package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/nospy/albion-openradar/internal/logger"
)

const radarURL = "http://localhost:5001"

func (app *App) openRadarInBrowser(appWindow bool) {
	if err := app.waitForHTTPReady(5 * time.Second); err != nil {
		logger.PrintWarn("APP", "Browser launch skipped: %v", err)
		return
	}
	if err := openBrowser(radarURL, appWindow); err != nil {
		logger.PrintWarn("APP", "Could not open browser automatically: %v", err)
		return
	}
	logger.PrintSuccess("APP", "Opened radar: %s", radarURL)
}

func (app *App) waitForHTTPReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if app.ctx.Err() != nil {
			return app.ctx.Err()
		}
		if atomic.LoadInt32(&app.httpRunning) == 1 {
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return fmt.Errorf("HTTP server did not become ready within %s", timeout)
}

func openBrowser(url string, appWindow bool) error {
	var candidates [][]string
	switch runtime.GOOS {
	case "windows":
		if appWindow {
			candidates = append(candidates,
				[]string{"msedge", "--app=" + url},
				[]string{"chrome", "--app=" + url},
				[]string{"cmd", "/c", "start", "", url},
			)
		} else {
			candidates = append(candidates, []string{"cmd", "/c", "start", "", url})
		}
	case "darwin":
		if appWindow {
			candidates = append(candidates,
				[]string{"open", "-na", "Microsoft Edge", "--args", "--app=" + url},
				[]string{"open", "-na", "Google Chrome", "--args", "--app=" + url},
			)
		}
		candidates = append(candidates, []string{"open", url})
	default:
		if appWindow {
			candidates = append(candidates,
				[]string{"microsoft-edge", "--app=" + url},
				[]string{"google-chrome", "--app=" + url},
				[]string{"chromium", "--app=" + url},
			)
		}
		candidates = append(candidates, []string{"xdg-open", url})
	}

	var lastErr error
	for _, candidate := range candidates {
		cmd := exec.Command(candidate[0], candidate[1:]...)
		if err := cmd.Start(); err != nil {
			lastErr = err
			continue
		}
		return nil
	}
	return lastErr
}

func (app *App) waitForShutdownSignal() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)
	select {
	case <-app.ctx.Done():
	case <-sigCh:
		app.cancel()
	}
}
