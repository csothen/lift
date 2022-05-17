package sonarqube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/csothen/tmdei-project/internal/fetcher/types"
	"golang.org/x/net/html"
)

const pluginsSource string = "https://update.sonarsource.org/"

type pluginData struct {
	Name     string              `json:"name"`
	Key      string              `json:"key"`
	Versions []pluginVersionData `json:"versions"`
}

type pluginVersionData struct {
	Version string `json:"version"`
	URL     string `json:"downloadURL"`
}

func (f *fetcher) fetchPlugins(wg *sync.WaitGroup, psURL *url.URL) {
	defer wg.Done()
	resp, err := http.Get(pluginsSource)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	errs := make([]error, 0)
	var m sync.Mutex

	page := html.NewTokenizer(resp.Body)
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			break
		}

		token := page.Token()
		switch tokenType {
		case html.StartTagToken:
			if token.DataAtom.String() != "a" {
				continue
			}

			for _, attr := range token.Attr {
				if attr.Key != "href" {
					continue
				}

				// we don't care about refs to non JSON files
				if !strings.HasSuffix(attr.Val, ".json") {
					continue
				}

				urlToCheck, err := url.Parse(attr.Val)
				if err != nil {
					errs = append(errs, err)
					continue
				}

				// we don't care about external links
				if urlToCheck.Host != psURL.Host && len(urlToCheck.Host) != 0 {
					continue
				}

				if len(urlToCheck.Host) == 0 {
					urlToCheck.Scheme = psURL.Scheme
					urlToCheck.Host = psURL.Host
				}

				wg.Add(1)
				go f.fetchPluginInformation(wg, &m, urlToCheck)
				break
			}
		}
	}
	for _, err := range errs {
		log.Println(err)
	}
}

func (f *fetcher) fetchPluginInformation(wg *sync.WaitGroup, m *sync.Mutex, urlToCheck *url.URL) {
	defer wg.Done()
	res, err := http.Get(urlToCheck.String())
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var pd pluginData
	err = json.Unmarshal(body, &pd)
	if err != nil {
		log.Println(err)
		return
	}

	m.Lock()
	for _, v := range pd.Versions {
		f.plugins[buildPluginKey(pd.Key, v.Version)] = &types.Plugin{
			Name:        pd.Name,
			Version:     v.Version,
			DownloadURL: v.URL,
		}
	}
	m.Unlock()
}

func buildPluginKey(name, version string) string {
	return fmt.Sprintf("%s:%s", name, version)
}
