package config

type Redis struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Db       int    `json:"db" yaml:"db"`
	Password string `json:"password" yaml:"password"`
}
