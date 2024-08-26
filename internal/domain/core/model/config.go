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
	Cmd     string   `yaml:"cmd"`
	EnvVars []string `yaml:"env"`
	Log     Log      `yaml:"log"`
}

type Log struct {
	Disabled      bool `yaml:"disabled"`
	SeperateFiles bool `yaml:"seperate_files"`
	Overwrite     bool `yaml:"overwrite"`
}
