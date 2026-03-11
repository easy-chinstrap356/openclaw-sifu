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

// LaunchOpenClaw launches the gateway in a new terminal window and opens the dashboard UI.
func (a *App) LaunchOpenClaw() error {
	if err := ensureGatewayModeLocal(); err != nil {
		return err
	}

	if err := startGatewayService(); err != nil {
		if err := launchGatewayFallback(); err != nil {
			return err
		}
	}

	if err := waitForGatewayReady(20 * time.Second); err != nil {
		return err
	}

	// Now open the dashboard (generates token URL + browser launch)
	dashboardCmd := exec.Command("openclaw", "dashboard")
	if err := dashboardCmd.Start(); err != nil {
		return fmt.Errorf("failed to open dashboard: %w", err)
	}

	return nil
}

func ensureGatewayModeLocal() error {
	cmd := exec.Command("openclaw", "config", "set", "gateway.mode", "local")
	output, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}

	message := sanitizeOutput(decodeOutputBytes(output))
	if message == "" {
		return fmt.Errorf("failed to set gateway.mode=local: %w", err)
	}
	return fmt.Errorf("failed to set gateway.mode=local: %s", message)
}

func startGatewayService() error {
	cmd := exec.Command("openclaw", "gateway", "start")
	output, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}

	message := sanitizeOutput(decodeOutputBytes(output))
	if message == "" {
		return fmt.Errorf("failed to start gateway service: %w", err)
	}
	return fmt.Errorf("failed to start gateway service: %s", message)
}

func launchGatewayFallback() error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		shellPath := detectPowerShell()
		if shellPath == "not-found" {
			shellPath = "powershell.exe"
		}
		cmd = exec.Command(
			shellPath,
			"-NoProfile",
			"-NonInteractive",
			"-WindowStyle",
			"Hidden",
			"-Command",
			"Start-Process -WindowStyle Hidden -FilePath 'cmd.exe' -ArgumentList @('/d','/c','openclaw','gateway','run','--allow-unconfigured')",
		)
	case "darwin":
		cmd = exec.Command("osascript", "-e", `tell application "Terminal" to do script "openclaw gateway run --allow-unconfigured"`)
	default:
		cmd = exec.Command("x-terminal-emulator", "-e", "openclaw gateway run --allow-unconfigured")
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to launch gateway process: %w", err)
	}

	return nil
}

func waitForGatewayReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		out, err := execOutput("openclaw", "daemon", "status", "--json")
		if err == nil && strings.Contains(out, `"ok": true`) {
			return nil
		}
		time.Sleep(time.Second)
	}

	return fmt.Errorf("gateway did not become ready before opening the dashboard")
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
