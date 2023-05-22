package cfgs

type Prompt struct {
	System       Expander `yaml:"System"`
	CodeSummary  Expander `yaml:"CodeSummary"`
	SummaryTable Expander `yaml:"SummaryTable"`
	FinalSummary Expander `yaml:"FinalSummary"`
}
