package mantis

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	ce "github.com/tamerfrombk/composite_error/pkg"
)

type ManPage struct {
	section          int
	flags            map[string]*flag.Flag
	title            string
	shortDescription string
	longDescription  string
	synopsis         string
	seeAlso          string
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

func (b ManPageBuilder) SeeAlso(str string) ManPageBuilder {
	b.instance.seeAlso = str

	return b
}

func (b ManPageBuilder) Build() (ManPage, error) {
	errs := ce.NewCompositeError()
	flagVisitor := func(f *flag.Flag) {
		if err := b.instance.addFlag(f.Name, f); err != nil {
			errs.Add(err)
		}
	}

	flag.VisitAll(flagVisitor)

	if b.instance.Section() < 1 || b.instance.Section() > 8 {
		errs.Add(errors.New("sections must be between 1 and 8"))
	}

	if b.instance.Title() == "" {
		errs.Add(errors.New("title must be set"))
	}

	if b.instance.ShortDescription() == "" {
		errs.Add(errors.New("short description must be set"))
	}

	if b.instance.Synopsis() == "" {
		errs.Add(errors.New("synopsis must be set"))
	}

	if b.instance.SeeAlso() == "" {
		errs.Add(errors.New("see also must be set"))
	}

	return *b.instance, errs.Value()
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

func (w *ManPage) SeeAlso() string {
	return w.seeAlso
}

// Save convenience function to write the man page to a file in the cwd
// The file's name will follow the convention <title>.<section>
func (m *ManPage) Save() error {
	f, err := os.Create(m.Title() + "." + strconv.Itoa(m.Section()))
	if err != nil {
		return err
	}
	defer f.Close()

	_, e := m.WriteTo(f)

	return e
}

// WriteTo writes the man page to the supplied writer
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
		m.writeSeeAlso,
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

	data := fmt.Sprintf(".TH %s %d %v %s %s\n", m.Title(), m.Section(), dateStr, "Linux", "Linux Programmer's Manual")

	w.Write([]byte(data))

	return nil
}

func (m *ManPage) writeName(w io.Writer) error {
	str := ".SH NAME" + "\n" + m.Title() + " \\- " + m.ShortDescription() + "\n"
	if _, err := w.Write([]byte(str)); err != nil {
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

func (m *ManPage) writeSeeAlso(w io.Writer) error {
	str := ".sh SEE ALSO" + "\n" + m.SeeAlso()

	if _, err := w.Write([]byte(str)); err != nil {
		return err
	}

	return nil
}
