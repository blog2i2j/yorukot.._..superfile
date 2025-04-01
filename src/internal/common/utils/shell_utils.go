package utils

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"runtime"
	"time"
)

// Choose correct shell as per OS
func ExecuteCommandInShell(timeLimitMsec int, cmdDir string, shellCommand string) (int, string, error) {
	// Linux and Darwin
	baseCmd := "/bin/sh"
	args := []string{"-c", shellCommand}

	if runtime.GOOS == "windows" {
		baseCmd = "powershell.exe"
		args[0] = "-Command"
	}

	return ExecuteCommand(timeLimitMsec, cmdDir, baseCmd, args...)
}

func ExecuteCommand(timeLimitMsec int, cmdDir string, baseCmd string, args ...string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeLimitMsec)*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, baseCmd, args...)
	cmd.Dir = cmdDir
	outputBytes, err := cmd.CombinedOutput()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		slog.Error("User's command timed out")
	}

	retCode := -1
	if err == nil {
		retCode = 0
	} else if exitErr, ok := err.(*exec.ExitError); ok {
		// We dont expect error to be Wrapped here, so we are using type
		// assertion not errors.As
		retCode = exitErr.ExitCode()
	} else {
		err = fmt.Errorf("unexpected Error in command execution : %w", err)
	}

	return retCode, string(outputBytes), err
}
