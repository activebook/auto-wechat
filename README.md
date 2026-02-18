# auto-wechat

A desktop automation tool that bulk-sends messages to every contact in WeChat's contact list using clipboard content.
It uses computer-vision (OpenCV via [gocv](https://gocv.io)) to locate UI elements and to drive the mouse and keyboard — no WeChat API or plugin required.

## How it works

The tool operates in two phases:

```
scan  →  calibrate UI element positions  →  cfgs/settings.yml
proc  →  iterate contacts & send message
```

1. **`scan`** — Auto scan each WeChat UI element (Contacts button, More button, Messages bar), and record the screen coordinates into `cfgs/settings.yml`.
2. **`proc`** — Reads `settings.yml`, activates WeChat, then loops through the contact list:
   - Navigates to the last visited contact.
   - Presses ↓ to advance to the next contact.
   - Right-clicks → "Send Message" → pastes the clipboard content → (optionally) hits Enter.
   - Repeats until `max_count` messages are sent or the end of the list is detected via image matching.

> **Tip:** Copy your message to the clipboard before running `proc`.

## Prerequisites

| Dependency | Notes |
|---|---|
| [Go](https://go.dev/doc/install) ≥ 1.24 | Build toolchain |
| [OpenCV 4](https://opencv.org/) | Image matching — `brew install opencv` on macOS |
| [pkg-config](https://www.freedesktop.org/wiki/Software/pkg-config/) | Required by gocv — `brew install pkg-config` |

On **macOS** you must also grant the terminal (or the binary) **Accessibility** and **Screen Recording** permissions in *System Settings → Privacy & Security*.

## Quick start

```bash
# 1. Clone and build
git clone https://github.com/activebook/auto-wechat.git
cd auto-wechat
make build

# 2. Calibrate UI positions (follow on-screen prompts)
./auto-wechat scan

# 3. Copy your message to the clipboard, then run
./auto-wechat proc
```

## Configuration — `cfgs/settings.yml`

| Key | Default | Description |
|---|---|---|
| `contacts.x/y` | — | Screen coords of the Contacts button |
| `more.x/y` | — | Screen coords of the More (⋯) button |
| `messages_bar.x/y` | — | Screen coords of the message input bar |
| `messages_popup.margin_x/y` | 20 / -10 | Offset from right-click to "Send Message" menu item |
| `app_title` | `WeChat` | Window title used to activate the app |
| `max_count` | 3 | Maximum number of messages to send per run |
| `interval` | 100 | Delay between actions (ms) |
| `check_end` | `true` | Detect end-of-list via image comparison |
| `auto_send` | `true` | Press Enter after pasting; set `false` to review first |
| `debug` | `false` | Pause 1 s after each action and print step logs |

## Performance Optimization

The tool uses OpenCV for two main purposes:
1. **Calibration**: Locating UI elements during the `scan` phase.
2. **End Detection**: Comparing images to detect when the contact list has ended (`check_end: true`).

If you have already calibrated the positions and set `check_end: false`, the tool will skip all OpenCV image processing and operate purely on coordinate-based automation, which is significantly faster.

## Stop at any time

Press **Ctrl + Shift + Q** while `proc` is running to abort immediately.

## Build & Release

```bash
make build        # dev build → ./auto-wechat
make release      # versioned build + assets → dist/<name>.zip
make test         # go test ./...
make fmt          # go fmt ./...
make tidy         # go mod tidy
make clean        # remove binary and dist/
make help         # list all targets
```

`make release` produces a self-contained zip under `dist/` containing the binary, all reference images (`imgs/`), and the config template (`cfgs/`).

