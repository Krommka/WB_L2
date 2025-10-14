package minishell

import (
	"fmt"
	"os"
	"strings"
)

// parseInput Парсит строку от пользователя
func parseInput(line string) (*Pipeline, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil
	}

	parts := splitByOperators(line)
	return parsePipeline(parts, 0)
}

// splitByOperators делит строку на части с обработкой кавычек
func splitByOperators(line string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)
	escape := false

	for i := 0; i < len(line); i++ {
		c := line[i]

		if escape {
			current.WriteByte(c)
			escape = false
			continue
		}

		if c == '\\' {
			escape = true
			continue
		}

		if c == '"' || c == '\'' {
			if !inQuotes {
				inQuotes = true
				quoteChar = c
			} else if c == quoteChar {
				inQuotes = false
			}
			current.WriteByte(c)
			continue
		}

		if !inQuotes && (c == '|' || c == '&') {
			if current.Len() > 0 {
				parts = append(parts, strings.TrimSpace(current.String()))
				current.Reset()
			}

			if i+1 < len(line) && line[i+1] == c {
				parts = append(parts, string(c)+string(c))
				i++
			} else if c == '|' {
				parts = append(parts, "|")
			} else {
				current.WriteByte(c)
			}
			continue
		}

		current.WriteByte(c)
	}
	if current.Len() > 0 {
		parts = append(parts, strings.TrimSpace(current.String()))
	}

	return parts
}

// parsePipeline Парсит части в пайплайн
func parsePipeline(parts []string, index int) (*Pipeline, error) {
	if index >= len(parts) {
		return nil, nil
	}

	pipe := &Pipeline{}

	for index < len(parts) {
		part := parts[index]

		if part == "|" || part == "&&" || part == "||" {
			pipe.Operator = part
			next, err := parsePipeline(parts, index+1)
			if err != nil {
				return nil, err
			}
			pipe.Next = next
			break
		}

		cmd, err := parseCommand(part)
		if err != nil {
			return nil, err
		}
		pipe.Commands = append(pipe.Commands, cmd)
		index++
	}

	return pipe, nil
}

// parseCommand создает команду
func parseCommand(cmdStr string) (Command, error) {
	var tokens []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)
	escape := false

	for i := 0; i < len(cmdStr); i++ {
		c := cmdStr[i]

		if escape {
			current.WriteByte(c)
			escape = false
			continue
		}

		switch c {
		case '\\':
			escape = true
		case '"', '\'':
			if !inQuotes {
				inQuotes = true
				quoteChar = c
			} else if c == quoteChar {
				inQuotes = false
			} else {
				current.WriteByte(c)
			}
		case '>', '<':
			if !inQuotes {
				if current.Len() > 0 {
					tokens = append(tokens, current.String())
					current.Reset()
				}

				if c == '>' && i+1 < len(cmdStr) && cmdStr[i+1] == '>' {
					tokens = append(tokens, ">>")
					i++
				} else {
					tokens = append(tokens, string(c))
				}
			} else {
				current.WriteByte(c)
			}
		case ' ', '\t':
			if inQuotes {
				current.WriteByte(c)
			} else {
				if current.Len() > 0 {
					tokens = append(tokens, current.String())
					current.Reset()
				}
			}
		default:
			current.WriteByte(c)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return parseTokens(tokens)
}

// parseTokens Парсит токены в команду
func parseTokens(tokens []string) (Command, error) {
	var args []string
	var redirects []Redirect

	i := 0
	for i < len(tokens) {
		token := tokens[i]

		if token == ">" || token == "<" || token == ">>" {
			if i+1 >= len(tokens) {
				return Command{}, fmt.Errorf("редирект %s требует файл", token)
			}

			file := expandEnvVars(tokens[i+1])

			redirects = append(redirects, Redirect{Type: token, File: file})
			i += 2
		} else {
			expanded := expandEnvVars(token)
			args = append(args, expanded)
			i++
		}
	}

	if len(args) == 0 {
		return Command{}, fmt.Errorf("пустая команда")
	}

	return Command{
		Name:      args[0],
		Args:      args[1:],
		Redirects: redirects,
	}, nil
}

// expandEnvVars возвращает значение переменных окружения при наличии
func expandEnvVars(s string) string {
	return os.Expand(s, func(key string) string {
		if value, exists := os.LookupEnv(key); exists {
			return value
		}
		return "$" + key
	})
}
