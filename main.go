package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"

	"github.com/LukasDeco/github-portfolio-generator/netlify"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

var reposOfInterest []string = []string{
	"lzw-compress",
	"token-buddy",
	"slicefuncs",
	"lego",
	"twilio-rust",
	"github-portfolio-generator",
}

func main() {
	ctx := context.Background()
	err := createPortfolio(ctx, "LukasDeco", "index.html")
	if err != nil {
		panic(err)
	}
	fmt.Println("portfolio website generated!")
}

func deployVercel(ctx context.Context, html, website, token string) error {
	return nil
}

func createWebsite(ctx context.Context, style string, githubJson []byte) (string, error) {
	var html string
	prompt := fmt.Sprintf(`
	Make a beautiful website. Your code should be visually stunning, 
	but also intuitive and mobile friendly. Include some animations where appropriate. 
	Make the cards have a hover box-shadow animation.
	Use gradients and fun vector images for some of the elements and background.
	Remember to prioritize readability and organization, and to create a polished final product.
	Make any links open in a new tab.
	Give it a %s style please.
	Include the following information from this json in the website. 
	Just use the json struct to generate the html directly, don't use any JS.
	---
	%s
	`, style, githubJson)

	llm, err := openai.New()
	if err != nil {
		return "", err
	}

	aiRes, err := llm.Generate(ctx, []string{prompt},
		llms.WithTemperature(0.2),
		llms.WithMaxTokens(2500),
	)
	if err != nil {
		return "", err
	}
	if len(aiRes) == 0 {
		return "", errors.New("ai response empty")
	}

	html = aiRes[0].Text

	return html, nil
}

func createPortfolio(ctx context.Context, userName, portfolioHtmlFilename string) error {
	gitubApiUrl := fmt.Sprintf("https://api.github.com/users/%s", userName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, gitubApiUrl, nil)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	userJson, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	// fmt.Println(string(userJson))

	var user GithubUser
	err = json.Unmarshal(userJson, &user)
	if err != nil {
		return err
	}

	repoUrl := fmt.Sprintf("https://api.github.com/users/%s/repos", userName)
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, repoUrl, nil)
	if err != nil {
		return err
	}
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	repoJson, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var repos []GithubRepo

	err = json.Unmarshal(repoJson, &repos)
	if err != nil {
		return err
	}

	// fmt.Println(string(repoJson))

	var filteredRepos []GithubRepo
	for _, r := range repos {
		if slices.Contains(reposOfInterest, r.Name) {
			filteredRepos = append(filteredRepos, r)
		}
	}

	portfolio := GithubPortfolio{
		User:  user,
		Repos: filteredRepos,
	}

	portfolioJson, err := json.Marshal(portfolio)
	if err != nil {
		return err
	}

	htmlContents, err := createWebsite(ctx, "dark theme with blue oceans", portfolioJson)
	if err != nil {
		return err
	}

	err = os.WriteFile(portfolioHtmlFilename, []byte(htmlContents), 0644)
	if err != nil {
		return err
	}

	netlify.DeployNetlify(portfolioHtmlFilename)

	return nil
}
