package endpoints

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"regexp"
	"strings"
	renderers "word_template_service/src/processors"
	"word_template_service/src/readers"
)

//go:embed templates/index_get.html
var indexGet []byte

//go:embed templates/index_post.html
var indexPost string

func IndexGet(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return c.Send(indexGet)
}

func IndexPost(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	if form.File["file"] == nil || len(form.File["file"]) == 0 {
		return errors.New("no file provided")
	}

	file := form.File["file"][0]
	docReader, err := readers.NewMultipartDocumentReader(file)
	if err != nil {
		return err
	}
	defer docReader.Close()

	files, err := renderers.ExtractAndProcessContentFiles(docReader)
	if err != nil {
		return err
	}

	inputTemplate := `<div class="control"><label for="%s">%s</label><input id="%s" name="%s" type="text" required/></div>`

	templateTags := ""
	for _, content := range files {
		reg := regexp.MustCompile("\\{\\{[^}]+}}")
		occs := reg.FindAllString(content, -1)
		for _, occ := range occs {
			occ = regexp.MustCompile("\\{\\{ *").ReplaceAllString(occ, "")
			occ = regexp.MustCompile(" *}}").ReplaceAllString(occ, "")
			if !strings.HasPrefix(occ, ".") {
				continue
			}
			occ = strings.TrimPrefix(occ, ".")
			templateTags += fmt.Sprintf(inputTemplate, occ, occ, occ, occ)
		}
	}

	fmt.Println(indexPost)
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	body := fmt.Sprintf(indexPost, templateTags)
	return c.SendString(body)
}
