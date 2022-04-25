package config

type SrvConfig struct {
	Name            string          `mapstructure:"name" json:"name"`
	Host            string          `mapstructure:"host" json:"host"`
	Port            int             `mapstructure:"port" json:"port"`
	Tags            []string        `mapstructure:"tags" json:"tags"`
	UseropSrvConfig UseropSrvConfig `json:"userop-srv"`
	GoodsSrvConfig  GoodsSrvConfig  `json:"goods-srv"`
	JWTConfig       JWTConfig       `mapstructure:"jwt" json:"jwt"`
	ConsulConfig    ConsulConfig    `mapstructure:"consul" json:"consul"`
}

type UseropSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type GoodsSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
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
