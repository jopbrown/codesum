package sumer

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jopbrown/codesum/pkg/cfgs"
	"github.com/jopbrown/codesum/pkg/utils"
	"github.com/jopbrown/codesum/pkg/utils/try"

	"github.com/jopbrown/gobase/errors"
	"github.com/jopbrown/gobase/fsutil"
	"github.com/jopbrown/gobase/log"
	"github.com/jopbrown/gobase/strutil"
	"github.com/sashabaranov/go-openai"
)

type Summarizer struct {
	cfg                  *cfgs.Config
	gptClient            *openai.Client
	fileFilter           strutil.Matcher
	pushMsgCallback      func(msg *openai.ChatCompletionMessage)
	streamAnswerCallback func(delta string)
}

func NewSummarizer(cfg *cfgs.Config, pushMsgCallback func(msg *openai.ChatCompletionMessage), streamAnswerCallbackOpt ...func(token string)) (*Summarizer, error) {
	sumer := &Summarizer{}
	sumer.cfg = cfg

	gptCfg := openai.DefaultConfig(cfg.ChatGpt.APIKey.String())
	gptCfg.BaseURL = cfg.ChatGpt.EndPoint.String()
	if proxy := cfg.ChatGpt.Proxy.String(); proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return nil, errors.ErrorAt(err)
		}
		gptCfg.HTTPClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}
	sumer.gptClient = openai.NewClientWithConfig(gptCfg)

	includer, err := strutil.CompileGitIgnoreLines(cfg.SummaryRules.Include...)
	if err != nil {
		return nil, errors.ErrorAt(err)
	}

	excluder, err := strutil.CompileGitIgnoreLines(cfg.SummaryRules.Exclude...)
	if err != nil {
		return nil, errors.ErrorAt(err)
	}

	sumer.fileFilter = strutil.MakeMultiMatcher(includer.Any(), excluder.None()).All()
	sumer.pushMsgCallback = pushMsgCallback
	if len(streamAnswerCallbackOpt) > 0 {
		sumer.streamAnswerCallback = streamAnswerCallbackOpt[0]
	}

	return sumer, nil
}

func (sumer *Summarizer) Summarize(ctx context.Context, codeFolder string) (reportPath string, err error) {
	codeFiles, err := fsutil.ListWithMatcher(codeFolder, sumer.fileFilter)
	if err != nil {
		return "", errors.ErrorAt(err)
	}

	cs := &CodeSummary{}
	defer func() {
		// save reports regardless of process success or failure
		reportPath, _ = sumer.saveMarkdownReport(codeFolder, cs)
	}()

	ps := cs.AddPartialSummaries(sumer.cfg.Prompt.System.String())
	sumer.pushMsgCallback(ps.System)

	for _, fname := range codeFiles {
		if isCtxCanceled(ctx) {
			break
		}

		fileName, content, err := getCodeFileContent(codeFolder, fname)
		if err != nil {
			return "", errors.ErrorAt(err)
		}

		question := FileSummaryQuestion(sumer.cfg.Prompt.CodeSummary, fileName, content)
		msgs := ps.RequestFileSummaryMessages(question)

		answer, err := sumer.sendRequest(ctx, msgs)

		if err != nil {
			if utils.IsErrMatchHTTPCode(err, 413) && len(msgs) > 3 {
				// out of max token, roll to next partial
				ps = cs.AddPartialSummaries(sumer.cfg.Prompt.System.String())
				msgs = ps.RequestFileSummaryMessages(question)
				answer, err = sumer.sendRequest(ctx, msgs)
				if err != nil {
					return "", errors.ErrorAt(err)
				}
			} else {
				return "", errors.ErrorAt(err)
			}
		}

		ps.AddFileSummary(fileName, question, answer)
	}

	// for _, ps := range cs.PartialList {
	for i := 0; i < len(cs.PartialList); /*i++*/ {
		ps = cs.PartialList[i]
		question := SummaryTableQuestion(sumer.cfg.Prompt.SummaryTable, ps.FileList())
		msgs := ps.RequestFileSummaryMessages(question)

		answer, err := sumer.sendRequest(ctx, msgs)

		if err != nil {
			if utils.IsErrMatchHTTPCode(err, 413) && len(msgs) > 3 {
				// out of max token, pop the latest file
				fs := ps.PopFileSummary()
				question = SummaryTableQuestion(sumer.cfg.Prompt.SummaryTable, ps.FileList())
				msgs = ps.RequestFileSummaryMessages(question)
				answer, err = sumer.sendRequest(ctx, msgs)
				if err != nil {
					return "", errors.ErrorAt(err)
				}

				// add the poped file to addition partial
				newPs := cs.AddPartialSummaries(sumer.cfg.Prompt.System.String())
				newPs.AddFileSummary(fs.FileName, fs.QA.Question.Content, fs.QA.Answer.Content)
			} else {
				return "", errors.ErrorAt(err)
			}
		}

		answer = trimJSONAnswer(answer)

		ps.SetSummaryQA(question, answer)

		// check answer is valid JSON, or try again(not move to next iter)
		if json.Valid([]byte(answer)) {
			i++
		}
	}

	{
		question := FinalSummaryQuestion(sumer.cfg.Prompt.FinalSummary)
		msgs := cs.RequestFinalSummaryMessages(question)
		answer, err := sumer.sendRequest(ctx, msgs)

		if err != nil {
			return "", errors.ErrorAt(err)
		}
		cs.SetFinalSummaryQA(question, answer)
	}

	if err != nil {
		return "", errors.ErrorAt(err)
	}

	return
}

func trimJSONAnswer(answer string) string {
	answer = strings.TrimSpace(answer)
	answer = strings.TrimPrefix(answer, "```json")
	answer = strings.TrimPrefix(answer, "```")
	answer = strings.TrimSuffix(answer, "```")
	return answer
}

func (sumer *Summarizer) sendRequest(ctx context.Context, ms []openai.ChatCompletionMessage) (string, error) {
	debugW := log.GetWriter(log.LevelDebug)

	sumer.pushMsgCallback(&ms[len(ms)-1]) // push question

	var stream *openai.ChatCompletionStream
	err := try.Do(func() error {
		var err error
		stream, err = sumer.gptClient.CreateChatCompletionStream(
			ctx,
			openai.ChatCompletionRequest{
				Model:    sumer.cfg.ChatGpt.Model.String(),
				Messages: ms,
			},
		)
		if err != nil {
			return errors.ErrorAt(err)
		}
		return nil
	}, try.Option().SetLimitTimes(3).SetInterval(5*time.Second).SetOnFail(func(err error) bool {
		return utils.IsErrMatchHTTPCode(err, 413) // out of max token, stop retry and return error
	}))

	if err != nil {
		return "", errors.ErrorAt(err)
	}
	defer stream.Close()

	sb := &strings.Builder{}

	log.Debugf("start recive response ...")
	for {
		if isCtxCanceled(ctx) {
			break
		}

		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return "", errors.ErrorAt(err)
		}

		delta := resp.Choices[0].Delta.Content
		sb.WriteString(delta)
		io.WriteString(debugW, delta)
		if sumer.streamAnswerCallback != nil {
			sumer.streamAnswerCallback(delta)
		}
	}

	answer := sb.String()

	// push answer
	sumer.pushMsgCallback(&Message{
		Role:    openai.ChatMessageRoleAssistant,
		Content: answer,
	})

	return answer, nil
}

func getCodeFileContent(codeFolder, fname string) (name string, content string, err error) {
	fdata, err := os.ReadFile(fname)
	if err != nil {
		return "", "", errors.ErrorAt(err)
	}

	content = string(fdata)
	relPath, err := filepath.Rel(codeFolder, fname)
	if err != nil {
		return "", "", errors.ErrorAt(err)
	}
	name = filepath.ToSlash(relPath)
	return
}

func (sumer *Summarizer) saveMarkdownReport(codeFolder string, cs *CodeSummary) (string, error) {
	reportPath := sumer.cfg.SummaryRules.GetReportPath(codeFolder)
	err := cs.SaveMarkdown(reportPath)
	if err != nil {
		return "", errors.ErrorAt(err)
	}

	return reportPath, nil
}

func isCtxCanceled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
