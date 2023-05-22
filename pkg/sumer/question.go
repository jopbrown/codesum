package sumer

import (
	"strings"

	"github.com/jopbrown/gobase/strutil"
)

func FileSummaryQuestion(prompt strutil.Expander, fileName, fileContent string) string {
	return prompt.ExpandByDict(map[string]string{
		"fileName":    fileName,
		"fileContent": fileContent,
	})
}

func SummaryTableQuestion(prompt strutil.Expander, fileList []string) string {
	return prompt.ExpandByDict(map[string]string{
		"filesCommaList": strings.Join(fileList, ","),
	})
}

func FinalSummaryQuestion(prompt strutil.Expander) string {
	return prompt.String()
}
