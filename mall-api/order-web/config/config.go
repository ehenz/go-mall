package config

type SrvConfig struct {
	Name           string         `mapstructure:"name" json:"name"`
	Host           string         `mapstructure:"host" json:"host"`
	Port           int            `mapstructure:"port" json:"port"`
	Tags           []string       `mapstructure:"tags" json:"tags"`
	GoodsSrvConfig GoodsSrvConfig `mapstructure:"goods-srv" json:"goods-srv"`
	OrderSrvConfig OrderSrvConfig `mapstructure:"order-srv" json:"order-srv"`
	StockSrvConfig StockSrvConfig `json:"stock-srv"`
	JWTConfig      JWTConfig      `mapstructure:"jwt" json:"jwt"`
	ConsulConfig   ConsulConfig   `mapstructure:"consul" json:"consul"`
	Alipay         Alipay         `json:"alipay"`
}

type GoodsSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type OrderSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type StockSrvConfig struct {
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

type Alipay struct {
	Appid        string `json:"appid"`
	PrivateKey   string `json:"private_key"`
	AliPublicKey string `json:"ali_public_key"`
	NotifyUrl    string `json:"notify_url"`
	ReturnUrl    string `json:"return_url"`
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
