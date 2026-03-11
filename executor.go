package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unicode/utf16"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
)

var ansiSequencePattern = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`)

func (a *App) OpenExternalTarget(target string) error {
	target = strings.TrimSpace(target)
	if target == "" {
		return fmt.Errorf("target is required")
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", target)
	case "darwin":
		cmd = exec.Command("open", target)
	default:
		cmd = exec.Command("xdg-open", target)
	}

	return cmd.Start()
}

// LaunchOpenClaw launches the gateway in a new terminal window and opens the dashboard UI.
func (a *App) LaunchOpenClaw() error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// Open a new command prompt, name the window, and run openclaw gateway
		cmd = exec.Command("cmd", "/c", "start", "OpenClaw Gateway", "cmd", "/k", "openclaw", "gateway")
	case "darwin":
		// On macOS, try to open Terminal and run the command
		cmd = exec.Command("osascript", "-e", `tell application "Terminal" to do script "openclaw gateway"`)
	default:
		// On Linux, try x-terminal-emulator
		cmd = exec.Command("x-terminal-emulator", "-e", "openclaw gateway")
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start gateway terminal: %w", err)
	}

	// Give the gateway a moment to start and write config strings if needed
	time.Sleep(1500 * time.Millisecond)

	// Now open the dashboard (generates token URL + browser launch)
	dashboardCmd := exec.Command("openclaw", "dashboard")
	if err := dashboardCmd.Start(); err != nil {
		return fmt.Errorf("failed to open dashboard: %w", err)
	}

	return nil
}

func sanitizeOutput(value string) string {
	value = ansiSequencePattern.ReplaceAllString(value, "")
	value = strings.ReplaceAll(value, "\u0007", "")
	value = strings.ReplaceAll(value, "\u0000", "")
	return strings.TrimSpace(value)
}

func decodeOutputBytes(input []byte) string {
	cleaned := bytes.TrimSpace(bytes.Clone(input))
	if len(cleaned) == 0 {
		return ""
	}

	if utf8.Valid(cleaned) {
		return string(cleaned)
	}

	if looksLikeUTF16LE(cleaned) {
		return decodeUTF16LE(cleaned)
	}

	if decoded, err := simplifiedchinese.GB18030.NewDecoder().Bytes(cleaned); err == nil && utf8.Valid(decoded) {
		return string(decoded)
	}

	return string(cleaned)
}

func looksLikeUTF16LE(input []byte) bool {
	if len(input) < 2 || len(input)%2 != 0 {
		return false
	}

	zeroCount := 0
	for i := 1; i < len(input); i += 2 {
		if input[i] == 0 {
			zeroCount++
		}
	}

	return zeroCount >= len(input)/4
}

func decodeUTF16LE(input []byte) string {
	if len(input)%2 != 0 {
		input = input[:len(input)-1]
	}

	units := make([]uint16, 0, len(input)/2)
	for i := 0; i < len(input); i += 2 {
		units = append(units, binary.LittleEndian.Uint16(input[i:i+2]))
	}

	return string(utf16.Decode(units))
}
