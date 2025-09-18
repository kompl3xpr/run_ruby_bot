package core

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"
)

type InterpretStatus int

const (
	SuccessWithOutput InterpretStatus = iota
	SuccessWithoutOutput
	CodeError
	InternalError
	Timeout
)

type InterpretResult struct {
	Msg    string
	Status InterpretStatus
}

func InitInterpreter() {
	slog.Info("removing old docker image...")
	cmd := exec.Command("docker", "rmi", Cfg().Docker.Name)
	if err := cmd.Run(); err != nil {
		slog.Warn("failed to remove image, ", "error", err.Error())
	}

	slog.Info("building new docker image...")
	cmd = exec.Command("docker", "build", "-t", Cfg().Docker.Name, ".")
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	if err != nil {
		slog.Error(
			"failed to rebuild image, ",
			"error", err.Error(),
			"stderr", stderrBuf.String(),
		)
		os.Exit(4)
	}

	slog.Info("docker image rebuilt")
}

func Interpret(code string) InterpretResult {
	f, err := os.CreateTemp("", "run_ruby_bot_script_*.rb")
	if err != nil {
		return InterpretResult{err.Error(), InternalError}
	}
	fPath := f.Name()
	if err := writeScript(f, code); err != nil {
		f.Close()
		os.Remove(fPath)
		return InterpretResult{err.Error(), InternalError}
	}
	f.Close()

	cmd := exec.Command(
		"docker", "run", "--rm",
		fmt.Sprintf("--memory=%v", Cfg().Docker.Memory),
		fmt.Sprintf("--cpus=%v", Cfg().Docker.Cpus),
		"-v", fmt.Sprintf("%v:/app/main.rb", fPath),
		Cfg().Docker.Name,
	)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	done := make(chan error)
	go func() {
		result := cmd.Run()
		os.Remove(fPath)
		done <- result
	}()

	select {
	case err := <-done:
		if _, ok := err.(*exec.ExitError); ok {
			return InterpretResult{stderrBuf.String(), CodeError}
		}
		stdout := stdoutBuf.String()
		if len(strings.TrimSpace(stdout)) == 0 {
			return InterpretResult{stdout, SuccessWithoutOutput}
		}
		return InterpretResult{stdout, SuccessWithOutput}

	case <-time.After(Cfg().Task.Timeout * time.Second):
		cmd.Process.Kill()
		<-done
		return InterpretResult{"", Timeout}
	}
}

func writeScript(f *os.File, code string) error {
	_, err := f.WriteString("puts -> do\n")
	if err != nil {
		return err
	}
	_, err = f.WriteString(code)
	if err != nil {
		return err
	}
	_, err = f.WriteString("\nend.call")
	if err != nil {
		return err
	}
	return nil
}
