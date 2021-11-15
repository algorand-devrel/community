package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	f, err := os.Create("REACTIONS.md")
	if err != nil {
		log.Fatalf("Failed to create README file: %+v", err)
	}

	fmt.Fprintf(f, "%s", template)
	for _, issue := range issues {
		fmt.Fprintf(f, "- *%s* ", *issue.Title)
		fmt.Fprintf(f, " :heavy_plus_sign: %d", *issue.Reactions.PlusOne)
		fmt.Fprintf(f, " :heavy_minus_sign: %d", *issue.Reactions.MinusOne)
		fmt.Fprintf(f, " :laughing: %d", *issue.Reactions.Laugh)
		fmt.Fprintf(f, " :confused: %d", *issue.Reactions.Confused)
		fmt.Fprintf(f, " :heart: %d", *issue.Reactions.Heart)
		fmt.Fprintf(f, " :tada: %d", *issue.Reactions.Hooray)
		fmt.Fprintf(f, " :rocket: %d", *issue.Reactions.Rocket)
		fmt.Fprintf(f, " :eyes: %d\n", *issue.Reactions.Eyes)
	}

	f.Close()
}
