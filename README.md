# mantis

Never write another man page for one of your Go tools again! 

`mantis` is a Go package for quickly generating man pages using information from the `flag` package.

### Quick Start

See the `Run()` function in [cli.go](pkg/cli/cli.go) as an example. 

Basically, it boils down to defining your flags and creating the man page:

```
import "github.com/tamerfrombk/mantis/pkg/mantis"
import "flag"

func MyFunc() {
    // define all flags
    flag.Bool("h", false, "displays this help message")

    // other flags
    // ...

    // generate the man page
    // ALL of these fields are required
    manPage, err := mantis.NewManPageBuilder().
    Title("foo").
    ShortDescription("bar").
    LongDescription("baz").
    Synopsis("synopsis").
    Section(1).
    SeeAlso("https://linux.die.net/man/7/man-pages").
    Build()

	if err != nil {
        // handle the error
	}

    // write the man page to disk
	if err := manPage.Save(); err != nil {
		// handle the error
	}
}
```

### Examples

`mantis` uses itself to generate its own manpage :).

Generate the man page by running `go run ./cmd/mantis/main.go`. This will generate a `mantis.1` man file in the current directory.

To view the man page, use the `man` command: `man -l ./mantis.1`.

### Limitations

1. `mantis` does _not_ install the man page to your system. It simply writes it out to the current working directory.  