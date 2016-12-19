package main

import (
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	token := flag.String("token", "", "the token of github")
	org := flag.String("org", "", "the name of organization")
	typeOfRepo := flag.String("type", "private", "private or public. default: private")
	flag.Parse()

	if *token == "" || *org == "" {
		fmt.Println("-token and -org are requered")
		os.Exit(1)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{Type: *typeOfRepo}
	repos, _, err := client.Repositories.ListByOrg(*org, opt)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	ch := make(chan bool)
	for _, repo := range repos {
		go func(repo github.Repository) {
			_, err := os.Stat(*repo.Name)
			if err == nil {
				ch <- false
				return
			}

			progress := fmt.Sprintf("git clone %s", *repo.SSHURL)
			fmt.Println(progress)

			cmd := exec.Command("git", "clone", *repo.SSHURL)
			cmd.Start()
			ch <- true
		}(*repo)
	}

	for i := 0; i < len(repos); i++ {
		<-ch
	}
}
