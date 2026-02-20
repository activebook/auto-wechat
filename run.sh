#!/bin/bash

# This script is only for macOS
# Exogenous heuristic detection for OpenCV runtime availability
if ! command -v pkg-config >/dev/null 2>&1 || ! pkg-config --exists opencv4; then
    # Fallback to evaluating Homebrew's internal manifest
    if ! command -v brew >/dev/null 2>&1 || ! brew ls --versions opencv >/dev/null 2>&1; then
        echo -e "\033[1;31m===========================================================\033[0m"
        echo -e "\033[1;31m[FATAL EXCEPTION] OpenCV Runtime Unavailability Detected\033[0m"
        echo -e "\033[1;31m===========================================================\033[0m"
        echo "The 'auto-wechat' executable dynamically links to the OpenCV framework,"
        echo "which is currently absent from your macOS operating environment."
        echo ""
        echo "To emancipate this dependency, please provision your system by executing:"
        echo -e "\033[1;32m    brew install opencv pkg-config\033[0m"
        echo -e "\033[1;31m===========================================================\033[0m"
        exit 1
    fi
fi

# Ascertain the absolute execution directory and relinquish control to the binary
EXEC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"

if [ -f "$EXEC_DIR/auto-wechat-macos-arm64" ]; then
    exec "$EXEC_DIR/auto-wechat-macos-arm64" "$@"
elif [ -f "$EXEC_DIR/auto-wechat" ]; then
    exec "$EXEC_DIR/auto-wechat" "$@"
else
    echo -e "\033[1;31m[FATAL EXCEPTION] No executable found.\033[0m"
    echo "Expected 'auto-wechat-macos-arm64' or 'auto-wechat' in: $EXEC_DIR"
    exit 1
fi