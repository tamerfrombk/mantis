package mantis

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type CompositeError struct {
	errors []error
}

func NewCompositeError() *CompositeError {
	return &CompositeError{
		errors: make([]error, 0),
	}
}

func (e *CompositeError) Add(err error) {
	e.errors = append(e.errors, err)
}

func (e *CompositeError) IsEmpty() bool {
	return len(e.errors) == 0
}

func (e *CompositeError) Error() string {
	if e.IsEmpty() {
		return ""
	}

	b := strings.Builder{}

	b.WriteString(e.errors[0].Error())
	for i := 1; i < len(e.errors); i++ {
		b.WriteString(", ")
		b.WriteString(e.errors[i].Error())
		if i%2 == 0 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

type ManPage struct {
	section     int
	flags       map[string]*flag.Flag
	Title       string
	Description string
}

func NewManPage() *ManPage {
	return &ManPage{
		section: 1,
		flags:   make(map[string]*flag.Flag),
	}
}

func (w *ManPage) SetSection(section int) error {
	if section < 1 || section > 8 {
		return errors.New("sections can only be between 1 and 8")
	}

	w.section = section

	return nil
}

func (w *ManPage) Section() int {
	return w.section
}

func (w *ManPage) Parse() error {
	errs := NewCompositeError()
	flagVisitor := func(f *flag.Flag) {
		if err := w.addFlag(f.Name, f); err != nil {
			errs.Add(err)
		}
	}

	flag.VisitAll(flagVisitor)

	if w.Title == "" {
		errs.Add(errors.New("title must be set"))
	}

	if w.Description == "" {
		errs.Add(errors.New("description must be set"))
	}

	if errs.IsEmpty() {
		return nil
	}

	return errs
}

// Write convenience function to write the man page to a default path in the cwd
// This method follows the conventions laid out in https://linux.die.net/man/7/man-pages
func (w *ManPage) Write() error {
	return w.WriteTo(w.Title + "." + "man")
}

// WriteTo writes the man page to the supplied path following the conventions
// laid out in https://linux.die.net/man/7/man-pages
func (m *ManPage) WriteTo(path string) error {
	if err := m.Parse(); err != nil {
		return errors.New("failed to parse: " + err.Error())
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := bufio.NewWriter(f)
	if err := m.writeTitleLine(buf); err != nil {
		return err
	}

	if err := m.writeName(buf); err != nil {
		return err
	}
	if err := m.writeSynopsis(buf); err != nil {
		return err
	}

	return buf.Flush()
}

func (w *ManPage) addFlag(name string, f *flag.Flag) error {
	if _, found := w.flags[name]; found {
		return errors.New("duplicate flag: " + name)
	}

	w.flags[name] = f

	return nil
}

func (m *ManPage) writeTitleLine(w io.Writer) error {
	dateStr := strings.Split(time.Now().String(), " ")[0]

	data := fmt.Sprintf(".TH %s %d %v %s %s\n", m.Title, m.section, dateStr, "Linux", "Linux Programmer's Manual")

	w.Write([]byte(data))

	return nil
}

func (m *ManPage) writeName(w io.StringWriter) error {
	if _, err := w.WriteString(".SH NAME\n"); err != nil {
		return err
	}

	if _, err := w.WriteString(m.Title + " \\- " + m.Description + "\n"); err != nil {
		return err
	}

	return nil
}

func (m *ManPage) writeSynopsis(w io.Writer) error {
	if _, err := w.Write([]byte(".SH SYNOPSIS\n")); err != nil {
		return err
	}

	prevOutput := flag.CommandLine.Output()

	flag.CommandLine.SetOutput(w)

	flag.PrintDefaults()

	flag.CommandLine.SetOutput(prevOutput)

	return nil
}
