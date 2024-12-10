package config

type Sw struct {
	Host  string `json:"host" yaml:"host"`
	Port  string `json:"port" yaml:"port"`
	Topic string `json:"topic" yaml:"topic"`
}
