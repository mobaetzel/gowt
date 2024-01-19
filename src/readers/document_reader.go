package readers

import (
	"archive/zip"
)

type DocumentReader interface {
	GetFiles() []*zip.File
	GetFile(name string) (*zip.File, error)
	Close() error
}
