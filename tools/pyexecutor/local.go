package pyexecutor

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/google/uuid"
)

func PythonFuncLocal(ctx context.Context, config *LocalOperatorConfig) func(ctx context.Context, params *PythonParams) (string, error) {
	return func(ctx context.Context, params *PythonParams) (string, error) {

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
			if strings.HasPrefix(err.Error(), "internal error") {
				return err.Error(), nil
			}
			return "", fmt.Errorf("execute error: %w", err)
		}
		return result.Stdout, nil
	}
}

func NewPythonToolLocal(ctx context.Context, config *LocalOperatorConfig) (tool.InvokableTool, error) {
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
