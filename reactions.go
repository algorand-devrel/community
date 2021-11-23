package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/google/go-github/v40/github"
)

// TODO: check only open issues
// TODO: sort by reactions?
// TODO: elipsis for title after N char
// TODO: more repos
// TODO: count comments
// TODO: unique contribs

const template = `
Community Interest Tracker
----------------------

Below is the title and reaction count for issues filed in all Algorand-owned repositories. Issues labeled with the "community interest" label will appear on this list. Anybody is welcome to add this label to an issue, and express their preference or disagreement with an issue by reacting with :+1: or :-1:.

This tracker is used to gauge community interest in different features or improvements. A high interest count (many :+1:) does not guarantee the issue will be worked on since there are many other factors that go into that decision, but these counts serve as input to prioritization decisions.


| Title | :+1: | :-1: |
| ----- | -- | ---- |
`

type Issue struct {
	text  string
	tally int
}

func main() {
	ctx := context.Background()
	client := github.NewClient(nil)

	issues, _, err := client.Issues.ListByRepo(ctx, "algorand-devrel", "community", nil)
	if err != nil {
		log.Fatalf("Failed to get issues: %+v", err)
	}

	var issueList []Issue
	for _, issue := range issues {
		line := fmt.Sprintf("| [%s](%s) |", *issue.Title, *issue.HTMLURL)
		line += fmt.Sprintf(" %d |", *issue.Reactions.PlusOne)
		line += fmt.Sprintf(" %d |", *issue.Reactions.MinusOne)

		issueList = append(issueList, Issue{text: line, tally: *issue.Reactions.PlusOne - *issue.Reactions.MinusOne})
	}

	// Sort by tally
	sort.SliceStable(issueList, func(i, j int) bool {
		return issueList[i].tally > issueList[j].tally
	})

	f, err := os.Create("README.md")
	if err != nil {
		log.Fatalf("Failed to create README file: %+v", err)
	}

	fmt.Fprintf(f, "%s", template)
	for _, issue := range issueList {
		fmt.Fprintf(f, "%s\n", issue.text)
	}

	f.Close()
}
