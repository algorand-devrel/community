package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/google/go-github/v40/github" // with go modules enabled (GO111MODULE=on or outside GOPATH)
)

const template = `
Issue Reaction Counter
----------------------


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
		line := fmt.Sprintf("- *%s* ", *issue.Title)
		line += fmt.Sprintf(" :heavy_plus_sign: %d", *issue.Reactions.PlusOne)
		line += fmt.Sprintf(" :heavy_minus_sign: %d", *issue.Reactions.MinusOne)
		line += fmt.Sprintf(" :laughing: %d", *issue.Reactions.Laugh)
		line += fmt.Sprintf(" :confused: %d", *issue.Reactions.Confused)
		line += fmt.Sprintf(" :heart: %d", *issue.Reactions.Heart)
		line += fmt.Sprintf(" :tada: %d", *issue.Reactions.Hooray)
		line += fmt.Sprintf(" :rocket: %d", *issue.Reactions.Rocket)
		line += fmt.Sprintf(" :eyes: %d", *issue.Reactions.Eyes)
		lines = append(lines, line)
	}
	sort.Strings(lines)

	fmt.Fprintf(f, "%s %s", template, strings.Join(lines, "\n"))
	f.Close()
}
