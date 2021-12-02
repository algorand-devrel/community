package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

// TODO: check only open issues
// TODO: elipsis for title after N char in title?
// TODO: count comments
// TODO: unique contribs

type Issue struct {
	text  string
	tally int
}

type Repo struct {
	Org    string   `json:"org"`
	Name   string   `json:"name"`
	Labels []string `json:"labels"`
}

const header = `
Community Interest Tracker
----------------------

Below is the title and reaction count for issues filed in all Algorand-owned repositories. Issues labeled with the "community interest" label will appear on this list. Anybody is welcome to add this label to an issue, and express their preference or disagreement with an issue by reacting with :+1: or :-1:.

This tracker is used to gauge community interest in different features or improvements. A high interest count (many :+1:) does not guarantee the issue will be worked on since there are many other factors that go into that decision, but these counts serve as input to prioritization decisions.
`

const template = `
## %s
| Title | :+1: | :-1: |
| ----- | -- | ---- |
`

func main() {

	ctx := context.Background()

	var tc *http.Client

	if f, err := os.Open(".token"); err == nil {
		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatalf("Failed to read token file: %+v", err)
		}

		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: string(b)})
		tc = oauth2.NewClient(ctx, ts)
	}

	client := github.NewClient(tc)

	f, err := os.Create("README.md")
	if err != nil {
		log.Fatalf("Failed to create README file: %+v", err)
	}
	f.Write([]byte(header))

	repoList := getRepoList()
	for _, repo := range repoList {

		var (
			issueList []Issue
			opts      = &github.IssueListByRepoOptions{State: "open", Labels: repo.Labels}
			finished  = false
		)

		for !finished {
			issues, resp, err := client.Issues.ListByRepo(ctx, repo.Org, repo.Name, opts)
			if err != nil {
				log.Fatalf("Failed to get issues: %+v", err)
			}

			if resp.NextPage == 0 {
				finished = true
			}

			for _, issue := range issues {
				line := fmt.Sprintf("| [%s](%s) |", *issue.Title, *issue.HTMLURL)
				line += fmt.Sprintf(" %d |", *issue.Reactions.PlusOne)
				line += fmt.Sprintf(" %d |", *issue.Reactions.MinusOne)

				issueList = append(issueList, Issue{text: line, tally: *issue.Reactions.PlusOne - *issue.Reactions.MinusOne})
			}

			opts.Page = resp.NextPage
		}

		// Sort by tally
		sort.SliceStable(issueList, func(i, j int) bool {
			return issueList[i].tally > issueList[j].tally
		})

		fmt.Fprintf(f, "%s", fmt.Sprintf(template, repo.Name))
		for _, issue := range issueList {
			fmt.Fprintf(f, "%s\n", issue.text)
		}
	}

	f.Close()
}

func getRepoList() []Repo {
	f, err := os.Open("repos.json")
	if err != nil {
		log.Fatalf("Failed to read repos file: %+v", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("Failed to read bytes from repos file: %+v", err)
	}

	r := &struct {
		Repos []Repo `json:"repos"`
	}{}

	if err := json.Unmarshal(b, r); err != nil {
		log.Fatalf("Failed to unmarshal json: %+v", err)
	}

	return r.Repos
}
