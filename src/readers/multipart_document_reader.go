package readers

import (
	"archive/zip"
	"errors"
	"mime/multipart"
)

type MultipartDocumentReader struct {
	documentFileHeader   *multipart.FileHeader
	documentFile         *multipart.File
	documentContentFiles []*zip.File
}

func NewMultipartDocumentReader(fileHeader *multipart.FileHeader) (*MultipartDocumentReader, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}

	zipFile, err := zip.NewReader(file, fileHeader.Size)
	if err != nil {
		return nil, err
	}

	return &MultipartDocumentReader{
		documentFileHeader:   fileHeader,
		documentFile:         &file,
		documentContentFiles: zipFile.File,
	}, nil
}

func (r *MultipartDocumentReader) GetFiles() []*zip.File {
	return r.documentContentFiles
}

func (r *MultipartDocumentReader) GetFile(name string) (*zip.File, error) {
	for _, file := range r.documentContentFiles {
		if file.Name == name {
			return file, nil
		}
	}
	return nil, errors.New("file not found: " + name)
}

func (r *MultipartDocumentReader) Close() error {
	return (*r.documentFile).Close()
}
