package sandbox

import (
	"context"
	"log"

	"github.com/LouYuanbo1/go-eino-agent/tools/pyexecutor/params"
	"github.com/cloudwego/eino-ext/components/tool/commandline"
	"github.com/cloudwego/eino-ext/components/tool/commandline/sandbox"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

func PythonFuncInSandbox(ctx context.Context, config *sandbox.Config) func(ctx context.Context, params *params.PythonParams) (string, error) {
	return func(ctx context.Context, params *params.PythonParams) (string, error) {
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
