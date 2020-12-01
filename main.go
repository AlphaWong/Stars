package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	Others         = "Others"
	MarkdownHeader = "%s|⭐️|%s\n"
	MarkdownColumn = "%s|%s|%s\n"
	MarkdownStar   = "[ [%s](%s) ]"

	GithubURI = "https://api.github.com/users/alphawong/starred?page="
)

var (
	token  = os.Getenv("TOKEN")
	client = &http.Client{}
)

type Repos struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

func main() {
	if len(token) == 0 {
		// simple check for missing github token
		fmt.Println("Missing Github token")
		return
	}
	m := GetCustomerGithubStars()
	PrintMarkdownHeader()
	PrintMarkdownColumn()
	PrintAsMarkdown(m)
}

func PrintMarkdownHeader() {
	fmt.Printf(MarkdownHeader, "Language", "Repos")
}

func PrintMarkdownColumn() {
	fmt.Printf(MarkdownColumn, "---", "---", "---")
}

func PrintAsMarkdown(m map[string][]Repos) {
	sk := SortMapByKeyAlphabatAsc(m)
	sort.Strings(sk)
	for _, v := range sk {
		fmt.Printf(MarkdownColumn, v, strconv.Itoa(len(m[v])), GetInnerReposStr(m[v]))
	}
}

func SortMapByKeyAlphabatAsc(m map[string][]Repos) (s []string) {
	for k, _ := range m {
		s = append(s, k)
	}
	return
}

func GetInnerReposStr(r []Repos) (s string) {
	var ss []string
	for _, v := range r {
		ss = append(ss, fmt.Sprintf(MarkdownStar, v.Name, v.URI))
	}
	return strings.Join(ss[:], ", ")
}

func GetCustomerGithubStars() map[string][]Repos {
	m := make(map[string][]Repos)

	var t []map[string]interface{}

	i := 0
	for true {
		req, _ := http.NewRequest(http.MethodGet, GithubURI+strconv.Itoa(i), nil)
		req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
		response, err := client.Do(req)
		if err != nil {
			log.Fatalf("%s", err)
		} else {
			defer response.Body.Close()

			json.NewDecoder(response.Body).Decode(&t)

			if len(t) == 0 {
				return m
			}

			for _, v := range t {
				lang := Others
				if nil != v["language"] {
					lang = v["language"].(string)
				}
				r := Repos{
					Name: v["full_name"].(string),
					URI:  v["html_url"].(string),
				}
				if _, ok := m[lang]; !ok {
					m[lang] = []Repos{r}
				} else {
					m[lang] = append(m[lang], r)
				}
			}
		}
		i++
	}
	return m
}
