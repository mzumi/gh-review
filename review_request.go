package main

import (
	"context"
	"os"
	"sync"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	token     = os.Getenv("TOKEN")
	org       = os.Getenv("ORGANIZATION")
	userLogin = os.Getenv("USER")
)

type Repository struct {
	*github.Repository

	PullRequestList []*github.PullRequest
}

func RepositoryListByReviewRequest() ([]Repository, error) {
	ctx := context.Background()
	client := generateClient(ctx)

	repos, _, err := client.Repositories.ListByOrg(ctx, org, nil)
	if err != nil {
		return nil, err
	}

	c := make(chan Repository, len(repos))
	wg := &sync.WaitGroup{}
	wg.Add(len(repos))

	for _, repo := range repos {
		go fetchRepository(ctx, c, wg, client, repo)
	}

	wg.Wait()
	close(c)

	repository := []Repository{}
	for r := range c {
		if len(r.PullRequestList) > 0 {
			repository = append(repository, r)
		}
	}

	return repository, nil
}

func generateClient(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func fetchRepository(ctx context.Context, c chan Repository, wc *sync.WaitGroup, client *github.Client, repo *github.Repository) {
	defer wc.Done()

	repository := Repository{Repository: repo}

	prs, _, err := client.PullRequests.List(ctx, org, *repo.Name, &github.PullRequestListOptions{State: "open"})
	if err != nil {
		return
	}

	for _, pr := range prs {
		reviewers, _, err := client.PullRequests.ListReviewers(ctx, org, repo.GetName(), pr.GetNumber(), nil)
		if err != nil {
			continue
		}
		if p := filterByUser(pr, repo, reviewers); p != nil {
			repository.PullRequestList = append(repository.PullRequestList, p)
		}
	}

	c <- repository
}

func filterByUser(pr *github.PullRequest, repo *github.Repository, r *github.Reviewers) *github.PullRequest {
	for _, user := range r.Users {
		if *user.Login == userLogin {
			return pr
		}
	}
	return nil
}
