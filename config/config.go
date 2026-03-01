package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Elasticsearch struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"elasticsearch"`
}

func InitConfig() (*Config, error) {
	// 配置文件设置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	path0 := filepath.Join(".")
	path1 := filepath.Join("config")
	path2 := filepath.Join("..", "config")
	path3 := filepath.Join("..", "..", "config")
	path4 := filepath.Join("..", "..", "..", "config")
	//Viper查询的路径是相对于当前工作目录（current working directory） 的，而不是相对于可执行文件的位置或源代码文件的位置。
	// 1.尝试加载main.go所在目录的上两级目录下config.yaml
	for _, path := range []string{path0, path1, path2, path3, path4} {
		viper.AddConfigPath(path)
	}
	// 2.尝试加载当前工作目录下的config.yaml
	viper.AddConfigPath(".")

	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)

	// 设置绝对路径
	// 3.尝试加载main.go可执行文件所在目录的上两级目录下的config.yaml
	pathExe0 := filepath.Join(exeDir)
	pathExe1 := filepath.Join(exeDir, "config")
	pathExe2 := filepath.Join(exeDir, "..", "config")
	pathExe3 := filepath.Join(exeDir, "..", "..", "config")
	// 4.尝试加载main.go可执行文件所在目录的上三级目录下的config.yaml
	pathExe4 := filepath.Join(exeDir, "..", "..", "..", "config")
	for _, path := range []string{pathExe0, pathExe1, pathExe2, pathExe3, pathExe4} {
		viper.AddConfigPath(path)
	}

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found, using defaults")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// 绑定环境变量
	viper.AutomaticEnv()

	// 监控配置变化
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()

	cfg, err := parseConfig()
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
	return cfg, nil
}

func parseConfig() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}
	return &cfg, nil
}
