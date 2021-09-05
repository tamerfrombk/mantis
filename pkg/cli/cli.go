package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/tamerfrombk/mantis/pkg/mantis"
)

func Run() int {
	isHelp := flag.Bool("h", false, "displays this help message")

	flag.Parse()

	if *isHelp {
		flag.Usage()
		return 0
	}

	manPage := mantis.NewManPage()
	manPage.Title = "mantis"
	manPage.Description = "golang man page generator"

	if err := manPage.Write(); err != nil {
		fmt.Fprintf(os.Stderr, "mantis: unable to write the man page %v\n", err)
		return 1
	}

	return 0
}
