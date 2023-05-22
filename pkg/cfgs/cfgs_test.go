package cfgs

import (
	"os"
	"testing"
	"time"

	"github.com/jopbrown/gobase/fsutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("CHATGPT_END_POINT", "https://api.openai.com/v1")
	os.Setenv("CHATGPT_ACCESS_KEY", "apikey")

	cfg, err := LoadConfig(`default/default.yml`)
	require.NoError(t, err)

	dict := map[string]string{
		"appDir":         fsutil.AppDir(),
		"timestamp":      time.Now().Format("20060102150405"),
		"folderBaseName": "mycode",
	}

	// err = cfg.SaveConfig(`tmp/testconfig.yml`)
	// require.NoError(t, err)

	assert.Equal(t, "https://api.openai.com/v1", cfg.ChatGpt.EndPoint.String())
	assert.Equal(t, "", cfg.ChatGpt.Proxy.String())
	assert.Equal(t, fsutil.AppDir()+"/summary_report", cfg.SummaryRules.OutDir.ExpandByDict(dict))
	assert.Equal(t, "summary_"+dict["timestamp"]+"_mycode.md", cfg.SummaryRules.OutFileName.ExpandByDict(dict))
}
