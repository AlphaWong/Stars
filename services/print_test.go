package services

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

func TestNewTplPrinter(t *testing.T) {
	require := require.New(t)
	baseTemplatePath, err := filepath.Abs("../template/starred.md")
	require.NoError(err)
	outputPath, err := filepath.Abs("../out.md")
	require.NoError(err)
	printer, err := NewTplPrinter(
		WithBaseTemplate(template.ParseFiles(baseTemplatePath)),
		WithOutputPath(outputPath),
	)
	require.NoError(err)
	require.NotNil(printer)
}

func TestNewTplPrinterWithEmptyBaseTemplatePath(t *testing.T) {
	require := require.New(t)
	printer, err := NewTplPrinter()
	require.Error(err, ErrorBaseTemplate)
	require.Nil(printer)
}

func TestNewTplPrinterWithEmptyOutputPath(t *testing.T) {
	require := require.New(t)
	baseTemplatePath, err := filepath.Abs("../template/starred.md")
	require.NoError(err)
	printer, err := NewTplPrinter(
		WithBaseTemplate(template.ParseFiles(baseTemplatePath)),
	)
	require.Error(err, ErrorOutputPath)
	require.Nil(printer)
}

func TestPrint2Template(t *testing.T) {
	require := require.New(t)
	input := []MarkDownRow{
		{
			Language: "Go",
			Stars:    "1",
			Items:    "[ [victorspringer/http-cache](https://github.com/victorspringer/http-cache) ]",
		},
		{
			Language: "JavaScript",
			Stars:    "2",
			Items:    "[ [stefanwuthrich/cached-google-places](https://github.com/stefanwuthrich/cached-google-places) ], [ [z](zxy) ]",
		},
	}
	var output strings.Builder
	tpl := template.Must(
		template.New("layout").
			Parse(`# Result
Language|⭐️|Repos
---|---|---
{{ range . }}{{.Language}}|{{.Stars}}|{{.Items}}
{{end}}`,
			))

	err := Print2Template(&output, tpl, input)
	require.NoError(err)
	expected := `# Result
Language|⭐️|Repos
---|---|---
Go|1|[ [victorspringer/http-cache](https://github.com/victorspringer/http-cache) ]
JavaScript|2|[ [stefanwuthrich/cached-google-places](https://github.com/stefanwuthrich/cached-google-places) ], [ [z](zxy) ]
`

	require.Equal(
		expected,
		output.String(),
	)
}

func TestPrintSlice(t *testing.T) {
	require := require.New(t)
	baseTemplatePath, _ := filepath.Abs("../template/starred.md")
	tmpfile, err := ioutil.TempFile("", "out.*.md")
	require.NoError(err)

	defer os.Remove(tmpfile.Name())
	outputPath := tmpfile.Name()
	printer, err := NewTplPrinter(
		WithBaseTemplate(template.ParseFiles(baseTemplatePath)),
		WithOutputPath(outputPath),
	)
	require.NoError(err)

	input := []MarkDownRow{
		{
			Language: "Go",
			Stars:    "1",
			Items:    "[ [victorspringer/http-cache](https://github.com/victorspringer/http-cache) ]",
		},
		{
			Language: "JavaScript",
			Stars:    "2",
			Items:    "[ [stefanwuthrich/cached-google-places](https://github.com/stefanwuthrich/cached-google-places) ], [ [z](zxy) ]",
		},
	}

	printer.PrintSlice(input)
	actual, err := ioutil.ReadFile(tmpfile.Name())
	require.NoError(err)

	require.Equal("![test](https://github.com/AlphaWong/Stars/workflows/test/badge.svg)[![codecov](https://codecov.io/gh/AlphaWong/Stars/branch/master/graph/badge.svg?token=xuILexY8TD)](https://codecov.io/gh/AlphaWong/Stars)\n# Stars\nDo you remember what you star ?\n\n# update\nchange to async request instead waterflow now.\n\n# Run \n```sh\nTOKEN=<GITHUB_TOKEN> go run ./main.go && cp -f ./out.md ./README.md\n```\n\n# GITHUB_TOKEN\n```\nsee https://github.com/settings/tokens\n```\n\n# Github doc\n```\nhttps://docs.github.com/en/free-pro-team@latest/rest/reference/activity#list-repositories-starred-by-a-user\n```\n# Result\nLanguage|⭐️|Repos\n---|---|---\nGo|1|[ [victorspringer/http-cache](https://github.com/victorspringer/http-cache) ]\nJavaScript|2|[ [stefanwuthrich/cached-google-places](https://github.com/stefanwuthrich/cached-google-places) ], [ [z](zxy) ]\n", string(actual))
}
