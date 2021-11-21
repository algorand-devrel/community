package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

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

func main() {
	ctx := context.Background()
	client := github.NewClient(nil)

	issues, _, err := client.Issues.ListByRepo(ctx, "algorand-devrel", "community", nil)
	if err != nil {
		log.Fatalf("Failed to get issues: %+v", err)
	}

	f, err := os.Create("README.md")
	if err != nil {
		log.Fatalf("Failed to create README file: %+v", err)
	}

	var lines []string
	for _, issue := range issues {
		line := fmt.Sprintf("| [%s](%s) |", *issue.Title, *issue.HTMLURL)
		line += fmt.Sprintf(" %d |", *issue.Reactions.PlusOne)
		line += fmt.Sprintf(" %d |", *issue.Reactions.MinusOne)
		lines = append(lines, line)
	}
	sort.Strings(lines)

	fmt.Fprintf(f, "%s%s", template, strings.Join(lines, "\n"))
	f.Close()
}
