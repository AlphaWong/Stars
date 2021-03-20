package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/AlphaWong/Stars/services"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type BaseConfig struct {
	Token        string `validate:"required"`
	UserName     string `validate:"required"`
	BaseTemplate string `validate:"required"`
	OutputPath   string `validate:"required"`
	mu           sync.Mutex
}

func main() {
	config := boot()
	run(config)
}

func boot() (config *BaseConfig) {
	// Githun access token
	token := os.Getenv("TOKEN")
	// Github username
	userName := "alphawong"
	config = &BaseConfig{
		Token:        token,
		UserName:     userName,
		BaseTemplate: "./template/starred.md",
		OutputPath:   "./out.md",
	}
	// ensure the config is valid
	validConfig(config)
	return
}

func validConfig(config *BaseConfig) {
	validate = validator.New()
	err := validate.Struct(config)
	if nil != err {
		for _, err := range err.(validator.ValidationErrors) {
			log.Printf(err.Error())
		}
		log.Panicln("Invalid config struct")
	}
}

func run(config *BaseConfig) {
	fetcher, err := services.NewGitHubFetcher(
		services.WithToken(config.Token),
		services.WithUserName(config.UserName),
	)
	if nil != err {
		fmt.Print(err.Error())
		return
	}
	results := fetcher.GetUsersStars()

	baseTemplatePath, _ := filepath.Abs(config.BaseTemplate)
	outputPath, _ := filepath.Abs(config.OutputPath)
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
