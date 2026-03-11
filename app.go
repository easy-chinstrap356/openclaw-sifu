package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type App struct {
	ctx context.Context
	mu  sync.RWMutex
}

type BootstrapPayload struct {
	AppName      string         `json:"appName"`
	Version      string         `json:"version"`
	ModeLabel    string         `json:"modeLabel"`
	Summary      string         `json:"summary"`
	BootTime     string         `json:"bootTime"`
	Agent        AgentProfile   `json:"agent"`
	Environment  Environment    `json:"environment"`
	Capabilities []Capability   `json:"capabilities"`
	Pipeline     []PipelineStep `json:"pipeline"`
	NextSteps    []string       `json:"nextSteps"`
}

type AgentProfile struct {
	Name         string `json:"name"`
	Summary      string `json:"summary"`
	EndpointMode string `json:"endpointMode"`
	PromptLayer  string `json:"promptLayer"`
}

type Environment struct {
	Hostname       string `json:"hostname"`
	Username       string `json:"username"`
	Platform       string `json:"platform"`
	Architecture   string `json:"architecture"`
	GoVersion      string `json:"goVersion"`
	WorkingDir     string `json:"workingDir"`
	ExecutablePath string `json:"executablePath"`
	TempDir        string `json:"tempDir"`
	PowerShellPath string `json:"powerShellPath"`
	WebView2State  string `json:"webView2State"`
}

type Capability struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Status string `json:"status"`
}

type PipelineStep struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	State       string `json:"state"`
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetBootstrapPayload() BootstrapPayload {
	return BootstrapPayload{
		AppName:   "OpenClaw-Sifu",
		Version:   "0.1.0",
		ModeLabel: "Generic Agent Bootstrap",
		Summary:   "A local assistant shell for remote agent-guided installation, validation, and recovery flows.",
		BootTime:  time.Now().Format(time.RFC3339),
		Agent: AgentProfile{
			Name:         "MiniMax Cloud Assistant",
			Summary:      "The cloud-side guide is now wired for MiniMax, while OpenClaw-specific prompt tuning still remains a later layer.",
			EndpointMode: "MiniMax chat completion + local executor bridge",
			PromptLayer:  "Installation playbook first, OpenClaw specialization later",
		},
		Environment: Environment{
			Hostname:       readHostname(),
			Username:       readUsername(),
			Platform:       runtime.GOOS,
			Architecture:   runtime.GOARCH,
			GoVersion:      runtime.Version(),
			WorkingDir:     readWorkingDir(),
			ExecutablePath: readExecutablePath(),
			TempDir:        os.TempDir(),
			PowerShellPath: detectPowerShell(),
			WebView2State:  detectWebView2(),
		},
		Capabilities: []Capability{
			{
				ID:     "screen-capture",
				Title:  "Screen and context capture",
				Detail: "Ready for screenshot ingestion, window metadata collection, and later OCR attachment.",
				Status: "ready",
			},
			{
				ID:     "local-execution",
				Title:  "Local execution bridge",
				Detail: "Designed for PowerShell commands, installer launch, and machine-state validation routines.",
				Status: "ready",
			},
			{
				ID:     "remote-session",
				Title:  "Remote agent session",
				Detail: "Reserved for websocket or gRPC orchestration once the cloud agent endpoint is wired in.",
				Status: "planned",
			},
			{
				ID:     "installer-recovery",
				Title:  "Recovery and resume",
				Detail: "Intended for reboot-safe checkpoints and resumable installation steps.",
				Status: "planned",
			},
			{
				ID:     "openclaw-layer",
				Title:  "OpenClaw prompt pack",
				Detail: "Explicitly postponed. The current app keeps the orchestration layer generic by design.",
				Status: "deferred",
			},
		},
		Pipeline: []PipelineStep{
			{
				Title:       "Bootstrap local shell",
				Description: "Launch a portable desktop shell without running a traditional installer flow.",
				State:       "ready",
			},
			{
				Title:       "Collect machine snapshot",
				Description: "Read environment facts so the remote agent can reason about the target system.",
				State:       "ready",
			},
			{
				Title:       "Connect remote agent",
				Description: "Attach the cloud-side generic agent pool and start session orchestration.",
				State:       "next",
			},
			{
				Title:       "Attach installation toolset",
				Description: "Expose commands for installer launch, verification, rollback, and recovery logic.",
				State:       "planned",
			},
			{
				Title:       "Layer OpenClaw specialization",
				Description: "Add OpenClaw-specific prompts and policies after the generic agent path is stable.",
				State:       "deferred",
			},
		},
		NextSteps: []string{
			"Wire a cloud session endpoint and handshake contract.",
			"Implement screenshot upload, command execution, and result streaming.",
			"Split elevated actions into a dedicated helper if admin flows become necessary.",
			"Add OpenClaw-specific system prompts only after the generic agent path is validated.",
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

func readUsername() string {
	value, err := os.UserHomeDir()
	if err != nil {
		return "unknown-user"
	}
	return filepath.Base(value)
}

func readWorkingDir() string {
	value, err := os.Getwd()
	if err != nil {
		return "unavailable"
	}
	return value
}

func readExecutablePath() string {
	value, err := os.Executable()
	if err != nil {
		return "unavailable"
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

func detectWebView2() string {
	candidates := []string{
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "Microsoft", "EdgeWebView", "Application"),
		filepath.Join(os.Getenv("ProgramFiles"), "Microsoft", "EdgeWebView", "Application"),
	}

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return "detected"
		}
	}

	return "not-detected"
}
