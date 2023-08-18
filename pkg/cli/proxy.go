package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/sulaiman-coder/go-error-with-exit-code/ecerror"
)

type Runner struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

var errKhulnasoftCantBeExecuted = errors.New(`the command "khulnasoft" can't be executed via khulnasoft-proxy to prevent the infinite loop`)

func (runner *Runner) Run(ctx context.Context, args ...string) error {
	cmdName := filepath.Base(args[0])
	if cmdName == "khulnasoft" {
		fmt.Fprintln(os.Stderr, "[ERROR] "+errKhulnasoftCantBeExecuted.Error())
		return errKhulnasoftCantBeExecuted
	}
	cmd := exec.CommandContext(ctx, "khulnasoft", append([]string{"exec", "--", cmdName}, args[1:]...)...) //nolint:gosec
	cmd.Stdin = runner.Stdin
	cmd.Stdout = runner.Stdout
	cmd.Stderr = runner.Stderr

	setCancel(cmd)

	if err := cmd.Run(); err != nil {
		return ecerror.Wrap(err, cmd.ProcessState.ExitCode())
	}
	return nil
}

const waitDelay = 1000 * time.Hour

func setCancel(cmd *exec.Cmd) {
	cmd.Cancel = func() error {
		return cmd.Process.Signal(os.Interrupt) //nolint:wrapcheck
	}
	cmd.WaitDelay = waitDelay
}

func absoluteKhulnasoftPath() (string, error) {
	khulnasoftPath, err := exec.LookPath("khulnasoft")
	if err != nil {
		return "", fmt.Errorf("khulnasoft isn't found: %w", err)
	}
	if filepath.IsAbs(khulnasoftPath) {
		return khulnasoftPath, nil
	}
	a, err := filepath.Abs(khulnasoftPath)
	if err != nil {
		return "", fmt.Errorf(`convert relative path "%s" to absolute path: %w`, khulnasoftPath, err)
	}
	return a, nil
}
