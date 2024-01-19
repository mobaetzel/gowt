package endpoints

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"word_template_service/src/processors"
	"word_template_service/src/readers"
	"word_template_service/src/utils"
)

func Process(c *fiber.Ctx) error {
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

	var values = make(utils.Values)
	for key, value := range form.Value {
		values[key] = value[0]
	}
	if err != nil {
		return err
	}

	bytes, err := renderers.Render(docReader, values)
	if err != nil {
		return err
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEOctetStream)
	c.Set(fiber.HeaderContentDisposition, "attachment; filename="+file.Filename)
	c.Status(fiber.StatusOK)

	return c.Send(bytes)
}
