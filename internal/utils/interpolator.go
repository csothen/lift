package utils

import "regexp"

var (
	reName = regexp.MustCompile(`\(\$(?:\s+)?name(?:\s+)?\$\)`)
)

type Interpolator struct {
	Name string
}

func (i *Interpolator) Inteprolate(data []byte) []byte {
	if i.Name != "" {
		data = reName.ReplaceAll(data, []byte(i.Name))
	}
	return data
}
