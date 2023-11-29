package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/jopbrown/codesum/pkg/cfgs"
	"github.com/jopbrown/codesum/pkg/sumer"
	"github.com/jopbrown/codesum/pkg/utils"

	"github.com/jopbrown/gobase/errors"
	"github.com/jopbrown/gobase/fsutil"
	"github.com/jopbrown/gobase/log"
	"github.com/jopbrown/gobase/log/rotate"
	"golang.org/x/sync/errgroup"
)

var (
	BuildName    = "myapp"
	BuildVersion = "v0.0.0"
	BuildHash    = "unknown"
	BuildTime    = "20060102150405"
)

var args struct {
	configPath string
	codeFolder string
}

func init() {
	flag.StringVar(&args.configPath, "c", filepath.Join(fsutil.AppDir(), "config.yml"), "config path")
	parseArgs()
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(errors.GetErrorDetails(err))
	}
}

func run() error {
	cfg, err := cfgs.LoadConfig(args.configPath)
	if err != nil {
		return errors.ErrorAt(err)
	}
	cfg.MergeDefault()

	err = applyLog(cfg)
	if err != nil {
		return errors.ErrorAt(err)
	}
	cfg.WriteConfig(log.GetWriter(log.LevelDebug))
	log.Infof("%s %v-%v-%v", BuildName, BuildVersion, BuildHash, BuildTime)

	accessToken := cfg.ChatGpt.AccessToken.String()
	if strings.HasPrefix(accessToken, "eyJhbGciOiJSUzI1NiI") {
		err = utils.UpdateApiServerAccessToken(cfg.ChatGpt.EndPoint.String(), accessToken)
		if err != nil {
			return errors.ErrorAt(err)
		}
	}

	sumer, err := sumer.NewSummarizer(cfg, pushMessage)
	if err != nil {
		return errors.ErrorAt(err)
	}

	err = startSummarize(sumer)
	if err != nil {
		return errors.ErrorAt(err)
	}

	return nil
}

func startSummarize(sumer *sumer.Summarizer) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer close(sigs)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		_, err := sumer.Summarize(ctx, args.codeFolder)
		if err != nil {
			return errors.ErrorAt(err)
		}
		cancel()
		return nil
	})

	select {
	case s := <-sigs:
		log.Errorf("%v", s)
		cancel()
	case <-ctx.Done():
	}

	err := eg.Wait()
	if err != nil {
		return errors.ErrorAt(err)
	}
	return nil
}

func pushMessage(msg *sumer.Message) {
	const maxContentLen = 300
	content := strings.ReplaceAll(msg.Content, "\n", " ")
	runeContent := []rune(content)
	if len(runeContent) > maxContentLen {
		content = string(runeContent[:maxContentLen]) + " ..."
	}
	log.Infof("%s:\n %s", msg.Role, content)
}

func parseArgs() {
	flag.Parse()

	nonFlagArgs := flag.Args()

	if len(nonFlagArgs) < 1 {
		log.Printlnf("Usage: %s [-c CONFIG_FILE] <CODE_FOLDER>", BuildName)
		flag.PrintDefaults()
		log.Fatal("invalid args")
	}

	args.codeFolder = nonFlagArgs[0]
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
