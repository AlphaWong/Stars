package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	GithubURI = "https://api.github.com/users/%s/starred"

	Others       = "Others"
	MarkdownStar = "[ [%s](%s) ]"
)

type Fetcher interface {
	GetUsersStars() []MarkDownRow
}

type GitHubFetcher struct {
	Token    string
	UserName string
	H        *http.Client
}

const (
	ErrorGithubToken = "Missing Github token"
	ErrorUserName    = "Missing user name"
)

// ensure interface implement is correct
var _ Fetcher = (*GitHubFetcher)(nil)

type GitHubFetcherOption func(*GitHubFetcher)

func WithToken(token string) GitHubFetcherOption {
	return func(g *GitHubFetcher) {
		g.Token = token
	}
}

func WithUserName(name string) GitHubFetcherOption {
	return func(g *GitHubFetcher) {
		g.UserName = name
	}
}

func NewGitHubFetcher(setters ...GitHubFetcherOption) (*GitHubFetcher, error) {
	g := &GitHubFetcher{
		Token:    "",
		UserName: "",
		H:        &http.Client{},
	}

	for _, setter := range setters {
		setter(g)
	}

	if g.Token == "" {
		return nil, errors.New(ErrorGithubToken)
	}

	if g.UserName == "" {
		return nil, errors.New(ErrorUserName)
	}

	return g, nil
}

func (self *GitHubFetcher) GetUsersStars() []MarkDownRow {
	totalPageCount := self.GetUserStarredRepositoriesTotalPage()
	starredRepositories := self.GetUserAllStarredRepositories(totalPageCount)
	repositories := GroupByProgrammingLanguage(starredRepositories)
	slices := Covert2Slice(repositories)
	return slices
}

func (self *GitHubFetcher) GetUserStarredRepositoriesTotalPage() (totalPage int) {
	rawURI := fmt.Sprintf(GithubURI, self.UserName)
	u, err := url.Parse(rawURI)
	if err != nil {
		log.Print(err.Error())
	}
	query := url.Values{
		"per_page": []string{"100"},
		"page":     []string{"1"},
	}
	u.RawQuery = query.Encode()
	req, _ := http.NewRequest(http.MethodGet, u.String(), nil)
	req.Header.Set("Authorization", fmt.Sprintf("token %s", self.Token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, err := self.H.Do(req)
	if err != nil {
		log.Fatalf("%s", err)
	}
	linkHeader := resp.Header.Get("link")
	totalPage = ParseRawLinkHeader(linkHeader)
	return
}

func ParseRawLinkHeader(rawHeader string) (totalPage int) {
	// rawHeader `<https://api.github.com/user/5622516/starred?per_page=100&page=2>; rel="next", <https://api.github.com/user/5622516/starred?per_page=100&page=18>; rel="last"`
	var links = strings.Split(rawHeader, ",")
	var regex = regexp.MustCompile(`\<(\S+)\>`)
	// only get matched group one for last links
	var lastPageRawURI = regex.FindStringSubmatch(links[1])[1]
	lastPageURI, err := url.Parse(lastPageRawURI)
	if err != nil {
		log.Println(err.Error())
		return
	}
	totalPage, err = strconv.Atoi(lastPageURI.Query().Get("page"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	return
}

func GetURI(baseURI string, userName string, query url.Values) (string, error) {
	rawURI := fmt.Sprintf(baseURI, userName)
	u, err := url.Parse(rawURI)
	if err != nil {
		return "", err
	}
	u.RawQuery = query.Encode()
	return u.String(), nil
}

func (self *GitHubFetcher) GetUserAllStarredRepositories(totalPage int) (userStarredRepositories UserStarredRepositories) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()
	ch := make(chan UserStarredRepositories, totalPage)
	for i := 1; i <= totalPage; i++ {
		go func(pageNum int) {
			// put the uri construction here to avoid data race
			query := url.Values{
				"per_page": []string{"100"},
				"page":     []string{strconv.Itoa(pageNum)},
			}
			uri, err := GetURI(GithubURI, self.UserName, query)
			if err != nil {
				log.Print(err.Error())
			}
			req, _ := http.NewRequest(http.MethodGet, uri, nil)
			req.Header.Set("Authorization", fmt.Sprintf("token %s", self.Token))
			req.Header.Set("Accept", "application/vnd.github.v3+json")
			resp, err := self.H.Do(req)
			if err != nil {
				log.Fatalf("%s", err)
			}
			defer resp.Body.Close()

			var singleUserStarredRepositoriesResponse UserStarredRepositories
			err = json.NewDecoder(resp.Body).Decode(&singleUserStarredRepositoriesResponse)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ch <- singleUserStarredRepositoriesResponse
		}(i)
	}
	var taskProgress = 0
task:
	for true {
		select {
		case repositories := <-ch:
			userStarredRepositories = append(userStarredRepositories, repositories...)
			taskProgress = taskProgress + 1
			if taskProgress == totalPage {
				break task
			}
		case <-ctx.Done():
			break task
		}
	}
	return userStarredRepositories
}

func GroupByProgrammingLanguage(userStarredRepositories UserStarredRepositories) map[string][]MarkDownRepo {
	var repositories = make(map[string][]MarkDownRepo)
	for _, v := range userStarredRepositories {
		var languageKey = v.Language
		if v.Language == "" {
			// handle repo without any language categorizing
			languageKey = Others
		}
		repositories[languageKey] = append(
			repositories[languageKey],
			MarkDownRepo{
				FullName: v.FullName,
				HtmlUrl:  v.HTMLURL,
				Language: v.Language,
			},
		)

	}
	return repositories
}

func Covert2Slice(repositories map[string][]MarkDownRepo) []MarkDownRow {
	keys := GetMapKeyASC(repositories)
	var rows = make([]MarkDownRow, 0, len(keys))
	for _, v := range keys {
		row := MarkDownRow{
			Language: v,
			Stars:    strconv.Itoa(len(repositories[v])),
			Items:    GetInnerReposStr(repositories[v]),
		}
		rows = append(rows, row)
	}
	return rows
}

func GetMapKeyASC(m map[string][]MarkDownRepo) (s []string) {
	for k := range m {
		s = append(s, k)
	}
	sort.Strings(s)
	return
}

func GetInnerReposStr(markDownRepo []MarkDownRepo) (s string) {
	var repositories []string
	for _, v := range markDownRepo {
		repositories = append(repositories, fmt.Sprintf(MarkdownStar, v.FullName, v.HtmlUrl))
	}
	return strings.Join(repositories[:], ", ")
}
