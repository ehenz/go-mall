package config

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type GoodsSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type StockSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type SrvConfig struct {
	Name           string         `mapstructure:"name" json:"name"`
	Host           string         `mapstructure:"host" json:"host"`
	Port           int            `mapstructure:"port" json:"port"`
	Tags           []string       `json:"tags"`
	MysqlConfig    MysqlConfig    `mapstructure:"mysql" json:"mysql"`
	ConsulConfig   ConsulConfig   `mapstructure:"consul" json:"consul"`
	GoodsSrvConfig GoodsSrvConfig `json:"goods_srv"`
	StockSrvConfig StockSrvConfig `json:"stock_srv"`
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
