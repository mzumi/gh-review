package main

import (
	"fmt"
	"os"
)

func main() {
	repositories, err := RepositoryListByReviewRequest()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	view := NewView(repositories)
	view.show()
}
