package sonarqube

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/csothen/tmdei-project/internal/fetcher/types"
)

const sonarqubeSource = "https://downloads-cdn-eu-central-1-prod.s3.eu-central-1.amazonaws.com/?delimiter=/&prefix=Distribution%2Fsonarqube%2F"

type sqVersionsXML struct {
	Contents []sqVersionContents `xml:"Contents"`
}

type sqVersionContents struct {
	Key string `xml:"Key"`
}

func (f *fetcher) fetchSonarqubeVersions(wg *sync.WaitGroup, ssURL *url.URL) {
	defer wg.Done()

	sonarqubeVersionRe := regexp.MustCompile(`sonar(qube)?-(.*).zip$`)

	resp, err := http.Get(sonarqubeSource)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	xmlBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var sqVersions sqVersionsXML
	xml.Unmarshal(xmlBody, &sqVersions)

	for _, sqVersion := range sqVersions.Contents {
		//we don't care about non ZIP files
		if !strings.HasSuffix(sqVersion.Key, ".zip") {
			continue
		}

		submatches := sonarqubeVersionRe.FindStringSubmatch(sqVersion.Key)
		if len(submatches) < 2 {
			continue
		}

		downloadURL := fmt.Sprintf("https://binaries.sonarsource.com/%s", sqVersion.Key)
		version := submatches[len(submatches)-1]

		f.SonarqubeVersions[version] = &types.AppVersion{
			Version:     version,
			DownloadURL: downloadURL,
		}
	}
}
