package main

import (
	"fmt"
	"io"
	"os"
)

var writer io.Writer

func init() {
	writer = os.Stdout
}

type View struct {
	Repositories []Repository
}

func NewView(repos []Repository) *View {
	return &View{
		Repositories: repos,
	}
}

func (v *View) Show() {
	reviewCount := 0
	for _, repo := range v.Repositories {
		reviewCount += len(repo.PullRequestList)
	}

	fprintf("review (%d)\n", reviewCount)
	fprintf("---\n")

	for _, repo := range v.Repositories {
		fprintf("%s | href=%s\n", repo.GetName(), repo.GetHTMLURL())
		for _, pr := range repo.PullRequestList {
			fprintf("- #%d | href=%s\n", pr.GetNumber(), pr.GetHTMLURL())
		}
	}
}

func fprintf(format string, a ...interface{}) {
	fmt.Fprintf(writer, format, a...)
}
