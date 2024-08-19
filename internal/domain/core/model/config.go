package model

type Config struct {
	Templates Templates `yaml:"templates"`
}

type Templates struct {
	Service        string `yaml:"service"`
	Application    string `yaml:"application"`
	Infrastructure string `yaml:"infrastructure"`
}
