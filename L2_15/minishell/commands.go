package minishell

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// executeCd Выполняет аналог программы смены директории
func executeCd(args []string, stdout io.Writer, stderr io.Writer) (bool, error) {
	var path string
	if len(args) == 0 || (len(args) == 1 && args[0] == "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(stderr, "cd: %v\n", err)
			return false, err
		}
		path = home
	} else {
		path = args[0]
	}

	if err := os.Chdir(path); err != nil {
		fmt.Fprintf(stderr, "cd: %v\n", err)
		return false, err
	}

	return true, nil
}

// executePwd Выполняет аналог команды вывода текущей директории
func executePwd(args []string, stdout io.Writer, stderr io.Writer) (bool, error) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(stderr, "pwd: %v\n", err)
		return false, err
	}
	fmt.Fprintln(stdout, dir)
	return true, nil
}

func executeEcho(args []string, stdout io.Writer, stderr io.Writer) (bool, error) {
	fmt.Fprintln(stdout, strings.Join(args, " "))
	return true, nil
}

// executeKill Выполняет аналог команды закрытия процесса по номеру PID
func executeKill(args []string, stdout io.Writer, stderr io.Writer) (bool, error) {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "kill: требуется PID процесса")
		return false, fmt.Errorf("kill: требуется PID процесса")
	}

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(stderr, "kill: неверный PID: %v\n", err)
		return false, fmt.Errorf("kill: неверный PID: %v", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Fprintf(stderr, "kill: процесс не найден: %v\n", err)
		return false, fmt.Errorf("kill: процесс не найден: %v", err)
	}

	if err = process.Signal(syscall.SIGTERM); err != nil {
		fmt.Fprintf(stderr, "kill: не удалось послать сигнал: %v\n", err)
		return false, fmt.Errorf("kill: не удалось послать сигнал: %v", err)
	}

	fmt.Fprintf(stdout, "Сигнал завершения послан процессу %d\n", pid)
	return true, nil
}

// executePs Выполняет аналог команды вывода текущих процессов
func executePs(args []string, stdout io.Writer, stderr io.Writer) (bool, error) {
	cmd := exec.Command("ps", "aux")
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(stderr, "ps: %v\n", err)
		return false, fmt.Errorf("ps: %v", err)
	}

	return true, nil
}

// executeExternalWithContext Выполняет внешние команды с CommandContext
func executeExternalWithContext(ctx context.Context, cmd Command, stdin io.Reader, stdout io.Writer, stderr io.Writer) (bool, error) {
	command := exec.CommandContext(ctx, cmd.Name, cmd.Args...)
	command.Stdin = stdin
	command.Stdout = stdout
	command.Stderr = stderr

	if err := command.Run(); err != nil {
		if ctx.Err() == context.Canceled {
			return false, context.Canceled
		}
		if _, ok := err.(*exec.ExitError); ok && len(cmd.Redirects) > 0 {
			return false, nil
		}
		fmt.Fprintf(stderr, "%s: %v\n", cmd.Name, err)
		return false, fmt.Errorf("%s: %v", cmd.Name, err)
	}

	return true, nil
}
