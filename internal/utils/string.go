package utils

import (
	"crypto/md5"
	"fmt"
	"time"
)

func BuildDeploymentCanonical(usecase, service string, number int) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%d%s", number, time.Now())))
	return fmt.Sprintf("%s-%s-%x", usecase, service, hash)
}
