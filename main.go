package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const (
	Others       = "Others"
	MarkdownStar = "[ [%s](%s) ]"

	GitHubUserName = "alphawong"
	// GithubURI store the base uri
	// "https://api.github.com/users/alphawong/starred"
	GithubURI = "https://api.github.com/users/%s/starred"
)

var (
	token  = os.Getenv("TOKEN")
	client = &http.Client{}
)

type UserStarredRepositories []struct {
	ID       int    `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	Owner    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"owner"`
	HTMLURL          string      `json:"html_url"`
	Description      string      `json:"description"`
	Fork             bool        `json:"fork"`
	URL              string      `json:"url"`
	ForksURL         string      `json:"forks_url"`
	KeysURL          string      `json:"keys_url"`
	CollaboratorsURL string      `json:"collaborators_url"`
	TeamsURL         string      `json:"teams_url"`
	HooksURL         string      `json:"hooks_url"`
	IssueEventsURL   string      `json:"issue_events_url"`
	EventsURL        string      `json:"events_url"`
	AssigneesURL     string      `json:"assignees_url"`
	BranchesURL      string      `json:"branches_url"`
	TagsURL          string      `json:"tags_url"`
	BlobsURL         string      `json:"blobs_url"`
	GitTagsURL       string      `json:"git_tags_url"`
	GitRefsURL       string      `json:"git_refs_url"`
	TreesURL         string      `json:"trees_url"`
	StatusesURL      string      `json:"statuses_url"`
	LanguagesURL     string      `json:"languages_url"`
	StargazersURL    string      `json:"stargazers_url"`
	ContributorsURL  string      `json:"contributors_url"`
	SubscribersURL   string      `json:"subscribers_url"`
	SubscriptionURL  string      `json:"subscription_url"`
	CommitsURL       string      `json:"commits_url"`
	GitCommitsURL    string      `json:"git_commits_url"`
	CommentsURL      string      `json:"comments_url"`
	IssueCommentURL  string      `json:"issue_comment_url"`
	ContentsURL      string      `json:"contents_url"`
	CompareURL       string      `json:"compare_url"`
	MergesURL        string      `json:"merges_url"`
	ArchiveURL       string      `json:"archive_url"`
	DownloadsURL     string      `json:"downloads_url"`
	IssuesURL        string      `json:"issues_url"`
	PullsURL         string      `json:"pulls_url"`
	MilestonesURL    string      `json:"milestones_url"`
	NotificationsURL string      `json:"notifications_url"`
	LabelsURL        string      `json:"labels_url"`
	ReleasesURL      string      `json:"releases_url"`
	DeploymentsURL   string      `json:"deployments_url"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	PushedAt         time.Time   `json:"pushed_at"`
	GitURL           string      `json:"git_url"`
	SSHURL           string      `json:"ssh_url"`
	CloneURL         string      `json:"clone_url"`
	SvnURL           string      `json:"svn_url"`
	Homepage         string      `json:"homepage"`
	Size             int         `json:"size"`
	StargazersCount  int         `json:"stargazers_count"`
	WatchersCount    int         `json:"watchers_count"`
	Language         string      `json:"language"`
	HasIssues        bool        `json:"has_issues"`
	HasProjects      bool        `json:"has_projects"`
	HasDownloads     bool        `json:"has_downloads"`
	HasWiki          bool        `json:"has_wiki"`
	HasPages         bool        `json:"has_pages"`
	ForksCount       int         `json:"forks_count"`
	MirrorURL        interface{} `json:"mirror_url"`
	Archived         bool        `json:"archived"`
	Disabled         bool        `json:"disabled"`
	OpenIssuesCount  int         `json:"open_issues_count"`
	License          struct {
		Key    string `json:"key"`
		Name   string `json:"name"`
		SpdxID string `json:"spdx_id"`
		URL    string `json:"url"`
		NodeID string `json:"node_id"`
	} `json:"license"`
	Forks         int    `json:"forks"`
	OpenIssues    int    `json:"open_issues"`
	Watchers      int    `json:"watchers"`
	DefaultBranch string `json:"default_branch"`
	Permissions   struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	} `json:"permissions"`
}

type MarkDownRepo struct {
	FullName string
	HtmlUrl  string
	Language string
}

type MarkDownRow struct {
	Language string
	Stars    string
	Items    string
}

func main() {
	if len(token) == 0 {
		// check for missing github token
		fmt.Println("Missing Github token")
		return
	}
	totalPageCount := GetUserStarredRepositoriesTotalPage()
	starredRepositories := GetUserAllStarredRepositories(totalPageCount)
	repos := GroupByProgrammingLanguage(starredRepositories)
	slices := Covert2Slice(repos)
	Print2File(slices)
}

func Print2File(markDownRows []MarkDownRow) error {
	os.Remove("./out.md")
	output, _ := os.Create("./out.md")
	defer output.Close()
	tpl := template.Must(template.ParseFiles("./template/starred.md"))
	return Print2Template(output, tpl, markDownRows)
}

func Print2Template(
	wr io.Writer,
	tpl *template.Template,
	markDownRows []MarkDownRow,
) error {
	return tpl.ExecuteTemplate(wr, "layout", markDownRows)
}

func Covert2Slice(repos map[string][]MarkDownRepo) []MarkDownRow {
	keys := GetMapKeyASC(repos)
	var rows = make([]MarkDownRow, 0, len(keys))
	for _, v := range keys {
		row := MarkDownRow{
			Language: v,
			Stars:    strconv.Itoa(len(repos[v])),
			Items:    GetInnerReposStr(repos[v]),
		}
		rows = append(rows, row)
	}
	return rows
}

func GetUserStarredRepositoriesTotalPage() (totalPage int) {
	rawURI := fmt.Sprintf(GithubURI, GitHubUserName)
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
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("%s", err)
	}
	linkHeader := resp.Header.Get("link")
	totalPage = ParseRawLinkHeader(linkHeader)
	return
}

func GetUserAllStarredRepositories(totalPage int) (userStarredRepositories UserStarredRepositories) {
	rawURI := fmt.Sprintf(GithubURI, GitHubUserName)
	u, err := url.Parse(rawURI)
	if err != nil {
		log.Print(err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()
	ch := make(chan UserStarredRepositories, totalPage)
	for i := 1; i <= totalPage; i++ {
		go func(pageNum int) {
			query := url.Values{
				"per_page": []string{"100"},
				"page":     []string{strconv.Itoa(pageNum)},
			}
			u.RawQuery = query.Encode()
			req, _ := http.NewRequest(http.MethodGet, u.String(), nil)
			req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
			req.Header.Set("Accept", "application/vnd.github.v3+json")
			resp, err := client.Do(req)
			if err != nil {
				log.Fatalf("%s", err)
			}
			defer resp.Body.Close()
			var singleUserStarredRepositoriesResponse UserStarredRepositories
			json.NewDecoder(resp.Body).Decode(&singleUserStarredRepositoriesResponse)
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

func ParseRawLinkHeader(rawHeader string) (totalPage int) {
	// rawHeader `<https://api.github.com/user/5622516/starred?per_page=100&page=2>; rel="next", <https://api.github.com/user/5622516/starred?per_page=100&page=18>; rel="last"`
	var links = strings.Split(rawHeader, ",")
	var regex = regexp.MustCompile(`\<(\S+)\>`)
	// only get matched group one for last links
	var lastPageRawURI = regex.FindStringSubmatch(links[1])[1]
	lastPageURI, err := url.Parse(lastPageRawURI)
	if err != nil {
		log.Println(err.Error())
	}
	totalPage, err = strconv.Atoi(lastPageURI.Query().Get("page"))
	if err != nil {
		log.Println(err.Error())
	}
	return
}

func GroupByProgrammingLanguage(userStarredRepositories UserStarredRepositories) map[string][]MarkDownRepo {
	var repos = make(map[string][]MarkDownRepo)
	for _, v := range userStarredRepositories {
		var languageKey = v.Language
		if v.Language == "" {
			// handle pure document repo
			languageKey = "other"
		}
		repos[languageKey] = append(
			repos[languageKey],
			MarkDownRepo{
				FullName: v.FullName,
				HtmlUrl:  v.HTMLURL,
				Language: v.Language,
			},
		)

	}
	return repos
}

func GetMapKeyASC(m map[string][]MarkDownRepo) (s []string) {
	for k := range m {
		s = append(s, k)
	}
	sort.Strings(s)
	return
}

func GetInnerReposStr(markDownRepo []MarkDownRepo) (s string) {
	var repos []string
	for _, v := range markDownRepo {
		repos = append(repos, fmt.Sprintf(MarkdownStar, v.FullName, v.HtmlUrl))
	}
	return strings.Join(repos[:], ", ")
}
