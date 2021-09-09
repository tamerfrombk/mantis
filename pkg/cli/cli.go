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

	manPage, err := mantis.NewManPageBuilder().Title("mantis").ShortDescription("golang man page generator").LongDescription("TODO").Synopsis("generate man pages from the flag package").Section(1).Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "mantis: %v\n", err)
		return 1
	}

	if err := manPage.Write(); err != nil {
		fmt.Fprintf(os.Stderr, "mantis: %v\n", err)
		return 1
	}

	return 0
}
