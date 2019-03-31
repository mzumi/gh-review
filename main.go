package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	ctx := context.Background()
	client := NewGithubClient(ctx)
	repositories, err := RepositoryListByReviewRequest(ctx, client)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	view := NewView(repositories)
	view.Show()
}
