package main

import (
	"embed"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/jopbrown/codesum/cmd/codesum-gui/app"
	"github.com/jopbrown/gobase/fsutil"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

var (
	BuildName    = "myapp"
	BuildVersion = "v0.0.0"
	BuildHash    = "unknown"
	BuildTime    = "20060102150405"
)

//go:embed all:frontend/dist
var assets embed.FS

var cli struct {
	configPath string
}

func init() {
	flag.StringVar(&cli.configPath, "c", filepath.Join(fsutil.AppDir(), "config.yml"), "config path")
	flag.Parse()
}

func main() {
	// Create an instance of the app structure
	a := app.NewApp(cli.configPath)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  fmt.Sprintf("%s %v-%v-%v", BuildName, BuildVersion, BuildHash, BuildTime),
		Logger: app.NewLogger(),
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        a.Startup,
		Bind: []interface{}{
			a,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
