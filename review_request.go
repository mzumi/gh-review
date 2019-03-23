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

const workerSize = 4

type Repository struct {
	*github.Repository

	PullRequestList []*github.PullRequest
}

func RepositoryListByReviewRequest() ([]Repository, error) {
	ctx := context.Background()
	client := generateClient(ctx)

	repoChan, err := generateRepoChan(ctx, client)
	if err != nil {
		return nil, err
	}
	prChan := generatePullRequestChan(ctx, repoChan, client)

	repository := []Repository{}
	for r := range prChan {
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

func generateRepoChan(ctx context.Context, client *github.Client) (<-chan *github.Repository, error) {
	repos, _, err := client.Repositories.ListByOrg(ctx, org, nil)
	if err != nil {
		return nil, err
	}

	queue := make(chan *github.Repository, len(repos))

	go func() {
		defer close(queue)

		for _, repo := range repos {
			queue <- repo
		}
	}()

	return queue, nil
}

func generatePullRequestChan(ctx context.Context, repoChan <-chan *github.Repository, client *github.Client) <-chan Repository {
	c := make(chan Repository, workerSize)
	wg := &sync.WaitGroup{}

	go func() {
		defer close(c)
		for i := 0; i < workerSize; i++ {
			wg.Add(1)
			go func() {
				for {
					repo, ok := <-repoChan
					if !ok {
						wg.Done()
						return
					}

					r := fetchRepository(ctx, repo, client)
					if r != nil {
						c <- *r
					}
				}
			}()
		}
		wg.Wait()
	}()

	return c
}

func fetchRepository(ctx context.Context, repo *github.Repository, client *github.Client) *Repository {
	repository := &Repository{Repository: repo}

	prs, _, err := client.PullRequests.List(ctx, org, *repo.Name, &github.PullRequestListOptions{State: "open"})
	if err != nil {
		return nil
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
	return repository
}

func filterByUser(pr *github.PullRequest, repo *github.Repository, r *github.Reviewers) *github.PullRequest {
	for _, user := range r.Users {
		if *user.Login == userLogin {
			return pr
		}
	}
	return nil
}
