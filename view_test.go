package main

import (
	"bytes"
	"testing"

	"github.com/google/go-github/v24/github"
)

var buffer *bytes.Buffer

func init() {
	buffer = &bytes.Buffer{}
	writer = buffer
}

func TestNoReviewee(t *testing.T) {
	buffer.Reset()

	repositories := []Repository{}
	view := NewView(repositories)
	view.Show()

	if buffer.String() != "review (0)\n---" {
		t.Errorf("Unexpected output: %s", buffer.String())
	}
}

func TestAssignReviewee(t *testing.T) {
	buffer.Reset()

	prNumber := 11
	prURL := "https://github.com/mzumi/gh-review/pull/11"

	pr := github.PullRequest{Number: &prNumber, HTMLURL: &prURL}

	pullRequests := []*github.PullRequest{}
	pullRequests = append(pullRequests, &pr)

	repositoryName := "gh-review"
	repositoryURL := "https://github.com/mzumi/gh-review"

	repo := Repository{
		Repository:      &github.Repository{Name: &repositoryName, HTMLURL: &repositoryURL},
		PullRequestList: pullRequests,
	}

	repositories := []Repository{}
	repositories = append(repositories, repo)

	view := NewView(repositories)
	view.Show()

	expected := "review (1)\n---gh-review | href=https://github.com/mzumi/gh-review\n- #11 | href=https://github.com/mzumi/gh-review/pull/11\n"

	if buffer.String() != expected {
		t.Errorf("Unexpected output: %s", buffer.String())
	}
}
