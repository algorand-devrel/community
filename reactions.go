package main

import (
	"context"
	"fmt"
	"log"

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

	fmt.Printf("%s", template)
	for _, issue := range issues {
		fmt.Printf("- *%s* ", *issue.Title)

		fmt.Printf(" :plusone: %d", *issue.Reactions.PlusOne)
		fmt.Printf(" :minusone: %d", *issue.Reactions.MinusOne)
		fmt.Printf(" :laugh: %d", *issue.Reactions.Laugh)
		fmt.Printf(" :confused: %d", *issue.Reactions.Confused)
		fmt.Printf(" :heart: %d", *issue.Reactions.Heart)
		fmt.Printf(" :hooray: %d", *issue.Reactions.Hooray)
		fmt.Printf(" :rocket: %d", *issue.Reactions.Rocket)
		fmt.Printf(" :eyes: %d\n", *issue.Reactions.Eyes)
	}
}
