package config

type AppConfig struct {
	RunMode     string `mapstructure:"run_mode" json:"run_mode"`
	LogSavePath string `mapstructure:"log_save_path" json:"host"`
	LogFileExt  string `mapstructure:"log_file_ext" json:"port"`
}

type ServerConfig struct {
	AppInfo      AppConfig      `mapstructure:"app" json:"app"`
	ErpMysqlInfo DatabaseConfig `mapstructure:"erp_mysql" json:"erp_mysql"`
}

type ServiceConfig struct {
	Url       string `mapstructure:"url" json:"url" yaml:"url"`
	SecretKey string `mapstructure:"secret_key" json:"secret_key" yaml:"secret"`
}

// cmd 启动脚本配置结构体
type CronServerConfig struct {
}
