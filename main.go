package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/AlphaWong/Stars/services"
)

const (
	Others       = "Others"
	MarkdownStar = "[ [%s](%s) ]"

	// GithubURI store the base uri
	// "https://api.github.com/users/alphawong/starred"
	GithubURI = "https://api.github.com/users/%s/starred"
)

var (
	// Githun access token
	token = os.Getenv("TOKEN")
	// Github username
	userName = "alphawong"
)

func main() {
	if len(token) == 0 {
		// check for missing github token
		fmt.Println("Missing Github token")
		return
	}
	fetcher, err := services.NewGitHubFetcher(
		services.WithToken(token),
		services.WithUserName(userName),
	)
	if nil != err {
		fmt.Print(err.Error())
		return
	}
	results := fetcher.GetUsersStars()

	baseTemplatePath, _ := filepath.Abs("./template/starred.md")
	outputPath, _ := filepath.Abs("./out.md")
	printer, err := services.NewTplPrinter(
		services.WithBaseTemplate(template.ParseFiles(baseTemplatePath)),
		services.WithOutputPath(outputPath),
	)
	if nil != err {
		fmt.Print(err.Error())
		return
	}

	printer.PrintSlice(results)
}
