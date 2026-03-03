package local

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cloudwego/eino-ext/components/tool/commandline"
)

type IDFormat string

const (
	IDFormatUUID = IDFormat("uuid")           // UUID 格式，例如：550e8400-e29b-41d4-a716-446655440000
	IDFormatTime = IDFormat("20060102150405") // 时间格式，例如：20230801143000
)

// OperatorConfig 定义了 LocalOperator 的配置选项
type OperatorConfig struct {
	TaskIDFormat IDFormat // 新增：任务ID格式，用于生成任务ID
	// RootDir 是所有文件操作允许的根目录。
	// 如果设置，所有文件路径必须位于此目录下，否则操作会被拒绝。
	// 留空表示不限制（仅推荐在受信任环境中使用）。
	RootDir string
	// 默认工作目录
	WorkDir string
	/*
		/usr/bin/python3 或 C:\Python39\python.exe 等
		Conda:
			Linux/macOS: ~/anaconda3/envs/你的环境名称/bin/python 或 ~/miniconda3/envs/你的环境名称/bin/python。
			Windows: C:\Users\你的用户名\anaconda3\envs\你的环境名称\python.exe
	*/
	ExecutablePath string // 新增：Python 可执行文件路径
	// DefaultFilePerm 是写入文件时的默认权限，默认为 0644。
	DefaultFilePerm os.FileMode
	// Shell 指定执行命令时使用的 shell 可执行文件路径。
	// 如果为空，则根据操作系统自动选择（Windows 为 "cmd.exe"，其他为 "/bin/sh"）。
	Shell string
	// ShellArgs 指定传递给 shell 的参数，默认为 ["-c"]。
	// 注意：命令字符串会被追加为最后一个参数。
	ShellArgs []string
}

// LocalOperator 是 commandline.Operator 的本地实现，支持文件操作和命令执行。
type LocalOperator struct {
	config *OperatorConfig
}

// NewLocalOperator 创建一个新的 LocalOperator，并应用默认配置。
func NewLocalOperator(config *OperatorConfig) *LocalOperator {
	if config.DefaultFilePerm == 0 {
		config.DefaultFilePerm = 0644
	}
	if config.Shell == "" {
		if runtime.GOOS == "windows" {
			config.Shell = "cmd.exe"
			config.ShellArgs = []string{"/C"}
		} else {
			config.Shell = "/bin/sh"
			config.ShellArgs = []string{"-c"}
		}
	} else if config.ShellArgs == nil {
		config.ShellArgs = []string{"-c"} // 如果用户指定了 shell 但未指定参数，使用默认参数
	}
	return &LocalOperator{config: config}
}

// safePath 检查并返回安全的绝对路径，确保其在 RootDir 内（如果配置了 RootDir）。
func (l *LocalOperator) safePath(path string) (string, error) {
	// 清理路径（去除 .. 和 . 等）
	cleanPath := filepath.Clean(path)

	// 如果未设置 RootDir，直接返回绝对路径（避免相对路径歧义）
	if l.config.RootDir == "" {
		abs, err := filepath.Abs(cleanPath)
		if err != nil {
			return "", fmt.Errorf("failed to get absolute path: %w", err)
		}
		return abs, nil
	}

	// 获取 RootDir 的绝对路径
	rootAbs, err := filepath.Abs(l.config.RootDir)
	if err != nil {
		return "", fmt.Errorf("invalid root directory: %w", err)
	}

	// 获取目标路径的绝对路径
	targetAbs, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	// 检查 targetAbs 是否在 rootAbs 内
	// 使用 filepath.Rel 判断相对路径是否以 ".." 开头
	rel, err := filepath.Rel(rootAbs, targetAbs)
	if err != nil {
		return "", fmt.Errorf("path relation error: %w", err)
	}
	if strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return "", fmt.Errorf("path %q is outside of root directory %q", targetAbs, rootAbs)
	}

	return targetAbs, nil
}

/*
实现了 eino 的 commandline.Operator 接口，支持文件操作和命令执行。
*/

// ReadFile 读取文件内容，返回字符串。
func (l *LocalOperator) ReadFile(ctx context.Context, path string) (string, error) {
	safePath, err := l.safePath(path)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(safePath)
	if err != nil {
		return "", fmt.Errorf("read file %q: %w", safePath, err)
	}
	return string(data), nil
}

// WriteFile 将内容写入文件，使用配置的默认权限。
func (l *LocalOperator) WriteFile(ctx context.Context, path string, content string) error {
	safePath, err := l.safePath(path)
	if err != nil {
		return err
	}
	// 确保父目录存在
	dir := filepath.Dir(safePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create parent directories for %q: %w", safePath, err)
	}
	return os.WriteFile(safePath, []byte(content), l.config.DefaultFilePerm)
}

// IsDirectory 判断路径是否为目录。
func (l *LocalOperator) IsDirectory(ctx context.Context, path string) (bool, error) {
	safePath, err := l.safePath(path)
	if err != nil {
		return false, err
	}
	info, err := os.Stat(safePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // 不存在，所以不是目录
		}
		return false, fmt.Errorf("stat %q: %w", safePath, err)
	}
	return info.IsDir(), nil
}

// Exists 判断路径是否存在。
func (l *LocalOperator) Exists(ctx context.Context, path string) (bool, error) {
	safePath, err := l.safePath(path)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(safePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("stat %q: %w", safePath, err)
	}
	return true, nil
}

// RunCommand 在本地 shell 中执行命令，并返回输出。
func (l *LocalOperator) RunCommand(ctx context.Context, command []string) (*commandline.CommandOutput, error) {
	// 获取工作目录（从上下文或配置）
	wd := l.config.WorkDir
	if wd == "" {
		// 如果上下文中未设置，使用 RootDir 作为工作目录（如果 RootDir 非空）
		if l.config.RootDir != "" {
			wd = l.config.RootDir
		} else {
			// 否则使用当前进程的工作目录
			var err error
			wd, err = os.Getwd()
			if err != nil {
				return nil, fmt.Errorf("failed to get current working directory: %w", err)
			}
		}
	}

	// 构建完整的 shell 命令
	// 注意：将用户命令作为单个字符串拼接后传递给 shell
	cmdLine := strings.Join(command, " ")
	shellCmd := append([]string{l.config.Shell}, append(l.config.ShellArgs, cmdLine)...)

	// 创建命令对象
	cmd := exec.CommandContext(ctx, shellCmd[0], shellCmd[1:]...)
	cmd.Dir = wd

	// 捕获输出
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf

	err := cmd.Run()
	if err != nil {
		// 如果是退出码错误，返回标准错误内容以便调试
		return &commandline.CommandOutput{
			Stdout: outBuf.String(),
			Stderr: errBuf.String(),
		}, fmt.Errorf("command execution failed: %w\nstderr: %s", err, errBuf.String())
	}

	return &commandline.CommandOutput{
		Stdout: outBuf.String(),
		Stderr: errBuf.String(),
	}, nil
}
