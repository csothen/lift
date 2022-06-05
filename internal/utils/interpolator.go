package utils

import (
	"fmt"
	"regexp"
)

var (
	reName  = regexp.MustCompile(`\(\$(?:\s+)?name(?:\s+)?\$\)`)
	reCount = regexp.MustCompile(`\(\$(?:\s+)?count(?:\s+)?\$\)`)
)

type Interpolator struct {
	Name  string
	Count int
}

func (i *Interpolator) Inteprolate(data []byte) []byte {
	if i.Name != "" {
		data = reName.ReplaceAll(data, []byte(i.Name))
	}
	if i.Count != 0 {
		data = reCount.ReplaceAll(data, []byte(fmt.Sprintf("%d", i.Count)))
	}
	return data
}
