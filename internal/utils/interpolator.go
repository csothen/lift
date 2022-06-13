package utils

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	reName        = regexp.MustCompile(`\(\$(?:\s+)?name(?:\s+)?\$\)`)
	reCount       = regexp.MustCompile(`\(\$(?:\s+)?count(?:\s+)?\$\)`)
	reDownloadURL = regexp.MustCompile(`\(\$(?:\s+)?download_url(?:\s+)?\$\)`)
	reVersion     = regexp.MustCompile(`\(\$(?:\s+)?version(?:\s+)?\$\)`)
	reAdminPass   = regexp.MustCompile(`\(\$(?:\s+)?admin_pass(?:\s+)?\$\)`)
	reDbPass      = regexp.MustCompile(`\(\$(?:\s+)?db_pass(?:\s+)?\$\)`)
	rePluginURLs  = regexp.MustCompile(`\(\$(?:\s+)?plugin_urls(?:\s+)?\$\)`)
)

type Interpolator struct {
	Name        string
	Count       int
	DownloadURL string
	Version     string
	AdminPass   string
	DbPass      string
	PluginURLs  []string
}

func (i *Interpolator) Interpolate(data []byte) []byte {
	if i.Name != "" {
		data = reName.ReplaceAll(data, []byte(i.Name))
	}
	if i.Count != 0 {
		data = reCount.ReplaceAll(data, []byte(fmt.Sprintf("%d", i.Count)))
	}
	if i.DownloadURL != "" {
		data = reDownloadURL.ReplaceAll(data, []byte(i.DownloadURL))
	}
	if i.Version != "" {
		data = reVersion.ReplaceAll(data, []byte(i.Version))
	}
	if i.AdminPass != "" {
		data = reAdminPass.ReplaceAll(data, []byte(i.AdminPass))
	}
	if i.DbPass != "" {
		data = reDbPass.ReplaceAll(data, []byte(i.DbPass))
	}
	if len(i.PluginURLs) > 0 {
		data = rePluginURLs.ReplaceAll(data, []byte(strings.Join(i.PluginURLs, "\n")))
	}
	return data
}
