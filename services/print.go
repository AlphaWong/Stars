package services

import (
	"errors"
	"io"
	"os"
	"text/template"
)

const (
	ErrorBaseTemplate = "Missing BaseTemplate"
	ErrorOutputPath   = "Missing OutputPath"
)

type Printer interface {
	PrintSlice([]MarkDownRow) error
}

// ensure interface implement is correct
var _ Printer = (*TplPrinter)(nil)

type TplPrinterOption func(tplPrinter *TplPrinter)

func WithBaseTemplate(t *template.Template, err error) TplPrinterOption {
	return func(tplPrinter *TplPrinter) {
		tplPrinter.BaseTemplate = t
	}
}

func WithOutputPath(outputPath string) TplPrinterOption {
	return func(tplPrinter *TplPrinter) {
		tplPrinter.OutputPath = outputPath
	}
}

type TplPrinter struct {
	BaseTemplate *template.Template
	OutputPath   string
}

func NewTplPrinter(setters ...TplPrinterOption) (*TplPrinter, error) {
	tplPrinter := &TplPrinter{
		BaseTemplate: nil,
		OutputPath:   "",
	}

	for _, setter := range setters {
		setter(tplPrinter)
	}

	if tplPrinter.BaseTemplate == nil {
		return nil, errors.New(ErrorBaseTemplate)
	}

	if tplPrinter.OutputPath == "" {
		return nil, errors.New(ErrorOutputPath)
	}

	return tplPrinter, nil
}

func (self *TplPrinter) PrintSlice(markDownRows []MarkDownRow) error {
	os.Remove(self.OutputPath)
	output, _ := os.Create(self.OutputPath)
	defer output.Close()
	// ignore the error from tpl parse
	tpl := template.Must(self.BaseTemplate, nil)
	return Print2Template(output, tpl, markDownRows)
}

func Print2Template(
	wr io.Writer,
	tpl *template.Template,
	markDownRows []MarkDownRow,
) error {
	return tpl.ExecuteTemplate(wr, "layout", markDownRows)
}
