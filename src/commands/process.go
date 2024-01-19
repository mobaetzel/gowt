package commands

import (
	"github.com/alexflint/go-arg"
	"log"
	"os"
	"word_template_service/src/processors"
	"word_template_service/src/readers"
	"word_template_service/src/utils"
)

type processArgs struct {
	Template string `arg:"-t,required" help:"odt/docx template file name"`
	Data     string `arg:"-d,required" help:"json data file name"`
	Output   string `arg:"-o,required" help:"output file name"`
}

func Process() {
	var args processArgs
	parser, err := arg.NewParser(arg.Config{}, &args)
	if err != nil {
		log.Fatal(err)
	}
	err = parser.Parse(os.Args[2:])
	if err != nil {
		parser.Fail(err.Error())
	}

	reader, err := readers.NewFilepathDocumentReader(args.Template)

	jsonData, err := os.ReadFile(args.Data)
	if err != nil {
		log.Fatal(err)
	}

	values, err := utils.ParseValues(string(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	templateBytes, err := renderers.Render(reader, values)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(args.Output, templateBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
