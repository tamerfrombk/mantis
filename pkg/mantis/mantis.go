package mantis

import (
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
	section          int
	flags            map[string]*flag.Flag
	title            string
	shortDescription string
	longDescription  string
	synopsis         string
}

type ManPageBuilder struct {
	instance *ManPage
}

func NewManPageBuilder() *ManPageBuilder {
	return &ManPageBuilder{
		instance: &ManPage{
			section: 1,
			flags:   make(map[string]*flag.Flag),
		},
	}
}

func (b ManPageBuilder) Section(section int) ManPageBuilder {
	b.instance.section = section

	return b
}

func (b ManPageBuilder) Title(title string) ManPageBuilder {
	b.instance.title = title

	return b
}

func (b ManPageBuilder) ShortDescription(desc string) ManPageBuilder {
	b.instance.shortDescription = desc

	return b
}

func (b ManPageBuilder) LongDescription(desc string) ManPageBuilder {
	b.instance.longDescription = desc

	return b
}

func (b ManPageBuilder) Synopsis(synopsis string) ManPageBuilder {
	b.instance.synopsis = synopsis

	return b
}

func (b ManPageBuilder) Build() (ManPage, error) {
	errs := NewCompositeError()
	flagVisitor := func(f *flag.Flag) {
		if err := b.instance.addFlag(f.Name, f); err != nil {
			errs.Add(err)
		}
	}

	flag.VisitAll(flagVisitor)

	if b.instance.section < 1 || b.instance.section > 8 {
		errs.Add(errors.New("sections can only be between 1 and 8"))
	}

	if b.instance.title == "" {
		errs.Add(errors.New("title must be set"))
	}

	if b.instance.shortDescription == "" {
		errs.Add(errors.New("short description must be set"))
	}

	if b.instance.synopsis == "" {
		errs.Add(errors.New("synopsis must be set"))
	}

	if errs.IsEmpty() {
		return *b.instance, nil
	}

	return *b.instance, errs
}

func (w *ManPage) Section() int {
	return w.section
}

func (w *ManPage) Title() string {
	return w.title
}

func (w *ManPage) ShortDescription() string {
	return w.shortDescription
}

func (w *ManPage) LongDescription() string {
	return w.longDescription
}

func (w *ManPage) Synopsis() string {
	return w.synopsis
}

// Write convenience function to write the man page to a default path in the cwd
// This method follows the conventions laid out in https://linux.die.net/man/7/man-pages
func (m *ManPage) Write() error {
	f, err := os.Create(m.title + "." + "man")
	if err != nil {
		return err
	}
	defer f.Close()

	_, e := m.WriteTo(f)

	return e
}

// WriteTo writes the man page to the supplied writer following the conventions
// laid out in https://linux.die.net/man/7/man-pages
func (m *ManPage) WriteTo(w io.Writer) (int64, error) {
	text, err := m.MarshalText()
	if err != nil {
		return 0, err
	}

	n, err := w.Write(text)

	return int64(n), err
}

func (m *ManPage) MarshalText() ([]byte, error) {
	writes := []func(io.Writer) error{
		m.writeTitleLine,
		m.writeName,
		m.writeSynopsis,
		m.writeDescription,
		m.writeOptions,
	}

	buf := strings.Builder{}
	for _, write := range writes {
		if err := write(&buf); err != nil {
			return nil, err
		}
	}

	return []byte(buf.String()), nil
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

	data := fmt.Sprintf(".TH %s %d %v %s %s\n", m.Title(), m.section, dateStr, "Linux", "Linux Programmer's Manual")

	w.Write([]byte(data))

	return nil
}

func (m *ManPage) writeName(w io.Writer) error {
	if _, err := w.Write([]byte(".SH NAME\n")); err != nil {
		return err
	}

	if _, err := w.Write([]byte(m.Title() + " \\- " + m.ShortDescription() + "\n")); err != nil {
		return err
	}

	return nil
}

func (m *ManPage) writeSynopsis(w io.Writer) error {
	synopsis := strings.Join([]string{
		".SH SYNOPSIS", m.Synopsis(),
	}, "\n") + "\n"

	if _, err := w.Write([]byte(synopsis)); err != nil {
		return err
	}

	return nil
}

func (m *ManPage) writeDescription(w io.Writer) error {
	description := strings.Join([]string{
		".SH DESCRIPTION", m.LongDescription(),
	}, "\n") + "\n"

	if _, err := w.Write([]byte(description)); err != nil {
		return err
	}

	return nil
}

func (m *ManPage) writeOptions(w io.Writer) error {
	if _, err := w.Write([]byte(".SH OPTIONS\n")); err != nil {
		return err
	}

	prevOutput := flag.CommandLine.Output()

	flag.CommandLine.SetOutput(w)

	flag.PrintDefaults()

	flag.CommandLine.SetOutput(prevOutput)

	return nil
}
