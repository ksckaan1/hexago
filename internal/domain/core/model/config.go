package model

type Config struct {
	Templates Templates          `yaml:"templates"`
	Runners   map[string]*Runner `yaml:"runners"`
}

type Templates struct {
	Service        string `yaml:"service"`
	Application    string `yaml:"application"`
	Infrastructure string `yaml:"infrastructure"`
	Package        string `yaml:"package"`
}

type Runner struct {
	Cmd              string   `yaml:"cmd"`
	EnvVars          []string `yaml:"env"`
	SeperateLogFiles bool     `yaml:"seperate_log_files"`
	DisableLogFiles  bool     `yaml:"disable_log_files"`
}
