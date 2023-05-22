package app

import (
	"time"

	"github.com/jopbrown/codesum/pkg/cfgs"
	"github.com/jopbrown/codesum/pkg/utils"
	"github.com/jopbrown/gobase/errors"
	"github.com/jopbrown/gobase/fsutil"
	"github.com/jopbrown/gobase/log"
	"github.com/jopbrown/gobase/log/rotate"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) SaveConfig(cfg *cfgs.Config) error {
	err := cfg.SaveConfig(a.cfgPath)
	if err != nil {
		return wrapGoError(err)
	}
	return nil
}

func applyLog(cfg *cfgs.Config) error {
	f, err := rotate.OpenFile(cfg.LogPath.ExpandByDict(map[string]string{"appDir": fsutil.AppDir()}), 24*time.Hour, 0)
	if err != nil {
		return errors.ErrorAt(err)
	}

	tee := log.NewTeeLogger(
		log.ConsoleLogger(cfg.DebugMode),
		log.FileLogger(f, log.FileLoggerFormat(), true),
	)
	log.SetGlobalLogger(tee)

	return nil
}

func (a *App) ApplyConfig(cfg *cfgs.Config) error {
	a.cfg = cfg

	accessToken := a.cfg.ChatGpt.AccessToken.String()
	if len(accessToken) > 0 {
		err := utils.UpdateApiServerAccessToken(a.cfg.ChatGpt.EndPoint.String(), accessToken)
		if err != nil {
			return wrapGoError(err)
		}
	}

	err := applyLog(cfg)
	if err != nil {
		return errors.ErrorAt(err)
	}
	return nil
}

func (a *App) ApplyAndSaveConfig(cfg *cfgs.Config) error {
	err := a.ApplyConfig(cfg)
	if err != nil {
		return wrapGoError(err)
	}

	err = a.SaveConfig(cfg)
	if err != nil {
		return wrapGoError(err)
	}
	return nil
}

func (a *App) LoadConfig() (*cfgs.Config, error) {
	if !fsutil.ExistsFile(a.cfgPath) {
		runtime.LogErrorf(a.ctx, "config not exist, load default only: %s", a.cfgPath)
		return cfgs.DefaultConfig(), nil
	}
	cfg, err := cfgs.LoadConfig(a.cfgPath)
	if err != nil {
		return nil, wrapGoError(err)
	}
	err = cfg.MergeDefault()
	if err != nil {
		return nil, wrapGoError(err)
	}

	return cfg, nil
}

func (a *App) LoadAndApplyConfig() (*cfgs.Config, error) {
	cfg, err := a.LoadConfig()
	if err != nil {
		return nil, wrapGoError(err)
	}

	err = a.ApplyConfig(cfg)
	if err != nil {
		return nil, wrapGoError(err)
	}

	return cfg, nil
}

func (a *App) GetConfig() *cfgs.Config {
	return a.cfg
}
