package services

import (
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
		MarkDownRow{
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
