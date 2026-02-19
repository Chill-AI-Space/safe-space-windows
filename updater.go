package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
}

// checkForUpdates queries GitHub for the latest release and returns the new
// version string if it is newer than appVersion, or "" otherwise.
func checkForUpdates() string {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/kobzevvv/vibe-sec-app-win/releases/latest")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ""
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return ""
	}

	tag := strings.TrimPrefix(release.TagName, "v")
	if isNewer(tag, appVersion) {
		return tag
	}
	return ""
}

// isNewer returns true if remote version is strictly greater than local.
func isNewer(remote, local string) bool {
	r := parseVersion(remote)
	l := parseVersion(local)
	for i := 0; i < 3; i++ {
		if r[i] > l[i] {
			return true
		}
		if r[i] < l[i] {
			return false
		}
	}
	return false
}

func parseVersion(v string) [3]int {
	parts := strings.SplitN(v, ".", 3)
	var result [3]int
	for i, p := range parts {
		if i >= 3 {
			break
		}
		result[i], _ = strconv.Atoi(p)
	}
	return result
}
