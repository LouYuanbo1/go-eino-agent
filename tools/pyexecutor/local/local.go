package local

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/LouYuanbo1/go-eino-agent/tools/pyexecutor/params"
	"github.com/cloudwego/eino-ext/components/tool/commandline"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/google/uuid"
)

func PythonFuncLocal(ctx context.Context, config *OperatorConfig) func(ctx context.Context, params *params.PythonParams) (string, error) {
	return func(ctx context.Context, params *params.PythonParams) (string, error) {
		op := NewLocalOperator(config)

		// 生成唯一 ID
		var fileID string
		switch IDFormat(config.TaskIDFormat) {
		case IDFormatUUID:
			taskUUID, err := uuid.NewUUID()
			if err != nil {
				return "", fmt.Errorf("generate uuid: %w", err)
			}
			fileID = taskUUID.String()
		default: // 默认使用时间戳
			fileID = time.Now().Format("20060102150405")
		}

		wd := config.WorkDir
		if wd == "" {
			return "", fmt.Errorf("work dir not found")
		}

		// 临时 Python 文件路径
		tempFileName := fmt.Sprintf("%s_temp_%s.py", config.FileName, fileID)
		tempFilePath := filepath.Join(wd, tempFileName)

		// 写入代码
		if err := op.WriteFile(ctx, tempFilePath, params.Code); err != nil {
			return "", fmt.Errorf("failed to create python file %s: %w", tempFilePath, err)
		}
		defer os.Remove(tempFilePath) // 确保最终删除

		// 准备执行
		pyExecutablePath := config.ExecutablePath
		if pyExecutablePath == "" {
			pyExecutablePath = "python"
		}

		result, err := op.RunCommand(ctx, []string{pyExecutablePath, tempFilePath})
		if err != nil {
			// 处理解释器未找到
			if strings.Contains(err.Error(), "executable file not found") {
				return "", fmt.Errorf("python interpreter not found: %w", err)
			}
			// 处理命令执行但返回非零退出码
			var execError *exec.ExitError
			if ok := errors.As(err, &execError); ok {
				output := buildOutput(result, execError)
				return output, nil
			}
			return "", fmt.Errorf("execute command error: %w", err)
		}

		// 正常执行完成
		output := buildOutput(result, nil)

		// 若执行成功（stderr 均为空），则保存源码为 .py 文件
		if result.ExitCode == 0 {
			pyFileName := fmt.Sprintf("%s_%s.py", config.FileName, fileID)
			pyFilePath := filepath.Join(wd, pyFileName)

			// 读取临时文件内容（使用 op 确保安全）
			content, err := op.ReadFile(ctx, tempFilePath)
			if err != nil {
				return "", fmt.Errorf("read temp file: %w", err)
			}
			// 写入目标文件
			if err := op.WriteFile(ctx, pyFilePath, content); err != nil {
				return "", fmt.Errorf("write target file: %w", err)
			}
		}

		return output, nil
	}
}

// buildOutput 辅助函数，用于组装最终输出字符串
func buildOutput(result *commandline.CommandOutput, execErr *exec.ExitError) string {
	var output string
	if result != nil {
		output = result.Stdout
	}
	if execErr != nil && len(execErr.Stderr) > 0 {
		if output != "" {
			output += "\n"
		}
		output += "STDERR:\n" + string(execErr.Stderr)
	} else if result != nil && result.Stderr != "" {
		if output != "" {
			output += "\n"
		}
		output += "STDERR:\n" + result.Stderr
	}
	exitCode := 0
	if execErr != nil {
		exitCode = execErr.ExitCode()
	} else if result != nil {
		exitCode = result.ExitCode
	}
	if exitCode != 0 {
		output = fmt.Sprintf("Exit code: %d\n%s", exitCode, output)
	}

	return output
}

func NewPythonToolLocal(ctx context.Context, config *OperatorConfig) (tool.InvokableTool, error) {
	description :=
		`
	 Python executor can run any Python code you provide. 
	 It is ideal for dynamic calculations, data processing, web scraping, file manipulation, or testing code snippets. 
	 The tool writes your code to a temporary file, executes it, and returns the standard output (stdout) and standard error (stderr). 
	 Use it whenever you need to compute something, generate results, or verify Python behavior. 
	 Examples: print("hello world"), 2+2, import requests; response = requests.get("https://api.example.com"), or any custom script.
	`
	pythonTool, err := utils.InferTool(
		"pythonExecutor",             // tool name
		description,                  // tool description
		PythonFuncLocal(ctx, config), // tool function
	)
	if err != nil {
		return nil, err
	}
	return pythonTool, nil
}
