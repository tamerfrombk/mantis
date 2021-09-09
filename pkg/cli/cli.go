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

	manPage, err := mantis.NewManPageBuilder().
		Title("mantis").
		ShortDescription("golang man page generator").
		LongDescription(`
mantis uses defined flags as well as programmer-supplied information to generate a man page
that contains the 4 mandatory sections and their conventions as defined here: https://linux.die.net/man/7/man-pages.
`).
		Synopsis("generate man pages using information from the flag package").
		Section(1).
		SeeAlso("https://linux.die.net/man/7/man-pages").
		Build()

	if err != nil {
		fmt.Fprintf(os.Stderr, "mantis: %v\n", err)
		return 1
	}

	if err := manPage.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "mantis: %v\n", err)
		return 1
	}

	return 0
}
