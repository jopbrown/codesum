package cfgs

import (
	"embed"
	"io"
	"os"

	"github.com/imdario/mergo"
	"github.com/jopbrown/gobase/errors"
	"github.com/jopbrown/gobase/fsutil"
	"github.com/jopbrown/gobase/strutil"
	"gopkg.in/yaml.v3"
)

type Expander = strutil.Expander

type Config struct {
	DebugMode    bool          `yaml:"DebugMode" json:"debug_mode,omitempty"`
	LogPath      Expander      `yaml:"LogPath" json:"log_path,omitempty"`
	ChatGpt      *ChatGpt      `yaml:"ChatGpt" json:"chat_gpt,omitempty"`
	SummaryRules *SummaryRules `yaml:"SummaryRules" json:"summary_rules,omitempty"`
	Prompt       *Prompt       `yaml:"Prompt" json:"prompt,omitempty"`
}

//go:embed default
var defaultCfgFs embed.FS

func DefaultConfig() *Config {
	r := errors.Must1(defaultCfgFs.Open("default/default.yml"))
	cfg := errors.Must1(ReadConfig(r))
	return cfg
}

func LoadConfig(fname string) (*Config, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, errors.ErrorAt(err)
	}
	defer f.Close()

	cfg, err := ReadConfig(f)
	if err != nil {
		return nil, errors.ErrorAt(err)
	}

	return cfg, nil
}

func ReadConfig(r io.Reader) (*Config, error) {
	cfg := &Config{}
	err := yaml.NewDecoder(r).Decode(cfg)
	if err != nil {
		return nil, errors.ErrorAt(err)
	}

	return cfg, nil
}

func (cfg *Config) MergeDefault() error {
	err := mergo.Merge(cfg, DefaultConfig())

	if err != nil {
		return errors.ErrorAt(err)
	}

	return nil
}

func (cfg *Config) SaveConfig(fname string) error {
	f, err := fsutil.OpenFileWrite(fname)
	if err != nil {
		return errors.ErrorAt(err)
	}
	defer f.Close()

	err = cfg.WriteConfig(f)
	if err != nil {
		return errors.ErrorAt(err)
	}

	return nil
}

func (cfg *Config) WriteConfig(w io.Writer) error {
	err := yaml.NewEncoder(w).Encode(cfg)
	if err != nil {
		return errors.ErrorAt(err)
	}

	return nil
}
