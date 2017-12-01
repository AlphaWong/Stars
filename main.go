package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const (
	Others         = "Others"
	MarkdownHeader = "%s|⭐️|%s\n"
	MarkdownColumn = "%s|%s|%s\n"
	MarkdownStar   = "[ [%s](%s) ]"
)

type Repos struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

func main() {
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
		response, err := http.Get("https://api.github.com/users/" + "alphawong" + "/starred?page=" + strconv.Itoa(i))
		if err != nil {
			log.Fatalf("%s", err)
		} else {
			defer response.Body.Close()

			contents, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatalf("%s", err)
			}

			json.NewDecoder(bytes.NewReader(contents)).Decode(&t)
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
