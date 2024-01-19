package readers

import (
	"archive/zip"
	"errors"
	"os"
)

type FilepathDocumentReader struct {
	documentFileHeader   *os.FileInfo
	documentFile         *os.File
	documentContentFiles []*zip.File
}

func NewFilepathDocumentReader(filepath string) (*FilepathDocumentReader, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	zipFile, err := zip.NewReader(file, fileInfo.Size())
	if err != nil {
		return nil, err
	}

	return &FilepathDocumentReader{
		documentFileHeader:   &fileInfo,
		documentFile:         file,
		documentContentFiles: zipFile.File,
	}, nil
}

func (r *FilepathDocumentReader) GetFiles() []*zip.File {
	return r.documentContentFiles
}

func (r *FilepathDocumentReader) GetFile(name string) (*zip.File, error) {
	for _, file := range r.documentContentFiles {
		if file.Name == name {
			return file, nil
		}
	}
	return nil, errors.New("file not found: " + name)
}

func (r *FilepathDocumentReader) GetDocumentSize() int64 {
	return (*r.documentFileHeader).Size()
}

func (r *FilepathDocumentReader) Close() error {
	return (*r.documentFile).Close()
}
