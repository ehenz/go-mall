package config

type SrvConfig struct {
	Name          string        `mapstructure:"name" json:"name"`
	Host          string        `mapstructure:"host" json:"host"`
	Port          int           `mapstructure:"port" json:"port"`
	Tags          []string      `mapstructure:"tags" json:"tags"`
	UserSrvConfig UserSrvConfig `mapstructure:"user-srv" json:"user-srv"`
	JWTConfig     JWTConfig     `mapstructure:"jwt" json:"jwt"`
	RedisConfig   RedisConfig   `mapstructure:"redis" json:"redis"`
	ConsulConfig  ConsulConfig  `mapstructure:"consul" json:"consul"`
}

type UserSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type RedisConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"data-id"`
	Group     string `mapstructure:"group"`
}
