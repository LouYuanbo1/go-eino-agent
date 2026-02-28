package spider

type SpiderConfig struct {
	UserDataDir          string `mapstructure:"user_data_dir"`
	UserMode             bool   `mapstructure:"user_mode"`
	Headless             bool   `mapstructure:"headless"`
	DisableBlinkFeatures string `mapstructure:"disable_blink_features"`
	Incognito            bool   `mapstructure:"incognito"`
	DisableDevShmUsage   bool   `mapstructure:"disable_dev_shm_usage"`
	NoSandbox            bool   `mapstructure:"no_sandbox"`
	DefaultPageWidth     int    `mapstructure:"default_page_width"`
	DefaultPageHeight    int    `mapstructure:"default_page_height"`
	UserAgent            string `mapstructure:"user_agent"`
	Leakless             bool   `mapstructure:"leakless"`
	Bin                  string `mapstructure:"bin"`
}
