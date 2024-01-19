package commands

import (
	_ "embed"
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"word_template_service/src/endpoints"
)

type serveArgs struct {
	Host  string `arg:"-x" default:"127.0.0.1" help:"host to serve on"`
	Port  string `arg:"-p" default:"3000" help:"port to serve on"`
	WebUI bool   `arg:"-w" default:"false" help:"enable web ui"`
}

func Serve() {
	var args serveArgs
	parser, err := arg.NewParser(arg.Config{}, &args)
	if err != nil {
		log.Fatal(err)
	}
	err = parser.Parse(os.Args[2:])
	if err != nil {
		parser.Fail(err.Error())
	}

	app := fiber.New(fiber.Config{})

	app.Post("/process", endpoints.Process)

	if args.WebUI {
		app.Get("/", endpoints.IndexGet)
		app.Post("/", endpoints.IndexPost)
	}

	addr := fmt.Sprintf("%s:%s", args.Host, args.Port)

	log.Fatal(app.Listen(addr))
}
