package cfgs

import (
	"path/filepath"
	"time"

	"github.com/jopbrown/gobase/fsutil"
)

type SummaryRules struct {
	Include     []string `yaml:"Include" json:"include,omitempty"`
	Exclude     []string `yaml:"Exclude" json:"exclude,omitempty"`
	OutDir      Expander `yaml:"OutDir" json:"out_dir,omitempty"`
	OutFileName Expander `yaml:"OutFileName" json:"out_file_name,omitempty"`
}

func (rule *SummaryRules) GetReportDir() string {
	dict := map[string]string{
		"appDir": fsutil.AppDir(),
	}
	return rule.OutDir.ExpandByDict(dict)
}

func (rule *SummaryRules) GetReportPath(codeFolder string) string {
	dict := map[string]string{
		"appDir":             fsutil.AppDir(),
		"timestamp":          time.Now().Format("20060102150405"),
		"codeFolderBaseName": filepath.Base(codeFolder),
	}
	return filepath.Join(rule.OutDir.ExpandByDict(dict), rule.OutFileName.ExpandByDict(dict))
}
