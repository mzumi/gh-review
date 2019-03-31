package main

import (
	"context"
	"os"
	"sync"

	"github.com/google/go-github/v24/github"
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

func RepositoryListByReviewRequest(ctx context.Context, client Client) ([]Repository, error) {
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

func generateRepoChan(ctx context.Context, client Client) (<-chan *github.Repository, error) {
	queue := make(chan *github.Repository, 10)
	go func() {
		defer close(queue)
		opt := &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{PerPage: 10}}

		for {
			repos, resp, err := client.RepositoryListByOrg(ctx, org, opt)
			if err != nil {
				return
			}

			for _, repo := range repos {
				queue <- repo
			}

			if resp.NextPage == 0 {
				return
			}
			opt.Page = resp.NextPage
		}
	}()

	return queue, nil
}

func generatePullRequestChan(ctx context.Context, repoChan <-chan *github.Repository, client Client) <-chan Repository {
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

func fetchRepository(ctx context.Context, repo *github.Repository, client Client) *Repository {
	repository := &Repository{Repository: repo}

	prs, _, err := client.PullRequestList(ctx, org, *repo.Name, &github.PullRequestListOptions{State: "open"})
	if err != nil {
		return nil
	}

	for _, pr := range prs {
		reviewers, _, err := client.PullRequestListReviewers(ctx, org, repo.GetName(), pr.GetNumber(), nil)
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
