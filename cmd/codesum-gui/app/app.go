package app

import (
	"context"
	"os"
	"path/filepath"

	"github.com/jopbrown/codesum/pkg/cfgs"
	"github.com/jopbrown/codesum/pkg/sumer"
	"github.com/jopbrown/gobase/errors"
	"github.com/jopbrown/gobase/log"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx     context.Context
	cfgPath string
	cfg     *cfgs.Config
	cancel  func()
}

// NewApp creates a new App application struct
func NewApp(cfgPath string) *App {
	a := &App{}
	a.cfgPath = cfgPath
	return a
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Stop() {
	if a.cancel != nil {
		a.cancel()
		a.cancel = nil
	}
}

func (a *App) SelectFolder() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
}

func wrapGoError(err error) error {
	log.ErrorAt(err)
	return errors.RootCause(err)
}

func (a *App) CodeSummarize(codeFolder string) (string, error) {
	if a.cfg == nil {
		return "", errors.Error(`the config not load yet`)
	}

	sumer, err := sumer.NewSummarizer(a.cfg, a.PushMessage, a.StreamAnswer)
	if err != nil {
		return "", wrapGoError(err)
	}

	ctx, cancel := context.WithCancel(a.ctx)
	a.cancel = cancel
	reportPath, err := sumer.Summarize(ctx, codeFolder)
	if err != nil {
		return "", wrapGoError(err)
	}

	return reportPath, nil
}

func (a *App) GetReportList() ([]string, error) {
	reportList, err := filepath.Glob(filepath.Join(a.cfg.SummaryRules.GetReportDir(), "*.md"))
	if err != nil {
		return nil, wrapGoError(err)
	}

	return reportList, nil
}

func (a *App) GetFileContent(fname string) (string, error) {
	data, err := os.ReadFile(fname)
	if err != nil {
		return "", wrapGoError(err)
	}

	return string(data), nil
}
