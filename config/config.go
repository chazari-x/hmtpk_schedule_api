package config

type Server struct {
	HTTP HTTP `yaml:"http"`
	GRPC GRPC `yaml:"grpc"`
}

type HTTP struct {
	Domain       string `yaml:"domain"`
	HTTPAddress  string `yaml:"http_address"`
	HTTPSAddress string `yaml:"https_address"`
	MiniAppToken string `yaml:"mini_app_token"`
}

type GRPC struct {
	Domain       string `yaml:"domain"`
	Address      string `yaml:"address"`
	MiniAppToken string `yaml:"mini_app_token"`
}

type Redis struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Pass string `yaml:"password"`
}
