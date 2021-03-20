package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Setenv("TOKEN", "TOKEN")
	m.Run()
}

func TestBoot(t *testing.T) {
	require := require.New(t)
	config := boot()
	require.NotNil(config)
}

func TestValidConfigWithMissingToken(t *testing.T) {
	require := require.New(t)
	require.Panics(func() {
		config := boot()
		config.Token = ""
		validConfig(config)
	})
}

func TestRun(t *testing.T) {
	require := require.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	response1Path, err := filepath.Abs("./mock_data/page_1.json")
	require.NoError(err)
	response2Path, err := filepath.Abs("./mock_data/page_2.json")
	require.NoError(err)
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.github.com/users/alphawong/starred?page=1&per_page=100",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, httpmock.File(response1Path).String())
			resp.Header.Set("link", `<https://api.github.com/user/5622516/starred?page=2>; rel="next", <https://api.github.com/user/5622516/starred?page=2>; rel="last"`)
			if err != nil {
				return httpmock.NewStringResponse(http.StatusInternalServerError, ""), nil
			}
			return resp, nil
		},
	)
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.github.com/users/alphawong/starred?page=2&per_page=100",
		httpmock.NewJsonResponderOrPanic(
			http.StatusOK,
			httpmock.File(response2Path),
		),
	)

	tmpfile, err := ioutil.TempFile(".", "tpl.*.md")
	require.NoError(err)

	_, err = tmpfile.WriteString(`{{define "layout"}}# Result
Language|⭐️|Repos
---|---|---
{{ range . }}{{.Language}}|{{.Stars}}|{{.Items}}
{{end}}{{end}}`,
	)
	require.NoError(err)
	defer os.Remove(tmpfile.Name())

	outputFile, err := ioutil.TempFile(".", "text-out.*.md")
	defer os.Remove(outputFile.Name())

	config := boot()
	config.mu.Lock()
	defer config.mu.Unlock()
	config.BaseTemplate = tmpfile.Name()
	config.OutputPath = outputFile.Name()
	run(config)

	info := httpmock.GetCallCountInfo()
	log.Println(info)

	actual, err := ioutil.ReadFile(outputFile.Name())
	require.NoError(err)
	require.Equal("# Result\nLanguage|⭐️|Repos\n---|---|---\nGo|1|[ [victorspringer/http-cache](https://github.com/victorspringer/http-cache) ]\nJavaScript|1|[ [stefanwuthrich/cached-google-places](https://github.com/stefanwuthrich/cached-google-places) ]\n", string(actual))
}
