package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestGetMapKeyASC(t *testing.T) {
	require := require.New(t)
	ramdonMap := map[string][]MarkDownRepo{
		"a": []MarkDownRepo{},
		"#": []MarkDownRepo{},
		"1": []MarkDownRepo{},
		"z": []MarkDownRepo{},
	}
	keys := GetMapKeyASC(ramdonMap)
	require.Equal([]string{"#", "1", "a", "z"}, keys)
}

func TestGetUserStarredRepositoriesTotalPage(t *testing.T) {
	require := require.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	token := "TOKEN"
	userName := "alphawong"
	fetcher, err := NewGitHubFetcher(
		WithToken(token),
		WithUserName(userName),
	)
	require.NoError(err)
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.github.com/users/alphawong/starred",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, "")
			resp.Header.Set("link", `<https://api.github.com/user/5622516/starred?page=2>; rel="next", <https://api.github.com/user/5622516/starred?page=63>; rel="last"`)
			if err != nil {
				return httpmock.NewStringResponse(http.StatusInternalServerError, ""), nil
			}
			return resp, nil
		},
	)
	actual := fetcher.GetUserStarredRepositoriesTotalPage()
	require.Equal(63, actual)
}

func TestParseRawLinkHeader(t *testing.T) {
	require := require.New(t)
	rawHeader := `<https://api.github.com/user/5622516/starred?page=2>; rel="next", <https://api.github.com/user/5622516/starred?page=63>; rel="last"`
	actual := ParseRawLinkHeader(rawHeader)
	require.Equal(63, actual)
}

func TestGetUserAllStarredRepositories(t *testing.T) {
	require := require.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	token := "TOKEN"
	userName := "alphawong"
	fetcher, err := NewGitHubFetcher(
		WithToken(token),
		WithUserName(userName),
	)
	require.NoError(err)
	response1Path, err := filepath.Abs("../mock_data/page_1.json")
	require.NoError(err)
	response2Path, err := filepath.Abs("../mock_data/page_2.json")
	require.NoError(err)
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.github.com/users/alphawong/starred?page=1&per_page=100",
		httpmock.NewJsonResponderOrPanic(
			http.StatusOK,
			httpmock.File(response1Path),
		),
	)
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.github.com/users/alphawong/starred?page=2&per_page=100",
		httpmock.NewJsonResponderOrPanic(
			http.StatusOK,
			httpmock.File(response2Path),
		),
	)
	actual := fetcher.GetUserAllStarredRepositories(2)
	path, err := filepath.Abs("../mock_data/page_total.json")
	require.NoError(err)
	b, err := ioutil.ReadFile(path)
	require.NoError(err)
	var expected UserStarredRepositories
	json.Unmarshal(b, &expected)
	require.Contains(expected, actual[0])
	require.Contains(expected, actual[1])
}

func TestGroupByProgrammingLanguage(t *testing.T) {
	require := require.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	token := "TOKEN"
	userName := "alphawong"
	fetcher, err := NewGitHubFetcher(
		WithToken(token),
		WithUserName(userName),
	)
	require.NoError(err)
	response1Path, err := filepath.Abs("../mock_data/page_1.json")
	require.NoError(err)
	response2Path, err := filepath.Abs("../mock_data/page_2.json")
	require.NoError(err)
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.github.com/users/alphawong/starred?page=1&per_page=100",
		httpmock.NewJsonResponderOrPanic(
			http.StatusOK,
			httpmock.File(response1Path),
		),
	)
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.github.com/users/alphawong/starred?page=2&per_page=100",
		httpmock.NewJsonResponderOrPanic(
			http.StatusOK,
			httpmock.File(response2Path),
		),
	)
	repos := fetcher.GetUserAllStarredRepositories(2)
	grouped := GroupByProgrammingLanguage(repos)
	require.Contains(grouped["Go"], MarkDownRepo{
		FullName: "victorspringer/http-cache",
		HtmlUrl:  "https://github.com/victorspringer/http-cache",
		Language: "Go",
	})
	require.Contains(grouped["JavaScript"], MarkDownRepo{
		FullName: "stefanwuthrich/cached-google-places",
		HtmlUrl:  "https://github.com/stefanwuthrich/cached-google-places",
		Language: "JavaScript",
	})
}

func TestCovert2Slice(t *testing.T) {
	require := require.New(t)
	input := map[string][]MarkDownRepo{
		"Go": {
			{
				FullName: "victorspringer/http-cache",
				HtmlUrl:  "https://github.com/victorspringer/http-cache",
				Language: "Go",
			},
		},
		"JavaScript": {
			{
				FullName: "stefanwuthrich/cached-google-places",
				HtmlUrl:  "https://github.com/stefanwuthrich/cached-google-places",
				Language: "JavaScript",
			}, {
				FullName: "z",
				HtmlUrl:  "zxy",
				Language: "JavaScript",
			},
		},
	}
	slice := Covert2Slice(input)
	require.Contains(slice,
		MarkDownRow{
			Language: "Go",
			Stars:    "1",
			Items:    "[ [victorspringer/http-cache](https://github.com/victorspringer/http-cache) ]",
		},
	)
	require.Contains(slice,
		MarkDownRow{
			Language: "JavaScript",
			Stars:    "2",
			Items:    "[ [stefanwuthrich/cached-google-places](https://github.com/stefanwuthrich/cached-google-places) ], [ [z](zxy) ]",
		},
	)
}
