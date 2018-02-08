package main

import (
	"fmt"
)

type View struct {
	Repositories []Repository
}

func NewView(repos []Repository) *View {
	return &View{
		Repositories: repos,
	}
}

func (v *View) show() {
	reviewCount := 0
	for _, repo := range v.Repositories {
		reviewCount += len(repo.PullRequestList)
	}

	fmt.Printf("review (%d)\n", reviewCount)
	fmt.Println("---")

	for _, repo := range v.Repositories {
		fmt.Printf("%s | href=%s\n", repo.GetName(), repo.GetHTMLURL())
		for _, pr := range repo.PullRequestList {
			fmt.Printf("- #%d | href=%s\n", pr.GetNumber(), pr.GetHTMLURL())
		}
	}
}
