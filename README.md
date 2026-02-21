# safe-space-windows

Windows system tray app for [safe-space](https://github.com/Chill-AI-Space/safe-space) — security scanner for AI coding agents.

![Windows](https://img.shields.io/badge/platform-Windows-blue)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8)

## Features

- **System tray icon** with scan status
- **Quick scan** — run security scan from tray menu
- **Open report** — view HTML report in browser
- **Injection Catcher status** — shows if real-time protection is active
- **Auto-update check** — notifies when new version is available

## Install

### Download

Download `safe-space.exe` from [Releases](https://github.com/Chill-AI-Space/safe-space-windows/releases/latest) and run it.

### Prerequisites

[safe-space](https://github.com/Chill-AI-Space/safe-space) must be installed:

```bash
npx safe-space
```

## Build from source

```bash
go build -ldflags="-w -s -H windowsgui" -o safe-space.exe .
```

## Menu

| Item | Action |
|------|--------|
| `safe-space v1.0.0` | Version (disabled) |
| Status | Shows findings count or "Clean" |
| `$ npx safe-space` | Copies install command to clipboard |
| Injection Catcher | Shows hook guard status |
| Open Report | Opens HTML report in browser |
| Scan Now | Runs `scan-logs.mjs --static-only` |
| Update available | Copies update command to clipboard |
| Quit safe-space | Exit |

## Related Projects

| Project | Platform | Description |
|---------|----------|-------------|
| [safe-space](https://github.com/Chill-AI-Space/safe-space) | All | Security scanner CLI (main project) |
| [safe-space-macos](https://github.com/Chill-AI-Space/safe-space-macos) | macOS | Menubar app |
| [safe-space-sandbox](https://github.com/Chill-AI-Space/safe-space-sandbox) | All | Test fixtures (fake secrets) |
