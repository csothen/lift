package sonarqube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	Up               StatusCode = "UP"
	Starting         StatusCode = "STARTING"
	Down             StatusCode = "DOWN"
	Restarting       StatusCode = "RESTARTING"
	MigrationNeeded  StatusCode = "DB_MIGRATION_NEEDED"
	MigrationRunning StatusCode = "DB_MIGRATION_RUNNING"
)

type StatusCode string

type statusCheckBody struct {
	Version string     `json:"version"`
	Status  StatusCode `json:"status"`
}

func Status(host string) (StatusCode, error) {
	statusCheckURL := fmt.Sprintf("http://%s:9000/api/system/status", host)

	client := &http.Client{}
	// create the request to create the user
	scReq, err := http.NewRequest(http.MethodGet, statusCheckURL, nil)
	if err != nil {
		return "", fmt.Errorf("could not build request: %w", err)
	}

	res, err := client.Do(scReq)
	if err != nil || res.StatusCode != 200 {
		return "", fmt.Errorf("could not check instance status")
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var check statusCheckBody
	err = json.Unmarshal(data, &check)
	if err != nil {
		return "", fmt.Errorf("failed to parse response body: %w", err)
	}

	return check.Status, nil
}

func Setup(host, username, password, adminUsername, adminPassword string) error {
	createUserURL := fmt.Sprintf("http://%s:9000/api/users/create", host)

	// Create request parameters
	parameters := url.Values{}
	parameters.Set("login", username)
	parameters.Add("name", username)
	parameters.Add("password", password)

	// Build URL with parameters
	createUserURL = fmt.Sprintf("%s?%s", createUserURL, parameters.Encode())

	client := &http.Client{}
	// create the request to create the user
	cuReq, err := http.NewRequest(http.MethodPost, createUserURL, nil)
	if err != nil {
		return fmt.Errorf("could not build request: %w", err)
	}

	// set the basic authentication
	cuReq.SetBasicAuth(adminUsername, adminPassword)
	res, err := client.Do(cuReq)
	if err != nil {
		return fmt.Errorf("could not do request: %w", err)
	}

	if res.StatusCode != 200 {
		errData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		return fmt.Errorf("failed to create user: \n%s", string(errData))
	}

	for _, permission := range []string{"scan", "provisioning"} {
		err = addPermissionToUser(host, username, permission, adminUsername, adminPassword)
		if err != nil {
			return err
		}
	}
	return nil
}

func addPermissionToUser(host, username, permission, adminUsername, adminPassword string) error {
	addPermissionsURL := fmt.Sprintf("http://%s:9000/api/permissions/add_user", host)

	// Create request parameters
	parameters := url.Values{}
	parameters.Set("login", username)
	parameters.Add("permission", permission)

	// Build URL with parameters
	addPermissionsURL = fmt.Sprintf("%s?%s", addPermissionsURL, parameters.Encode())

	client := &http.Client{}
	// create the request to add the permissions to the user
	apReq, err := http.NewRequest(http.MethodPost, addPermissionsURL, nil)
	if err != nil {
		return fmt.Errorf("could not build request: %w", err)
	}

	// set the basic authentication
	apReq.SetBasicAuth(adminUsername, adminPassword)
	res, err := client.Do(apReq)
	if err != nil {
		return fmt.Errorf("could not do request: %w", err)
	}

	if res.StatusCode != 204 {
		errData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to add permissions to user: %w", err)
		}
		return fmt.Errorf("failed to add permissions to user: \n%s", string(errData))
	}
	return nil
}
