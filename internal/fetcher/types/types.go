package types

type Plugin struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	DownloadURL string `json:"download_url"`
}

type AppVersion struct {
	Version     string `json:"version"`
	DownloadURL string `json:"download_url"`
}
