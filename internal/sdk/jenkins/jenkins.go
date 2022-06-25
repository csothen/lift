package jenkins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/csothen/lift/internal/utils"
)

type StatusCode string

const (
	Up   StatusCode = "up"
	Down StatusCode = "down"
)

type CrumbResponse struct {
	Crumb             string `json:"crumb"`
	CrumbRequestField string `json:"crumbRequestField"`
}

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
	jobName := "job-default"

	staticFolderPath, err := utils.BuildStaticFolderPath()
	if err != nil {
		return err
	}

	pathToJob := fmt.Sprintf("%s/jenkins/default.xml", staticFolderPath)
	data, err := os.ReadFile(pathToJob)
	if err != nil {
		return err
	}

	client := &http.Client{}

	cr, err := getCrumb(client, host)
	if err != nil {
		return err
	}

	err = createJob(client, jobName, host, data, cr)
	if err != nil {
		return err
	}

	err = triggerJob(client, jobName, host, cr)
	if err != nil {
		return err
	}
	return nil
}

func getCrumb(client *http.Client, host string) (*CrumbResponse, error) {
	getCrumbURL := fmt.Sprintf("http://%s:8080/crumbIssuer/api/json", host)

	req, err := http.NewRequest(http.MethodGet, getCrumbURL, nil)
	if err != nil {
		return nil, fmt.Errorf("could not build request: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	crumbData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var cr CrumbResponse
	err = json.Unmarshal(crumbData, &cr)
	if err != nil {
		return nil, err
	}
	return &cr, nil
}

func createJob(client *http.Client, jobName, host string, body []byte, crumb *CrumbResponse) error {
	createJobURL := fmt.Sprintf("http://%s:8080/createItem?name=%s", host, jobName)

	// create the request to create the job
	req, err := http.NewRequest(http.MethodPost, createJobURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("could not build request: %w", err)
	}

	// Add content type header
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add(crumb.CrumbRequestField, crumb.Crumb)
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create job")
	}
	return nil
}

func triggerJob(client *http.Client, jobName, host string, crumb *CrumbResponse) error {
	triggerJobURL := fmt.Sprintf("http://%s:8080/job/%s/build", host, jobName)

	req, err := http.NewRequest(http.MethodPost, triggerJobURL, nil)
	if err != nil {
		return fmt.Errorf("could not build request: %w", err)
	}

	req.Header.Add(crumb.CrumbRequestField, crumb.Crumb)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to trigger job")
	}
	return nil
}
