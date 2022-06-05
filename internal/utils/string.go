package utils

import (
	"crypto/md5"
	"fmt"
	"os"
	"path"
	"time"
)

func BuildDeploymentCanonical(usecase, service string, number int) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%d%s", number, time.Now())))
	return fmt.Sprintf("%s-%s-%x", usecase, service, hash)
}

func BuildDeploymentFolderPath(canonical string) (string, error) {
	return pathFromProjectRoot("deployment_files", canonical)
}

func BuildStaticFolderPath() (string, error) {
	return pathFromProjectRoot("static")
}

func BuildTemplatesFolderPath() (string, error) {
	return pathFromProjectRoot("templates")
}

func pathFromProjectRoot(folders ...string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not read working directory: %w", err)
	}
	// prepend the working directory
	folders = append([]string{wd}, folders...)
	return path.Join(folders...), nil
}
