package sumer

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/jopbrown/gobase/errors"
	"github.com/jopbrown/gobase/fsutil"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Message = openai.ChatCompletionMessage

type CodeSummary struct {
	PartialList    []*PartialSummaries
	FinalSummaryQA *QA
}

func NewCodeSummary() *CodeSummary {
	cs := &CodeSummary{}
	cs.PartialList = make([]*PartialSummaries, 0, 10)
	return cs
}

func (cs *CodeSummary) SaveMarkdown(fname string) error {
	f, err := fsutil.OpenFileWrite(fname)
	if err != nil {
		return errors.ErrorAt(err)
	}
	defer f.Close()

	err = cs.WriteMarkdown(f)
	if err != nil {
		return errors.ErrorAt(err)
	}
	return nil
}

func (cs *CodeSummary) WriteMarkdown(w io.Writer) error {
	fmt.Fprintln(w, "# Code Summary")
	fmt.Fprintln(w)
	fmt.Fprintln(w, cs.FinalSummary())
	fmt.Fprintln(w)

	slices.SortFunc(cs.PartialList, func(a, b *PartialSummaries) bool {
		return a.FilesSummary[0].FileName < b.FilesSummary[0].FileName
	})

	cs.WriteFileSummaryTable(w)
	fmt.Fprintln(w)

	for _, report := range cs.PartialList {
		for _, fa := range report.FilesSummary {
			fmt.Fprintln(w)
			fmt.Fprintln(w, "##", fa.FileName)
			fmt.Fprintln(w)
			fmt.Fprintln(w, fa.GetSummary())
		}
	}

	return nil
}

func (cs *CodeSummary) WriteFileSummaryTable(w io.Writer) {

	fmt.Fprintln(w, "| File | Description |")
	fmt.Fprintln(w, "| --- | --- |")
	for _, report := range cs.PartialList {
		summaryDict := make(map[string]string)
		summary := report.GetSummary()
		err := json.Unmarshal([]byte(summary), &summaryDict)
		if err != nil {
			fmt.Fprint(w, err.Error())
			fmt.Fprint(w, summary)
		} else {
			keys := maps.Keys(summaryDict)
			slices.Sort(keys)
			for _, key := range keys {
				descript := strings.ReplaceAll(summaryDict[key], "\n", "<br>")
				fmt.Fprintf(w, "| %s | %s |\n", key, descript)
			}
		}
	}
}

func (cs *CodeSummary) GetFileSummaryTable() string {
	sb := &strings.Builder{}
	cs.WriteFileSummaryTable(sb)
	return sb.String()
}

func (cs *CodeSummary) AddPartialSummaries(systemPrompt string) *PartialSummaries {
	ps := &PartialSummaries{}
	ps.SetSystemPrompt(systemPrompt)
	cs.PartialList = append(cs.PartialList, ps)

	return ps
}

func (cs *CodeSummary) RequestFinalSummaryMessages(question string) []Message {
	if len(cs.PartialList) == 0 {
		panic("unable to request final summary without any partial analyses")
	}

	ms := make([]Message, 0, 10)
	ms = append(ms, *cs.PartialList[0].System)
	for _, pa := range cs.PartialList {
		ms = append(ms, *pa.SummaryQA.Question, *pa.SummaryQA.Answer)
	}
	m := Message{
		Role:    openai.ChatMessageRoleUser,
		Content: question,
	}
	ms = append(ms, m)

	return ms
}

func (cs *CodeSummary) SetFinalSummaryQA(question, answer string) *QA {
	cs.FinalSummaryQA = NewQA(question, answer)
	return cs.FinalSummaryQA
}

func (cs *CodeSummary) FinalSummary() string {
	if cs.FinalSummaryQA == nil {
		return ""
	}

	if cs.FinalSummaryQA.Answer == nil {
		return ""
	}

	return cs.FinalSummaryQA.Answer.Content
}

type PartialSummaries struct {
	System       *Message
	FilesSummary []*FileSummary
	SummaryQA    *QA
}

type FileSummary struct {
	FileName string
	QA       *QA
}

type QA struct {
	Question *Message
	Answer   *Message
}

func NewQA(question, answer string) *QA {
	q := &Message{
		Role:    openai.ChatMessageRoleUser,
		Content: question,
	}
	a := &Message{
		Role:    openai.ChatMessageRoleAssistant,
		Content: answer,
	}

	return &QA{
		Question: q,
		Answer:   a,
	}
}

func (ps *PartialSummaries) SetSystemPrompt(prompt string) *Message {
	m := &Message{
		Role:    openai.ChatMessageRoleSystem,
		Content: prompt,
	}
	ps.System = m
	return m
}

func (ps *PartialSummaries) RequestFileSummaryMessages(question string) []Message {
	ms := make([]Message, 0, 10)
	ms = append(ms, *ps.System)
	for _, fa := range ps.FilesSummary {
		ms = append(ms, *fa.QA.Question)
		ms = append(ms, *fa.QA.Answer)
	}
	m := Message{
		Role:    openai.ChatMessageRoleUser,
		Content: question,
	}
	ms = append(ms, m)
	return ms
}

func (ps *PartialSummaries) AddFileSummary(fileName, question, answer string) *FileSummary {
	fa := &FileSummary{}
	fa.FileName = fileName
	fa.QA = NewQA(question, answer)
	ps.FilesSummary = append(ps.FilesSummary, fa)
	return fa
}

func (ps *PartialSummaries) PopFileSummary() *FileSummary {
	length := len(ps.FilesSummary)
	if length == 0 {
		return nil
	}

	var fa *FileSummary
	fa, ps.FilesSummary = ps.FilesSummary[length-1], ps.FilesSummary[:length-1]
	return fa
}

func (ps *FileSummary) GetSummary() string {
	if ps.QA == nil {
		return ""
	}

	if ps.QA.Answer == nil {
		return ""
	}

	return ps.QA.Answer.Content
}

func (pa *PartialSummaries) SetSummaryQA(question, answer string) *QA {
	pa.SummaryQA = NewQA(question, answer)
	return pa.SummaryQA
}

func (pa *PartialSummaries) GetSummary() string {
	if pa.SummaryQA == nil {
		return ""
	}

	if pa.SummaryQA.Answer == nil {
		return ""
	}

	return pa.SummaryQA.Answer.Content
}

func (r *PartialSummaries) FileList() []string {
	files := make([]string, 0, len(r.FilesSummary))
	for _, fa := range r.FilesSummary {
		files = append(files, fa.FileName)
	}

	return files
}
