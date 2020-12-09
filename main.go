package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	Others         = "Others"
	MarkdownHeader = "%s|⭐️|%s\n"
	MarkdownColumn = "%s|%s|%s\n"
	MarkdownStar   = "[ [%s](%s) ]"

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

func main() {
	if len(token) == 0 {
		// check for missing github token
		fmt.Println("Missing Github token")
		return
	}
	fmt.Printf("%v", GetUserAllStarredRepositoriesTotalPage())
	// m := GetCustomerGithubStars()
	// PrintMarkdownHeader()
	// PrintMarkdownColumn()
	// PrintAsMarkdown(m)
}

func GetUserAllStarredRepositoriesTotalPage() (totalPage int) {
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
	// defer resp.Body.Close()
	// var userStarredRepositories UserStarredRepositories
	// json.NewDecoder(resp.Body).Decode(&userStarredRepositories)
	// return userStarredRepositories
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

// func PrintMarkdownHeader() {
// 	fmt.Printf(MarkdownHeader, "Language", "Repos")
// }

// func PrintMarkdownColumn() {
// 	fmt.Printf(MarkdownColumn, "---", "---", "---")
// }

// func PrintAsMarkdown(m map[string][]Repos) {
// 	sk := SortMapByKeyAlphabatAsc(m)
// 	sort.Strings(sk)
// 	for _, v := range sk {
// 		fmt.Printf(MarkdownColumn, v, strconv.Itoa(len(m[v])), GetInnerReposStr(m[v]))
// 	}
// }

// func SortMapByKeyAlphabatAsc(m map[string][]Repos) (s []string) {
// 	for k, _ := range m {
// 		s = append(s, k)
// 	}
// 	return
// }

// func GetInnerReposStr(r []Repos) (s string) {
// 	var ss []string
// 	for _, v := range r {
// 		ss = append(ss, fmt.Sprintf(MarkdownStar, v.Name, v.URI))
// 	}
// 	return strings.Join(ss[:], ", ")
// }

// func GetCustomerGithubStars() map[string][]Repos {
// 	m := make(map[string][]Repos)

// 	var t []map[string]interface{}

// 	i := 0
// 	for true {
// 		req, _ := http.NewRequest(http.MethodGet, GithubURI+strconv.Itoa(i), nil)
// 		req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
// 		response, err := client.Do(req)
// 		if err != nil {
// 			log.Fatalf("%s", err)
// 		} else {
// 			defer response.Body.Close()

// 			json.NewDecoder(response.Body).Decode(&t)

// 			if len(t) == 0 {
// 				return m
// 			}

// 			for _, v := range t {
// 				lang := Others
// 				if nil != v["language"] {
// 					lang = v["language"].(string)
// 				}
// 				r := Repos{
// 					Name: v["full_name"].(string),
// 					URI:  v["html_url"].(string),
// 				}
// 				if _, ok := m[lang]; !ok {
// 					m[lang] = []Repos{r}
// 				} else {
// 					m[lang] = append(m[lang], r)
// 				}
// 			}
// 		}
// 		i++
// 	}
// 	return m
// }
