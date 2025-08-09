package config

type DatabaseConfig struct {
	DbType       string `mapstructure:"db_type" json:"db_type"`
	Host         string `mapstructure:"host" json:"host"`
	Port         int    `mapstructure:"port" json:"port"`
	Name         string `mapstructure:"db" json:"db"`
	User         string `mapstructure:"user" json:"user"`
	Password     string `mapstructure:"password" json:"password"`
	TablePrefix  string `mapstructure:"table_prefix" json:"table_prefix"`
	Charset      string `mapstructure:"charset" json:"charset"`
	ParseTime    bool   `mapstructure:"parse_time" json:"parse_time"`
	MaxIdleConns int    `mapstructure:"max_idle_conns" json:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns" json:"max_open_conns"`
}
