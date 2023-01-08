package samdoc

import (
	"os"
)

const (
	PREFIX  = "samdoc-"
	TEMPDIR = "/tmp"
)

func NewTempFile() (*os.File, error) {
	return os.CreateTemp(TEMPDIR, PREFIX)
}
