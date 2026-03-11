package main

import (
	"context"
	"os"
	"os/exec"
	"runtime"
)

type App struct {
	ctx context.Context
}

type BootstrapPayload struct {
	Environment Environment `json:"environment"`
}

type Environment struct {
	Hostname     string `json:"hostname"`
	Platform     string `json:"platform"`
	Architecture string `json:"architecture"`
}

type PostInstallActionResult struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Error     string `json:"error,omitempty"`
	Cancelled bool   `json:"cancelled"`
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetBootstrapPayload() BootstrapPayload {
	return BootstrapPayload{
		Environment: Environment{
			Hostname:     readHostname(),
			Platform:     runtime.GOOS,
			Architecture: runtime.GOARCH,
		},
	}
}

func readHostname() string {
	value, err := os.Hostname()
	if err != nil {
		return "unknown-host"
	}
	return value
}

func detectPowerShell() string {
	if path, err := exec.LookPath("pwsh.exe"); err == nil {
		return path
	}
	if path, err := exec.LookPath("powershell.exe"); err == nil {
		return path
	}
	return "not-found"
}
