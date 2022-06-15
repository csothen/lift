package utils

import (
	"github.com/csothen/lift/internal/models"
)

func GetAdminCredentials(st models.Type, host string) (string, string) {
	switch st {
	case models.SonarqubeService:
		return "admin", "admin"
	case models.JenkinsService:
		// TODO: Read from file created after the deployment is complete
		return "admin", "password"
	default:
		return "", ""
	}
}
