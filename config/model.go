package config

type store struct {
	Runners   map[string]*Runner `yaml:"runners"`
	Templates templates          `yaml:"templates"`
}

type templates struct {
	Service        string `yaml:"service"`
	Application    string `yaml:"application"`
	Infrastructure string `yaml:"infrastructure"`
	Package        string `yaml:"package"`
}

type Runner struct {
	Cmd     string   `yaml:"cmd"`
	EnvVars []string `yaml:"env"`
	Log     log      `yaml:"log"`
}

type log struct {
	Disabled      bool `yaml:"disabled"`
	SeperateFiles bool `yaml:"seperate_files"`
	Overwrite     bool `yaml:"overwrite"`
}
