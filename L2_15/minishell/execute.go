package minishell

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
)

// executePipelineWithContext  Выполняет пайплайн
func executePipelineWithContext(ctx context.Context, pipeline *Pipeline, prevSuccess bool) (bool, error) {
	if pipeline == nil {
		return prevSuccess, nil
	}

	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	commandsToExecute, nextPipeline, operator := collectCommands(pipeline)

	var success bool
	var err error

	if len(commandsToExecute) == 1 {
		success, err = executeCommandWithContext(ctx, commandsToExecute[0], os.Stdin, os.Stdout, os.Stderr)
	} else if len(commandsToExecute) > 1 {
		success, err = executeCommandsPipelineWithContext(ctx, commandsToExecute)
	} else {
		success = prevSuccess
	}

	if err != nil && pipeline.Next == nil {
		return false, err
	}

	if nextPipeline != nil {
		switch operator {
		case "&&":
			if success {
				return executePipelineWithContext(ctx, nextPipeline, success)
			}
			return false, nil

		case "||":
			if !success {
				return executePipelineWithContext(ctx, nextPipeline, success)
			}
			return true, nil

		default:
			return executePipelineWithContext(ctx, nextPipeline, success)
		}
	}

	return success, nil
}

//
//func executePipelineWithContext(ctx context.Context, pipeline *Pipeline, prevSuccess bool) (bool, error) {
//	if pipeline == nil {
//		return prevSuccess, nil
//	}
//
//	select {
//	case <-ctx.Done():
//		return false, ctx.Err()
//	default:
//	}
//
//	var success bool
//	var err error
//
//	if len(pipeline.Commands) > 0 {
//		if len(pipeline.Commands) == 1 {
//			success, err = executeCommandWithContext(ctx, pipeline.Commands[0], os.Stdin, os.Stdout, os.Stderr)
//		} else {
//			success, err = executeCommandsPipelineWithContext(ctx, pipeline.Commands)
//		}
//	} else {
//		success = prevSuccess
//	}
//
//	if err != nil {
//		return false, err
//	}
//	//if err != nil && pipeline.Next == nil {
//	//	return false, err
//	//}
//
//	if pipeline.Next != nil {
//		switch pipeline.Operator {
//		case "&&":
//			if success {
//				return executePipelineWithContext(ctx, pipeline.Next, success)
//			}
//			return false, nil
//		case "||":
//			if !success {
//				return executePipelineWithContext(ctx, pipeline.Next, success)
//			}
//			return true, nil
//		default:
//			return executePipelineWithContext(ctx, pipeline.Next, success)
//		}
//	}
//
//	return success, nil
//}

// executeCommandsPipelineWithContext Запускает команды пайплайна
func executeCommandsPipelineWithContext(ctx context.Context, commands []Command) (bool, error) {
	if len(commands) == 0 {
		return true, nil
	}
	pipes := make([]*io.PipeWriter, len(commands)-1)
	readers := make([]*io.PipeReader, len(commands)-1)

	for i := 0; i < len(commands)-1; i++ {
		readers[i], pipes[i] = io.Pipe()
	}

	errors := make(chan error, len(commands))
	var wg sync.WaitGroup

	for i := 0; i < len(commands); i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			var stdin io.Reader = os.Stdin
			var stdout io.Writer = os.Stdout

			if i > 0 {
				stdin = readers[i-1]
			}

			if i < len(commands)-1 {
				stdout = pipes[i]
			}

			_, err := executeCommandWithContext(ctx, commands[i], stdin, stdout, os.Stderr)
			if err != nil {
				errors <- fmt.Errorf("%s: %v", commands[i].Name, err)
			} else {
				errors <- nil
			}

			if i < len(commands)-1 {
				pipes[i].Close()
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

// collectCommands Собирает команды и определяет следующий пайплайн
func collectCommands(pipeline *Pipeline) ([]Command, *Pipeline, string) {
	var commands []Command
	current := pipeline
	var next *Pipeline
	var operator string

	for current != nil {
		commands = append(commands, current.Commands...)

		if current.Operator == "|" && current.Next != nil {
			current = current.Next
		} else {
			next = current.Next
			operator = current.Operator
			break
		}
	}

	return commands, next, operator
}

// executeCommandWithContext Запускает выполнение builtins или внешней подпрограммы
func executeCommandWithContext(ctx context.Context, cmd Command, stdin io.Reader, stdout io.Writer, stderr io.Writer) (bool, error) {
	stdin, stdout, stderr, cleanup, err := setupRedirects(cmd, stdin, stdout, stderr)
	if err != nil {
		return false, err
	}
	defer cleanup()

	switch cmd.Name {
	case "cd":
		return executeCd(cmd.Args, stdout, stderr)
	case "pwd":
		return executePwd(cmd.Args, stdout, stderr)
	case "echo":
		return executeEcho(cmd.Args, stdout, stderr)
	case "kill":
		return executeKill(cmd.Args, stdout, stderr)
	case "ps":
		return executePs(cmd.Args, stdout, stderr)
	default:
		return executeExternalWithContext(ctx, cmd, stdin, stdout, stderr)
	}
}

// setupRedirects Устанавливает редиректы > < >>
func setupRedirects(cmd Command, stdin io.Reader, stdout io.Writer, stderr io.Writer) (io.Reader, io.Writer, io.Writer, func(), error) {
	var files []*os.File
	cleanup := func() {
		for _, f := range files {
			if f != nil {
				f.Close()
			}
		}
	}

	for _, redirect := range cmd.Redirects {
		switch redirect.Type {
		case ">":
			file, err := os.OpenFile(redirect.File, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return nil, nil, nil, cleanup, err
			}
			files = append(files, file)
			stdout = file

		case ">>":
			file, err := os.OpenFile(redirect.File, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				return nil, nil, nil, cleanup, err
			}
			files = append(files, file)
			stdout = file

		case "<":
			file, err := os.Open(redirect.File)
			if err != nil {
				return nil, nil, nil, cleanup, err
			}
			files = append(files, file)
			stdin = file
		}
	}

	return stdin, stdout, stderr, cleanup, nil
}
