package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// ScanResult holds parsed state from the latest vibe-sec report.
type ScanResult struct {
	Score       int    // -1 = no results, 0 = clean, >0 = findings count
	Date        string // "2026-02-19" extracted from report filename
	IsInstalled bool   // true if vibe-sec scripts directory exists
}

// configDir returns the vibe-sec configuration directory.
// Checks %USERPROFILE%\.config\vibe-sec first (Node.js default),
// then %APPDATA%\vibe-sec as fallback.
func configDir() string {
	home := os.Getenv("USERPROFILE")
	if home == "" {
		home = os.Getenv("HOME")
	}

	primary := filepath.Join(home, ".config", "vibe-sec")
	if info, err := os.Stat(primary); err == nil && info.IsDir() {
		return primary
	}

	if appdata := os.Getenv("APPDATA"); appdata != "" {
		alt := filepath.Join(appdata, "vibe-sec")
		if info, err := os.Stat(alt); err == nil && info.IsDir() {
			return alt
		}
	}

	return primary
}

// claudeSettingsPath returns path to Claude Code settings.json.
func claudeSettingsPath() string {
	home := os.Getenv("USERPROFILE")
	if home == "" {
		home = os.Getenv("HOME")
	}
	return filepath.Join(home, ".claude", "settings.json")
}

// readScanResult reads the latest vibe-sec scan report and returns parsed state.
func readScanResult() *ScanResult {
	dir := configDir()

	// Check if vibe-sec is installed (scripts directory exists).
	scriptsDir := filepath.Join(dir, "scripts")
	isInstalled := false
	if _, err := os.Stat(filepath.Join(scriptsDir, "scan-logs.mjs")); err == nil {
		isInstalled = true
	}

	// Find latest markdown report.
	matches, err := filepath.Glob(filepath.Join(dir, "vibe-sec-log-report-*.md"))
	if err != nil || len(matches) == 0 {
		return &ScanResult{Score: -1, IsInstalled: isInstalled}
	}

	sort.Strings(matches)
	latest := matches[len(matches)-1]

	data, err := os.ReadFile(latest)
	if err != nil {
		return &ScanResult{Score: -1, IsInstalled: isInstalled}
	}

	// Extract date from filename: vibe-sec-log-report-2026-02-19.md
	date := ""
	if m := dateRe.FindString(filepath.Base(latest)); m != "" {
		date = m
	}

	return &ScanResult{
		Score:       parseScore(string(data)),
		Date:        date,
		IsInstalled: isInstalled,
	}
}

var (
	dateRe             = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	findingsCommentRe  = regexp.MustCompile(`<!-- findings: (\d+) -->`)
	noIssuesRe         = regexp.MustCompile(`(?i)no static issues found`)
	critHighRe         = regexp.MustCompile(`(\d+)\s+critical\s+and\s+(\d+)\s+high`)
	issuesRe           = regexp.MustCompile(`(\d+)\s+(?:issues?|findings?)`)
)

// parseScore extracts the number of findings from report content.
func parseScore(content string) int {
	// Primary: <!-- findings: N -->
	if m := findingsCommentRe.FindStringSubmatch(content); len(m) > 1 {
		if n, err := strconv.Atoi(m[1]); err == nil {
			return n
		}
	}
	// "no static issues found" → 0
	if noIssuesRe.MatchString(content) {
		return 0
	}
	// "N critical and M high" → sum
	if m := critHighRe.FindStringSubmatch(content); len(m) > 2 {
		c, _ := strconv.Atoi(m[1])
		h, _ := strconv.Atoi(m[2])
		return c + h
	}
	// "N issues" / "N findings"
	if m := issuesRe.FindStringSubmatch(content); len(m) > 1 {
		if n, err := strconv.Atoi(m[1]); err == nil {
			return n
		}
	}
	return -1
}

// isHookInstalled checks if vibe-sec hook is registered in Claude settings.
func isHookInstalled() bool {
	data, err := os.ReadFile(claudeSettingsPath())
	if err != nil {
		return false
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return false
	}

	hooks, ok := settings["hooks"].(map[string]interface{})
	if !ok {
		return false
	}
	preToolUse, ok := hooks["PreToolUse"].([]interface{})
	if !ok {
		return false
	}

	for _, entry := range preToolUse {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		hooksList, ok := entryMap["hooks"].([]interface{})
		if !ok {
			continue
		}
		for _, h := range hooksList {
			hookMap, ok := h.(map[string]interface{})
			if !ok {
				continue
			}
			cmd, _ := hookMap["command"].(string)
			if strings.Contains(cmd, "hook.mjs") {
				return true
			}
		}
	}
	return false
}

// runScan executes the vibe-sec scanner.
func runScan() {
	dir := configDir()
	scriptPath := filepath.Join(dir, "scripts", "scan-logs.mjs")
	cmd := exec.Command("node", scriptPath, "--static-only", "--source", "app")
	cmd.Dir = dir
	_ = cmd.Run()
}

// openReport opens the latest HTML report in the default browser.
func openReport() {
	dir := configDir()

	// Look for HTML report first.
	matches, _ := filepath.Glob(filepath.Join(dir, "vibe-sec-log-report-*.html"))
	if len(matches) > 0 {
		sort.Strings(matches)
		openBrowser(matches[len(matches)-1])
		return
	}

	// Fall back to starting the report server.
	scriptPath := filepath.Join(dir, "scripts", "serve-report.mjs")
	if _, err := os.Stat(scriptPath); err == nil {
		cmd := exec.Command("node", scriptPath)
		cmd.Dir = dir
		_ = cmd.Start()
		openBrowser("http://localhost:7777")
	}
}
