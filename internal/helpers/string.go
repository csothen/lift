package helpers

import "fmt"

func BuildDeploymentCanonical(usecase, service string, number int) string {
	return fmt.Sprintf("%s-%s-%d", usecase, service, number)
}
