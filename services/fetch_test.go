package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestNewGitHubFetcherFailWithMissingToken(t *testing.T) {
	require := require.New(t)
	fetcher, err := NewGitHubFetcher(
		WithUserName("alphawong"),
	)
	require.Error(err, ErrorGithubToken)
	require.Nil(fetcher)
}

func TestNewGitHubFetcherFailWithMissingUserName(t *testing.T) {
	require := require.New(t)
	fetcher, err := NewGitHubFetcher(
		WithToken("TOKEN"),
	)
	require.Error(err, ErrorUserName)
	require.Nil(fetcher)
}

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

func TestParseRawLinkHeaderFailWithInvalidLastPageHeader(t *testing.T) {
	require := require.New(t)
	rawHeader := `<https://api.github.com/user/5622516/starred?page=2>; rel="next", <::!2312:#>; rel="last"`
	var s strings.Builder
	log.SetOutput(&s)
	actual := ParseRawLinkHeader(rawHeader)
	require.True(strings.Contains(s.String(), "missing protocol scheme"))
	defer func() {
		log.SetOutput(os.Stdout)
	}()
	require.Equal(0, actual)
}

func TestParseRawLinkHeaderFailWithInvalidLastPageNumberHeader(t *testing.T) {
	require := require.New(t)
	rawHeader := `<https://api.github.com/user/5622516/starred?page=2>; rel="next", <https://api.github.com/user/5622516/starred?page=a>; rel="last"`
	var s strings.Builder
	log.SetOutput(&s)
	actual := ParseRawLinkHeader(rawHeader)
	require.True(strings.Contains(s.String(), "invalid syntax"))
	defer func() {
		log.SetOutput(os.Stdout)
	}()
	require.Equal(0, actual)
}

func TestParseRawLinkHeaderFailWithInvalidHeader(t *testing.T) {
	require := require.New(t)
	rawHeader := ``
	require.Panics(func() {
		ParseRawLinkHeader(rawHeader)
	}, "it should panic if InvalidHeader go to ParseRawLinkHeader")
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
	response3Path, err := filepath.Abs("../mock_data/page_3.json")
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
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.github.com/users/alphawong/starred?page=3&per_page=100",
		httpmock.NewJsonResponderOrPanic(
			http.StatusOK,
			httpmock.File(response3Path),
		),
	)
	repos := fetcher.GetUserAllStarredRepositories(3)
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
	require.Contains(grouped[Others], MarkDownRepo{
		FullName: "victorspringer/http-cache",
		HtmlUrl:  "https://github.com/victorspringer/http-cache",
		Language: "",
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

func TestGetURI(t *testing.T) {
	require := require.New(t)
	query := url.Values{
		"per_page": []string{"100"},
		"page":     []string{"1"},
	}
	actual, err := GetURI(GithubURI, "alphawong", query)
	require.NoError(err)
	require.Equal(
		"https://api.github.com/users/alphawong/starred?page=1&per_page=100",
		actual,
	)
}

func TestGetURIWithInvalidBaseURI(t *testing.T) {
	require := require.New(t)
	query := url.Values{
		"per_page": []string{"100"},
		"page":     []string{"1"},
	}
	actual, err := GetURI("::!2312:#", "alphawong", query)
	require.Error(err, "missing protocol scheme")
	require.Equal("", actual)
}
