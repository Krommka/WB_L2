package minishell

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func TestBasicCommands(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "pwd command",
			input:    "pwd",
			expected: getCurrentDir(),
		},
		{
			name:     "echo command",
			input:    "echo hello world",
			expected: "hello world\n",
		},
		{
			name:     "cd command",
			input:    "cd /tmp && pwd",
			expected: "/tmp\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			pipeline, err := parseInput(tt.input)
			if err != nil {
				t.Fatalf("parseInput failed: %v", err)
			}

			success, err := executePipelineWithContext(ctx, pipeline, true)

			if (err != nil) != tt.wantErr {
				t.Errorf("executePipelineWithContext() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !success && !tt.wantErr {
				t.Error("pipeline execution failed")
			}
		})
	}
}

func TestPipelines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "simple pipeline",
			input:    "echo hello | grep hello",
			contains: "hello",
		},
		{
			name:     "multiple pipes",
			input:    "echo -e 'hello\\nworld' | grep hello | wc -l",
			contains: "1",
		},
		{
			name:     "ps with grep",
			input:    "ps aux | grep $$ | wc -l",
			contains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			pipeline, err := parseInput(tt.input)
			if err != nil {
				t.Fatalf("parseInput failed: %v", err)
			}

			var output bytes.Buffer
			originalStdout := os.Stdout
			originalStderr := os.Stderr

			r, w, _ := os.Pipe()
			os.Stdout = w
			os.Stderr = w

			success, err := executePipelineWithContext(ctx, pipeline, true)

			w.Close()
			os.Stdout = originalStdout
			os.Stderr = originalStderr

			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output.Write(buf[:n])

			if err != nil {
				t.Errorf("executePipelineWithContext() error = %v", err)
			}

			if !success {
				t.Error("pipeline execution failed")
			}

			if tt.contains != "" && !strings.Contains(output.String(), tt.contains) {
				t.Errorf("output %q doesn't contain %q", output.String(), tt.contains)
			}
		})
	}
}

func TestRedirects(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		checkFile   string
		shouldExist bool
		contains    string
		lineCount   int
		wantErr     bool
	}{
		{
			name:        "output redirect should create file",
			input:       "echo test content > test_output.txt",
			checkFile:   "test_output.txt",
			shouldExist: true,
			contains:    "test content",
		},
		{
			name:        "output redirect should overwrite file",
			input:       "echo new content > test_output.txt",
			checkFile:   "test_output.txt",
			shouldExist: true,
			contains:    "new content",
		},
		{
			name:        "append redirect should add content",
			input:       "echo line1 >> test_append.txt && echo line2 >> test_append.txt",
			checkFile:   "test_append.txt",
			shouldExist: true,
			contains:    "line1",
			lineCount:   2,
		},
		{
			name:        "pwd redirect to file",
			input:       "pwd > pwd_output.txt",
			checkFile:   "pwd_output.txt",
			shouldExist: true,
		},
		{
			name:        "redirect to non-existent directory should fail",
			input:       "echo test > /nonexistent/dir/file.txt",
			checkFile:   "/nonexistent/dir/file.txt",
			shouldExist: false,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Remove(tt.checkFile)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			pipeline, err := parseInput(tt.input)
			if err != nil {
				t.Fatalf("parseInput failed: %v", err)
			}

			success, err := executePipelineWithContext(ctx, pipeline, true)

			_, statErr := os.Stat(tt.checkFile)
			fileExists := statErr == nil

			if tt.shouldExist && !fileExists {
				t.Errorf("file %s should exist but doesn't", tt.checkFile)
			}

			if !tt.shouldExist && fileExists {
				t.Errorf("file %s should not exist but does", tt.checkFile)
			}

			if fileExists && tt.shouldExist {
				content, err := os.ReadFile(tt.checkFile)
				if err != nil {
					t.Errorf("failed to read file: %v", err)
				}

				contentStr := string(content)
				if tt.contains != "" && !strings.Contains(contentStr, tt.contains) {
					t.Errorf("file content %q doesn't contain %q", contentStr, tt.contains)
				}

				if tt.lineCount > 0 {
					lines := strings.Split(strings.TrimSpace(contentStr), "\n")
					if len(lines) != tt.lineCount {
						t.Errorf("expected %d lines, got %d", tt.lineCount, len(lines))
					}
				}
			}

			if tt.wantErr {
				if err == nil && success {
					t.Error("expected error but got success")
				}
			} else {
				if err != nil {
					t.Errorf("executePipelineWithContext() error = %v", err)
				}
				if !success {
					t.Error("pipeline execution failed")
				}
			}

			os.Remove(tt.checkFile)
		})
	}
}

func TestOperators(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expected       bool
		expectError    bool
		shouldPrint    string
		shouldNotPrint string
		shouldError    string
	}{
		{
			name:        "AND operator success - both commands should execute",
			input:       "pwd && echo success",
			expected:    true,
			shouldPrint: "success",
		},
		{
			name:           "AND operator failure - second command should NOT execute",
			input:          "cd nonexistent_dir && echo should_not_print",
			expected:       false,
			expectError:    false,
			shouldError:    "cd:",
			shouldNotPrint: "should_not_print",
		},
		{
			name:        "OR operator success - second command should execute",
			input:       "cd nonexistent_dir || echo fallback",
			expected:    true,
			expectError: false,
			shouldError: "cd:",
			shouldPrint: "fallback",
		},
		{
			name:           "OR operator failure - second command should NOT execute",
			input:          "pwd || echo should_not_print",
			expected:       true,
			shouldNotPrint: "should_not_print",
		},
		{
			name:        "single command with error should return error",
			input:       "cd nonexistent_dir",
			expected:    false,
			expectError: true,
			shouldError: "cd:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			pipeline, err := parseInput(tt.input)
			if err != nil {
				t.Fatalf("parseInput failed: %v", err)
			}

			var stdoutBuf, stderrBuf bytes.Buffer
			originalStdout := os.Stdout
			originalStderr := os.Stderr

			rOut, wOut, _ := os.Pipe()
			rErr, wErr, _ := os.Pipe()
			os.Stdout = wOut
			os.Stderr = wErr

			success, err := executePipelineWithContext(ctx, pipeline, true)

			wOut.Close()
			wErr.Close()
			os.Stdout = originalStdout
			os.Stderr = originalStderr

			stdoutBuf.ReadFrom(rOut)
			stderrBuf.ReadFrom(rErr)

			stdoutStr := stdoutBuf.String()
			stderrStr := stderrBuf.String()

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("executePipelineWithContext() error = %v", err)
				}
			}

			if success != tt.expected {
				t.Errorf("expected success=%v, got %v", tt.expected, success)
			}

			if tt.shouldError != "" && !strings.Contains(stderrStr, tt.shouldError) {
				t.Errorf("expected stderr to contain %q, but got %q", tt.shouldError, stderrStr)
			}

			if tt.shouldPrint != "" && !strings.Contains(stdoutStr, tt.shouldPrint) {
				t.Errorf("expected stdout to contain %q, but got %q", tt.shouldPrint, stdoutStr)
			}

			if tt.shouldNotPrint != "" && strings.Contains(stdoutStr, tt.shouldNotPrint) {
				t.Errorf("expected stdout NOT to contain %q, but got %q", tt.shouldNotPrint, stdoutStr)
			}
		})
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "unknown command",
			input:   "nonexistent_command_xyz",
			wantErr: true,
		},
		{
			name:    "invalid redirect",
			input:   "echo test >",
			wantErr: true,
		},
		{
			name:    "syntax error",
			input:   "echo hello | | grep hello",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			pipeline, err := parseInput(tt.input)
			if err != nil && !tt.wantErr {
				t.Fatalf("parseInput failed: %v", err)
			}

			if pipeline != nil {
				success, execErr := executePipelineWithContext(ctx, pipeline, true)

				if tt.wantErr {
					if execErr == nil && success {
						t.Error("expected error but got success")
					}
				} else {
					if execErr != nil {
						t.Errorf("unexpected error: %v", execErr)
					}
				}
			}
		})
	}
}

func TestBuiltinCommands(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		validate func(stdout, stderr string) bool
	}{
		{
			name:  "kill command with invalid PID",
			input: "kill invalid_pid",
			validate: func(stdout, stderr string) bool {
				// Ошибка kill выводится в stderr
				return strings.Contains(stderr, "неверный PID") ||
					strings.Contains(stderr, "kill:")
			},
		},
		{
			name:  "ps command",
			input: "ps",
			validate: func(stdout, stderr string) bool {
				// ps выводит в stdout
				return strings.Contains(stdout, "PID") ||
					len(stdout) > 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			pipeline, err := parseInput(tt.input)
			if err != nil {
				t.Fatalf("parseInput failed: %v", err)
			}

			// Захватываем stdout и stderr
			var stdoutBuf, stderrBuf bytes.Buffer
			originalStdout := os.Stdout
			originalStderr := os.Stderr

			rOut, wOut, _ := os.Pipe()
			rErr, wErr, _ := os.Pipe()
			os.Stdout = wOut
			os.Stderr = wErr

			success, err := executePipelineWithContext(ctx, pipeline, true)

			// Восстанавливаем и читаем вывод
			wOut.Close()
			wErr.Close()
			os.Stdout = originalStdout
			os.Stderr = originalStderr

			stdoutBuf.ReadFrom(rOut)
			stderrBuf.ReadFrom(rErr)

			stdoutStr := stdoutBuf.String()
			stderrStr := stderrBuf.String()

			if tt.validate != nil && !tt.validate(stdoutStr, stderrStr) {
				t.Errorf("output validation failed:\nstdout: %q\nstderr: %q", stdoutStr, stderrStr)
			}

			_ = success
		})
	}
}

func getCurrentDir() string {
	dir, _ := os.Getwd()
	return dir
}

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}
