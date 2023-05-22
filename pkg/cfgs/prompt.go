package cfgs

type Prompt struct {
	System       Expander `yaml:"System" json:"system,omitempty"`
	CodeSummary  Expander `yaml:"CodeSummary" json:"code_summary,omitempty"`
	SummaryTable Expander `yaml:"SummaryTable" json:"summary_table,omitempty"`
	FinalSummary Expander `yaml:"FinalSummary" json:"final_summary,omitempty"`
}
