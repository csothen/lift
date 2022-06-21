package jenkins

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type StatusCode string

const (
	Up   StatusCode = "up"
	Down StatusCode = "down"
)

func Status(host string) (StatusCode, error) {
	statusCheckURL := fmt.Sprintf("http://%s:8080", host)

	client := &http.Client{}
	// create the request to create the user
	scReq, err := http.NewRequest(http.MethodGet, statusCheckURL, nil)
	if err != nil {
		return Down, fmt.Errorf("could not build request: %w", err)
	}

	res, err := client.Do(scReq)
	if err != nil || res.StatusCode != 200 {
		return Down, fmt.Errorf("could not check instance status")
	}
	return Up, nil
}

func Setup(host string) error {
	jobName := "default"
	createJobURL := fmt.Sprintf("http://%s:8080/job/createItem?name=%s", host, jobName)
	triggerJobURL := fmt.Sprintf("http://%s:8080/job/%s/build", host, jobName)

	data, err := ioutil.ReadFile("config.xml")
	if err != nil {
		return err
	}

	client := &http.Client{}
	// create the request to create the job
	cjReq, err := http.NewRequest(http.MethodPost, createJobURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("could not build request: %w", err)
	}

	// Add content type header
	cjReq.Header.Add("Content-Type", "text/xml")
	_, err = client.Do(cjReq)
	if err != nil {
		return err
	}

	tjReq, err := http.NewRequest(http.MethodPost, triggerJobURL, nil)
	if err != nil {
		return fmt.Errorf("could not build request: %w", err)
	}

	_, err = client.Do(tjReq)
	if err != nil {
		return err
	}
	return nil
}
