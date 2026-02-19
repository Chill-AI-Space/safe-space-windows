# vibe-sec-app-win

Windows system tray app for [vibe-sec](https://github.com/kobzevvv/vibe-sec) — security scanner for AI coding agents.

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

Download `vibe-sec.exe` from [Releases](https://github.com/kobzevvv/vibe-sec-app-win/releases/latest) and run it.

### Prerequisites

[vibe-sec](https://github.com/kobzevvv/vibe-sec) must be installed:

```bash
npx vibe-sec
```

## Build from source

```bash
go build -ldflags="-w -s -H windowsgui" -o vibe-sec.exe .
```

## Menu

| Item | Action |
|------|--------|
| `vibe-sec v1.0.0` | Version (disabled) |
| Status | Shows findings count or "Clean" |
| `$ npx vibe-sec` | Copies install command to clipboard |
| Injection Catcher | Shows hook guard status |
| Open Report | Opens HTML report in browser |
| Scan Now | Runs `scan-logs.mjs --static-only` |
| Update available | Copies update command to clipboard |
| Quit vibe-sec | Exit |

## Related Projects

| Project | Platform | Description |
|---------|----------|-------------|
| [vibe-sec](https://github.com/kobzevvv/vibe-sec) | All | Security scanner CLI (main project) |
| [vibe-sec-app](https://github.com/kobzevvv/vibe-sec-app) | macOS | Menubar app |
| [vibe-sec-dirty-machine](https://github.com/kobzevvv/vibe-sec-dirty-machine) | All | Test fixtures (fake secrets) |
