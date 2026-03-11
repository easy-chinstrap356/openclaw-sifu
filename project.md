# OpenClaw-Sifu

## Goal

Build a portable Windows desktop assistant that launches from a single executable and helps users complete installation flows under guidance from a cloud-hosted generic agent.

## Product stance

- The local application is only the shell and execution bridge.
- The cloud side is agent-oriented and model-agnostic for now.
- OpenClaw-specific system prompts, policies, and flow tuning are deferred to a later phase.

## Local responsibilities

- Collect machine facts and runtime status.
- Capture screenshots and other task context.
- Execute local commands and installer actions.
- Resume interrupted sessions after failures or reboot.

## Cloud responsibilities

- Interpret screenshots and machine state.
- Decide the next action step.
- Return structured actions for the local shell to execute.
- Maintain per-user installation session state.

## Current milestone

Create a Wails desktop shell that can:

- boot as a portable executable
- show local environment details
- expose a placeholder generic-agent strategy
- reserve extension points for screenshot, execution, and session orchestration
