package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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
	token = "TOKEN"
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
	actual := GetUserStarredRepositoriesTotalPage()
	require.Equal(63, actual)
}

func TestParseRawLinkHeader(t *testing.T) {
	require := require.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	token = "TOKEN"
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
	actual := GetUserStarredRepositoriesTotalPage()
	require.Equal(63, actual)
}

func TestGetUserAllStarredRepositories(t *testing.T) {
	require := require.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	token = "TOKEN"
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.github.com/users/alphawong/starred?page=1&per_page=100",
		httpmock.NewJsonResponderOrPanic(
			http.StatusOK,
			httpmock.File("./mock_data/page_1.json"),
		),
	)
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.github.com/users/alphawong/starred?page=2&per_page=100",
		httpmock.NewJsonResponderOrPanic(
			http.StatusOK,
			httpmock.File("./mock_data/page_2.json"),
		),
	)
	actual := GetUserAllStarredRepositories(2)
	b, err := ioutil.ReadFile("./mock_data/page_total.json")
	require.NoError(err)
	var expected UserStarredRepositories
	json.Unmarshal(b, &expected)
	require.Contains(expected, actual[0])
	require.Contains(expected, actual[1])
}
