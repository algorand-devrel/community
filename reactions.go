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

const template = `
Issue Reaction Counter
----------------------

Below is the title and reaction count for issues filed in this repository

Not sure what the reactions mean yet but I want to get a sense for how popular these issues are to help with prioritization

Feel free to file another issue and I'll re-run this script, maybe setup some github action to do it automatically?


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

	// TODO: check only open issues
	// TODO: sort by reactions?
	// TODO: elipsis for title after N char
	var lines []string
	for _, issue := range issues {
		line := fmt.Sprintf("| [%s](%s) |", *issue.Title, *issue.HTMLURL)
		line += fmt.Sprintf(" %d |", *issue.Reactions.PlusOne)
		line += fmt.Sprintf(" %d |", *issue.Reactions.MinusOne)
		//line += fmt.Sprintf(" :laughing: %d", *issue.Reactions.Laugh)
		//line += fmt.Sprintf(" :confused: %d", *issue.Reactions.Confused)
		//line += fmt.Sprintf(" :heart: %d", *issue.Reactions.Heart)
		//line += fmt.Sprintf(" :tada: %d", *issue.Reactions.Hooray)
		//line += fmt.Sprintf(" :rocket: %d", *issue.Reactions.Rocket)
		//line += fmt.Sprintf(" :eyes: %d", *issue.Reactions.Eyes)
		lines = append(lines, line)
	}
	sort.Strings(lines)

	fmt.Fprintf(f, "%s%s", template, strings.Join(lines, "\n"))
	f.Close()
}
