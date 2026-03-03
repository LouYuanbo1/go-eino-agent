package local

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/LouYuanbo1/go-eino-agent/tools/pyexecutor/params"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/google/uuid"
)

func PythonFuncLocal(ctx context.Context, config *OperatorConfig) func(ctx context.Context, params *params.PythonParams) (string, error) {
	return func(ctx context.Context, params *params.PythonParams) (string, error) {

		fmt.Printf("调用Python执行工具\n")
		fmt.Printf("python code: %s\n", params.Code)

		op := NewLocalOperator(config)
		var taskID string
		switch IDFormat(config.TaskIDFormat) {
		case IDFormatUUID:
			taskUUID, err := uuid.NewUUID()
			if err != nil {
				return "", fmt.Errorf("generate uuid: %w", err)
			}
			taskID = taskUUID.String()
		case IDFormatTime:
			taskID = time.Now().Format("20060102150405")
		default:
			taskID = time.Now().Format("20060102150405")
		}
		wd := config.WorkDir
		if wd == "" {
			return "", fmt.Errorf("work dir not found")
		}
		filePath := filepath.Join(wd, fmt.Sprintf("%s.py", taskID))
		if err := op.WriteFile(ctx, filePath, params.Code); err != nil {
			return fmt.Sprintf("failed to create python file %s: %v", filePath, err), nil
		}

		fmt.Printf("python file path: %s\n", filePath)

		pyExecutablePath := config.ExecutablePath
		if pyExecutablePath == "" {
			pyExecutablePath = "python"
		}
		result, err := op.RunCommand(ctx, []string{pyExecutablePath, filePath})
		if err != nil {
			if strings.Contains(err.Error(), "executable file not found") {
				return "", fmt.Errorf("python interpreter not found: %w", err)
			}
			// 如果 err 是 *exec.ExitError，说明命令运行了但返回非0退出码
			// 此时应该将 stdout 和 stderr 合并返回，而不是返回错误
			var execError *exec.ExitError
			if ok := errors.As(err, &execError); ok {
				// 安全地从 result 中获取 stdout（如果 result 非 nil）
				var output string
				if result != nil {
					output = result.Stdout
				}

				// 优先使用 execError.Stderr（它包含了命令的标准错误输出）
				if len(execError.Stderr) > 0 {
					if output != "" {
						output += "\n"
					}
					output += "STDERR:\n" + string(execError.Stderr)
				} else if result != nil && result.Stderr != "" {
					// 如果 execError 没有 Stderr，但 result 中有，则使用 result 的（兜底）
					if output != "" {
						output += "\n"
					}
					output += "STDERR:\n" + result.Stderr
				}

				output = fmt.Sprintf("Exit code: %d\n%s", execError.ExitCode(), output)
				return output, nil
			}
			// 其他未知错误
			return "", fmt.Errorf("execute command error: %w", err)
		}
		output := result.Stdout
		if result.Stderr != "" {
			if output != "" {
				output += "\n"
			}
			output += "STDERR:\n" + result.Stderr
		}
		if result.ExitCode != 0 {
			output = fmt.Sprintf("Exit code: %d\n%s", result.ExitCode, output)
		}
		return output, nil
	}
}

func NewPythonToolLocal(ctx context.Context, config *OperatorConfig) (tool.InvokableTool, error) {
	pythonTool, err := utils.InferTool(
		"pythonExecutor", // tool name
		`Python executor are used to execute Python code; they can be used to execute dynamic Python code, 
		such as print("hello world")". 
		It writes the Python code to a file, executes the file, and returns "stdout" and "stderr".`, // tool description`
		PythonFuncLocal(ctx, config), // tool function
	)
	if err != nil {
		return nil, err
	}
	return pythonTool, nil
}
