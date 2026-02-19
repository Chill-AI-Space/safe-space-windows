package main

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/getlantern/systray"
)

const appVersion = "1.0.0"

type appState struct {
	mu        sync.Mutex
	result    *ScanResult
	scanning  bool
	latestVer string
}

var state = &appState{}

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTooltip("vibe-sec — security scanner for AI coding agents")

	mTitle := systray.AddMenuItem(fmt.Sprintf("vibe-sec v%s", appVersion), "")
	mTitle.Disable()

	mStatus := systray.AddMenuItem("Checking...", "")
	mStatus.Disable()

	mCmd := systray.AddMenuItem("$ npx vibe-sec", "Copy command to clipboard")

	systray.AddSeparator()

	mHookStatus := systray.AddMenuItem("Injection Catcher: checking...", "")
	mHookStatus.Disable()

	mHookCmd := systray.AddMenuItem("$ npx vibe-sec setup", "Copy setup command")
	mHookCmd.Hide()

	systray.AddSeparator()

	mOpenReport := systray.AddMenuItem("Open Report", "Open latest HTML report in browser")
	mScanNow := systray.AddMenuItem("Scan Now", "Run security scan")

	systray.AddSeparator()

	mUpdate := systray.AddMenuItem("", "")
	mUpdate.Hide()

	mQuit := systray.AddMenuItem("Quit vibe-sec", "")

	// refresh reads current state and updates all menu items.
	refresh := func() {
		state.mu.Lock()
		result := readScanResult()
		state.result = result
		ver := state.latestVer
		state.mu.Unlock()

		if !result.IsInstalled {
			mStatus.SetTitle("Not installed")
			mCmd.SetTitle("$ npx vibe-sec")
			mCmd.Enable()
			mCmd.Show()
		} else if result.Score < 0 {
			mStatus.SetTitle("No scan results yet")
			mCmd.SetTitle("$ npx vibe-sec scan")
			mCmd.Enable()
			mCmd.Show()
		} else if result.Score == 0 {
			mStatus.SetTitle("✓ Clean — no issues")
			if result.Date != "" {
				mCmd.SetTitle("Last scan: " + result.Date)
				mCmd.Disable()
				mCmd.Show()
			} else {
				mCmd.Hide()
			}
		} else {
			mStatus.SetTitle(fmt.Sprintf("● %d findings", result.Score))
			if result.Date != "" {
				mCmd.SetTitle("Last scan: " + result.Date)
				mCmd.Disable()
				mCmd.Show()
			} else {
				mCmd.Hide()
			}
		}

		if isHookInstalled() {
			mHookStatus.SetTitle("✓ Injection Catcher: active")
			mHookCmd.Hide()
		} else {
			mHookStatus.SetTitle("Injection Catcher: not installed")
			mHookCmd.Show()
		}

		if ver != "" {
			mUpdate.SetTitle(fmt.Sprintf("↑ Update available: v%s", ver))
			mUpdate.Show()
		}
	}

	// Initial state refresh.
	go refresh()

	// Periodic refresh every 60 seconds.
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			refresh()
		}
	}()

	// Check for updates on start and every 6 hours.
	go func() {
		if ver := checkForUpdates(); ver != "" {
			state.mu.Lock()
			state.latestVer = ver
			state.mu.Unlock()
			refresh()
		}
		ticker := time.NewTicker(6 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if ver := checkForUpdates(); ver != "" {
				state.mu.Lock()
				state.latestVer = ver
				state.mu.Unlock()
				refresh()
			}
		}
	}()

	// Handle menu clicks.
	go func() {
		for {
			select {
			case <-mCmd.ClickedCh:
				text := "npx vibe-sec"
				state.mu.Lock()
				if state.result != nil && state.result.IsInstalled && state.result.Score < 0 {
					text = "npx vibe-sec scan"
				}
				state.mu.Unlock()
				copyToClipboard(text)
				prev := mCmd.String()
				mCmd.SetTitle("✓ Copied — paste in Terminal")
				time.Sleep(2 * time.Second)
				_ = prev
				refresh()

			case <-mHookCmd.ClickedCh:
				copyToClipboard("npx vibe-sec setup")
				mHookCmd.SetTitle("✓ Copied — paste in Terminal")
				time.Sleep(2 * time.Second)
				refresh()

			case <-mOpenReport.ClickedCh:
				go openReport()

			case <-mScanNow.ClickedCh:
				go func() {
					state.mu.Lock()
					if state.scanning {
						state.mu.Unlock()
						return
					}
					state.scanning = true
					state.mu.Unlock()

					mScanNow.SetTitle("Scanning...")
					mScanNow.Disable()

					runScan()

					mScanNow.SetTitle("Scan Now")
					mScanNow.Enable()

					state.mu.Lock()
					state.scanning = false
					state.mu.Unlock()

					refresh()
				}()

			case <-mUpdate.ClickedCh:
				state.mu.Lock()
				ver := state.latestVer
				state.mu.Unlock()
				if ver != "" {
					updateCmd := `powershell -Command "& { $r = Invoke-RestMethod 'https://api.github.com/repos/kobzevvv/vibe-sec-app-win/releases/latest'; $url = $r.assets[0].browser_download_url; Invoke-WebRequest $url -OutFile vibe-sec.exe; Start-Process vibe-sec.exe }"`
					copyToClipboard(updateCmd)
					mUpdate.SetTitle("✓ Update command copied!")
					time.Sleep(2 * time.Second)
					mUpdate.SetTitle(fmt.Sprintf("↑ Update available: v%s", ver))
				}

			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {}

func copyToClipboard(text string) {
	cmd := exec.Command("clip")
	cmd.Stdin = strings.NewReader(text)
	_ = cmd.Run()
}

func openBrowser(url string) {
	exec.Command("cmd", "/c", "start", "", url).Start()
}
