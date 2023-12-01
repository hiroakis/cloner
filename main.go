package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	token := flag.String("token", "", "the token of github")
	org := flag.String("org", "", "the name of organization")
	typeOfRepo := flag.String("type", "private", "private or public. default: private")
	page := flag.Int("page", 1, "the page num. default: 1")
	perPage := flag.Int("per", 100, "the number of results to include per page. default: 100")
	flag.Parse()

	if *org == "" {
		// help message for user
		flag.Usage()
		os.Exit(1)
	}

	if *token == "" {
		// if you can get output, you will set token
		fmt.Println("if you can gh command, cloner will be set your token using 'gh auth token' command.")
		cmd := exec.Command("gh", "auth", "token")
		cmd.Stderr = os.Stderr
		output, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
			fmt.Println("-org are required. And you can not get token using gh command")
			fmt.Println("Please set -token using github.com/settings/tokens")
			os.Exit(1)
		}
		*token = strings.TrimSpace(string(output))
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)
	opt := &github.RepositoryListByOrgOptions{
		Type:        *typeOfRepo,
		ListOptions: github.ListOptions{Page: *page, PerPage: *perPage},
	}

	fmt.Println("org:", *org)
	fmt.Println("repository type:", *typeOfRepo)

	repos, response, err := client.Repositories.ListByOrg(context.Background(), *org, opt)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if response.StatusCode != 200 {
		fmt.Println("status code is not 200")
		fmt.Println(response.Status, response.Body, response.Rate)
		os.Exit(1)
	}

	fmt.Println("repositories:", len(repos))
	if len(repos) == 0 {
		fmt.Printf("%s repository is not found\n", *typeOfRepo)
		os.Exit(0)
	}

	// you will get get repository list from github
	for _, repo := range repos {
		fmt.Println(*repo.SSHURL)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	ch := make(chan bool)
	for _, repo := range repos {
		go func(repo github.Repository) {
			cmd := exec.Command("ghq", "get", "-p", *repo.SSHURL)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			fmt.Println(cmd.String())
			err = cmd.Run()
			if err != nil {
				fmt.Println(err)
				ch <- false
				return
			}
			ch <- true
		}(*repo)
	}

	for i := 0; i < len(repos); i++ {
		<-ch
	}
}
