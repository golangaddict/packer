package packer

import (
	"io/ioutil"
	"strings"
)

type errStringBuilder struct {
	builder strings.Builder
	err     error
}

func (sb *errStringBuilder) writeBytesFromFile(f string) {
	if sb.err != nil {
		return
	}

	b, err := ioutil.ReadFile(f)
	if err != nil {
		sb.err = err
		return
	}

	sb.builder.Write(b)
}

func (sb *errStringBuilder) String() string {
	return sb.builder.String()
}
