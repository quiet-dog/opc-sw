package config

type Config struct {
	Sqlite Sqlite `json:"sqlite" yaml:"sqlite"`
	Server Server `json:"server" yaml:"server"`
	Sw     Sw     `json:"sw" yaml:"sw"`
}
