# OpenClaw-Sifu

OpenClaw-Sifu is a Wails-based desktop installer shell designed for a fluid, local installation experience for OpenClaw. It runs completely locally on the user's machine to guide them through prerequisite checks, component acquisition, and initialization of the OpenClaw environment.

## Current direction

- The desktop app functions as an automated local shell installer.
- It detects the native OS environment (Windows, macOS, Linux) and executes the corresponding deterministic setup steps.

## Stack

- Wails v2
- React + Vite + TypeScript
- TailwindCSS
- Go backend for execution mapping and robust cross-platform system calls

## Local development

```bash
# Start Wails dev server (hot reload for both frontend and Go)
wails dev
```

## Production build

```bash
# Build a standalone executable
wails build
```
