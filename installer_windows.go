package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// readRegistryPath reads a PATH-type string value from the Windows registry.
//
//   - For machine-level: subKey = `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`
//   - For user-level:    subKey = `Environment`
func readRegistryPath(subKey, valueName string) (string, error) {
	var root registry.Key
	if strings.HasPrefix(strings.ToUpper(subKey), "SYSTEM") {
		root = registry.LOCAL_MACHINE
	} else {
		root = registry.CURRENT_USER
	}

	k, err := registry.OpenKey(root, subKey, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("open registry key %s: %w", subKey, err)
	}
	defer k.Close()

	val, _, err := k.GetStringValue(valueName)
	if err != nil {
		return "", fmt.Errorf("read registry value %s\\%s: %w", subKey, valueName, err)
	}

	return val, nil
}

// addToUserPath appends a directory to the user's PATH via the registry,
// then refreshes the current process environment.
func addToUserPath(dir string, emit func(string, string, string)) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		if emit != nil {
			emit("path", "warn", fmt.Sprintf("Cannot open user Environment key: %v", err))
		}
		return
	}
	defer k.Close()

	currentPath, _, err := k.GetStringValue("Path")
	if err != nil {
		currentPath = ""
	}

	// Check if dir is already present (case-insensitive)
	parts := strings.Split(currentPath, ";")
	for _, p := range parts {
		if strings.EqualFold(strings.TrimSpace(p), dir) {
			return // already on PATH
		}
	}

	newPath := currentPath
	if newPath != "" && !strings.HasSuffix(newPath, ";") {
		newPath += ";"
	}
	newPath += dir

	if err := k.SetExpandStringValue("Path", newPath); err != nil {
		if emit != nil {
			emit("path", "warn", fmt.Sprintf("Failed to update user PATH: %v", err))
		}
		return
	}

	if emit != nil {
		emit("path", "ok", fmt.Sprintf("Added %s to user PATH (restart terminal if command not found)", dir))
	}

	// Refresh current process
	refreshSystemPath()

	// Also update the current os.Environ so exec.LookPath picks it up immediately
	env := os.Getenv("Path")
	if !strings.Contains(strings.ToLower(env), strings.ToLower(dir)) {
		os.Setenv("Path", env+";"+dir)
	}
}
