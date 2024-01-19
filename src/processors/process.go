package renderers

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"regexp"
	"strings"
	"text/template"
	"word_template_service/src/readers"
	"word_template_service/src/utils"
)

var (
	tagRegex      = regexp.MustCompile("(\\{[^{]*\\{)((<[^>]+>)|([^}^<]+))*(}[^}]*})")
	tagStartRegex = regexp.MustCompile("\\{[^{]*\\{")
	tagEndRegex   = regexp.MustCompile("}[^}]*}")
	xmlTagRegex   = regexp.MustCompile("<[^>]+>")

	ifParagraphRegexOdt    = regexp.MustCompile("<text:p[^{/]+\\{\\{ *if ([^} ]+) *}}.*?</text:p>")
	rangeParagraphRegexOdt = regexp.MustCompile("<text:p[^{/]+\\{\\{ *range ([^} ]+) *}}.*?</text:p>")
	elseParagraphRegexOdt  = regexp.MustCompile("<text:p[^{/]*\\{\\{ *else *}}.*?</text:p>")
	endParagraphRegexOdt   = regexp.MustCompile("<text:p[^{/]*\\{\\{ *end *}}.*?</text:p>")
)

func Render(reader readers.DocumentReader, dirtyValues utils.Values) ([]byte, error) {
	cleanValues := utils.CleanData(dirtyValues).(utils.Values)

	contentFiles, err := ExtractAndProcessContentFiles(reader)
	if err != nil {
		return nil, err
	}

	transformedFiles := make(map[string][]byte)
	for filename, content := range contentFiles {
		tmpl, err := template.New("template").Parse(content)
		if err != nil {
			return nil, err
		}
		var buffer bytes.Buffer
		err = tmpl.Execute(&buffer, cleanValues)
		if err != nil {
			return nil, err
		}
		transformedFiles[filename] = buffer.Bytes()
	}

	var archive bytes.Buffer
	zipWriter := zip.NewWriter(&archive)

	for _, file := range reader.GetFiles() {
		fileWriter, err := zipWriter.Create(file.Name)
		if err != nil {
			return nil, err
		}

		if _, ok := transformedFiles[file.Name]; ok {
			_, err = fileWriter.Write(transformedFiles[file.Name])
			if err != nil {
				return nil, err
			}
		} else {
			fileReader, err := file.Open()
			if err != nil {
				return nil, err
			}

			_, err = io.Copy(fileWriter, fileReader)
			if err != nil {
				return nil, err
			}
		}
	}

	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return archive.Bytes(), nil
}

func ExtractAndProcessContentFiles(reader readers.DocumentReader) (map[string]string, error) {
	docType, err := determineDocumentType(reader)
	if err != nil {
		return nil, err
	}

	contentFiles, err := extractContentFiles(reader, docType)
	if err != nil {
		return nil, err
	}

	transformedFiles := make(map[string]string)

	for _, contentFile := range contentFiles {
		stringData, err := processContentFile(contentFile, docType)
		if err != nil {
			return nil, err
		}
		transformedFiles[contentFile.Name] = stringData
	}

	return transformedFiles, err
}

func processContentFile(contentFile *zip.File, docType utils.DocumentType) (string, error) {
	doc, err := contentFile.Open()
	if err != nil {
		return "", err
	}
	defer doc.Close()

	data, err := io.ReadAll(doc)
	if err != nil {
		return "", err
	}
	stringData := string(data)

	for _, occ := range tagRegex.FindAll(data, -1) {
		occStr := tagStartRegex.ReplaceAllString(string(occ), "{{ ")
		occStr = tagEndRegex.ReplaceAllString(occStr, " }}")
		occStr = xmlTagRegex.ReplaceAllString(occStr, "")

		stringData = strings.ReplaceAll(stringData, string(occ), occStr)
	}

	if docType == utils.DocumentTypeDocx {

	} else if docType == utils.DocumentTypeOdt {
		for _, ifOcc := range ifParagraphRegexOdt.FindAllString(stringData, -1) {
			content := strings.Split(strings.Split(ifOcc, "{{")[1], "}}")[0]
			stringData = strings.ReplaceAll(stringData, ifOcc, "{{ "+content+" }}")
		}

		for _, rangeOcc := range rangeParagraphRegexOdt.FindAllString(stringData, -1) {
			content := strings.Split(strings.Split(rangeOcc, "{{")[1], "}}")[0]
			stringData = strings.ReplaceAll(stringData, rangeOcc, "{{ "+content+" }}")
		}

		stringData = elseParagraphRegexOdt.ReplaceAllString(stringData, "{{ else }}")
		stringData = endParagraphRegexOdt.ReplaceAllString(stringData, "{{ end }}")
	}

	return stringData, nil
}

// determineDocumentType determines the type of the document based on the files contained in the zip file.
func determineDocumentType(reader readers.DocumentReader) (utils.DocumentType, error) {
	for _, file := range reader.GetFiles() {
		switch file.Name {
		case "word/document.xml":
			return utils.DocumentTypeDocx, nil
		case "content.xml":
			return utils.DocumentTypeOdt, nil
		}
	}
	return utils.DocumentTypeDocx, errors.New("input file is not a valid docx or odt file")
}

// extractContentFiles extracts the files containing the document content from a given zip file.
func extractContentFiles(reader readers.DocumentReader, docType utils.DocumentType) ([]*zip.File, error) {
	switch docType {
	case utils.DocumentTypeDocx:
		return extractDocxContentFiles(reader)
	case utils.DocumentTypeOdt:
		return extractOdtContentFiles(reader)
	default:
		return nil, errors.New("invalid document type")
	}
}

// extractOdtContentFiles extracts the files containing the document content from a given zip file.
func extractOdtContentFiles(reader readers.DocumentReader) ([]*zip.File, error) {
	contentFile, err := reader.GetFile("content.xml")
	if err != nil {
		return nil, err
	}

	stylesFile, err := reader.GetFile("styles.xml")
	if err != nil {
		return nil, err
	}

	contentFiles := []*zip.File{
		contentFile,
		stylesFile,
	}

	return contentFiles, nil
}

// extractDocxContentFiles extracts the files containing the document content from a given zip file.
func extractDocxContentFiles(reader readers.DocumentReader) ([]*zip.File, error) {
	file, err := reader.GetFile("word/document.xml")
	if err != nil {
		return nil, err
	}

	contentFiles := []*zip.File{
		file,
	}

	for _, file := range reader.GetFiles() {
		if strings.HasPrefix(file.Name, "word/header") {
			contentFiles = append(contentFiles, file)
		}

		if strings.HasPrefix(file.Name, "word/footer") {
			contentFiles = append(contentFiles, file)
		}
	}

	return contentFiles, nil
}
