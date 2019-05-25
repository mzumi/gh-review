package main

import (
	"context"

	"github.com/google/go-github/v25/github"
	"golang.org/x/oauth2"
)

type Client interface {
	RepositoryListByOrg(ctx context.Context, org string, opt *github.RepositoryListByOrgOptions) ([]*github.Repository, *github.Response, error)
	PullRequestList(ctx context.Context, owner string, repo string, opt *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
	PullRequestListReviewers(ctx context.Context, owner, repo string, number int, opt *github.ListOptions) (*github.Reviewers, *github.Response, error)
}

type GithubClient struct {
	*github.Client
}

func NewGithubClient(ctx context.Context) *GithubClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	return &GithubClient{
		Client: github.NewClient(tc),
	}
}

func (client *GithubClient) RepositoryListByOrg(ctx context.Context, org string, opt *github.RepositoryListByOrgOptions) ([]*github.Repository, *github.Response, error) {
	return client.Repositories.ListByOrg(ctx, org, opt)
}

func (client *GithubClient) PullRequestList(ctx context.Context, owner string, repo string, opt *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error) {
	return client.PullRequests.List(ctx, org, repo, &github.PullRequestListOptions{State: "open"})
}

func (client *GithubClient) PullRequestListReviewers(ctx context.Context, owner, repo string, number int, opt *github.ListOptions) (*github.Reviewers, *github.Response, error) {
	return client.PullRequests.ListReviewers(ctx, owner, repo, number, opt)
}
