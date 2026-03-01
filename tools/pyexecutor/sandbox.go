package pyexecutor

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/tool/commandline"
	"github.com/cloudwego/eino-ext/components/tool/commandline/sandbox"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type PythonParams struct {
	Code string `json:"code" jsonschema:"description=用于指定要执行的 Python 代码"`
}

func PythonFuncInSandbox(ctx context.Context, config *sandbox.Config) func(ctx context.Context, params *PythonParams) (string, error) {
	return func(ctx context.Context, params *PythonParams) (string, error) {
		op, err := sandbox.NewDockerSandbox(ctx, config)
		if err != nil {
			log.Fatal(err)
		}
		// you should ensure that docker has been started before create a docker container
		err = op.Create(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer op.Cleanup(ctx)

		exec, err := commandline.NewPyExecutor(ctx, &commandline.PyExecutorConfig{Operator: op}) // use python3 by default
		if err != nil {
			log.Fatal(err)
		}
		result, err := exec.Execute(ctx, &commandline.Input{Code: params.Code})
		if err != nil {
			log.Fatal(err)
		}
		return result.Stdout, nil
	}
}

func NewPythonToolInSandbox(ctx context.Context, config *sandbox.Config) (tool.InvokableTool, error) {
	pythonTool, err := utils.InferTool(
		"pythonExecutor", // tool name
		`Python executor are used to execute Python code; 
		they can be used to execute dynamic Python code, eg: print("hello world")`, // tool description
		PythonFuncInSandbox(ctx, config), // tool function
	)
	if err != nil {
		return nil, err
	}
	return pythonTool, nil
}
