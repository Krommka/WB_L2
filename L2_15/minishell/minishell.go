package minishell

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Command содержит название команды, аргументы и редиректы
type Command struct {
	Name      string
	Args      []string
	Redirects []Redirect
}

// Redirect содержит тип редиректа и файл
type Redirect struct {
	Type string // ">", "<", ">>"
	File string
}

// Pipeline содержит слайс команд с операторами и указателем на следующий пайплайн
type Pipeline struct {
	Commands []Command
	Operator string
	Next     *Pipeline
}

// Run содержит основной цикл
func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)

	pipelineCancel := make(chan context.CancelFunc, 1)

	commands := read(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case sig := <-sigCh:
			if sig == syscall.SIGINT {
				fmt.Println("\nПолучен сигнал Ctrl+C - прерывание текущего пайплайна...")
				select {
				case cancelFn := <-pipelineCancel:
					if cancelFn != nil {
						cancelFn()
					}
				default:
				}
			}
		case line, ok := <-commands:
			if !ok {
				return
			}

			if line == "" {
				continue
			}

			pipelineCtx, pipelineCancelFn := context.WithCancel(ctx)

			select {
			case pipelineCancel <- pipelineCancelFn:
			default:
				<-pipelineCancel
				pipelineCancel <- pipelineCancelFn
			}

			go func(ctx context.Context, line string, cancelFn context.CancelFunc) {
				defer cancelFn()

				pipeline, err := parseInput(line)
				if err != nil {
					fmt.Printf("Ошибка парсинга: %v\n", err)
					return
				}

				if pipeline != nil {
					_, err = executePipelineWithContext(ctx, pipeline, true)
					if err != nil {
						if err == context.Canceled {
							fmt.Printf("Пайплайн прерван\n")
						}
					}
				}
			}(pipelineCtx, line, pipelineCancelFn)
		}
	}
}

// read читает ввод от пользователя
func read(ctx context.Context) <-chan string {

	res := make(chan string)
	go func() {
		defer close(res)
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			line := scanner.Text()
			select {
			case <-ctx.Done():
				log.Println("Контекст отменен, завершение чтения")
				return
			case res <- line:
			}
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Ошибка чтения: %v\n", err)
		}
	}()
	return res
}
