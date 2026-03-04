package pythonAgent

import (
	"fmt"
	"os/exec"
)

// GetPipPackages 通过 pip 获取当前 Python 环境的包列表
func GetPackages(executablePath string) (string, error) {
	cmd := exec.Command(executablePath, "-m", "pip", "list", "--format=json")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("执行 pip list 失败: %w", err)
	}
	return string(output), nil
}
