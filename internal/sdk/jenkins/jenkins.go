package jenkins

import (
	"fmt"
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
