package cfgs

type SummaryRules struct {
	Include     []string `yaml:"Include"`
	Exclude     []string `yaml:"Exclude"`
	OutDir      Expander `yaml:"OutDir"`
	OutFileName Expander `yaml:"OutFileName"`
}
