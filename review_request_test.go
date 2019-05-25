package main

import (
	"context"
	"strconv"
	"testing"

	"github.com/google/go-github/v24/github"
)

func TestNoRepositories(t *testing.T) {
	ctx := context.Background()
	client := createNoRepositoryListMockGithubClient(ctx)
	repositories, _ := RepositoryListByReviewRequest(ctx, client)

	if len(repositories) != 0 {
		t.Errorf("Unexpected output: %v", len(repositories))
	}
}

func TestNoPullRequests(t *testing.T) {
	ctx := context.Background()
	client := createNoPullRequestListMockGithubClient(ctx)
	repositories, _ := RepositoryListByReviewRequest(ctx, client)

	if len(repositories) != 0 {
		t.Errorf("Unexpected output: %v", len(repositories))
	}
}

func TestNoAssignmentReviewee(t *testing.T) {
	repoName := "test repo"
	loginName := "test user"
	repositoryURL := "https://github.com/mzumi/gh-review"
	userLoginOld := userLogin
	defer func() { userLogin = userLoginOld }()

	userLogin = "test user2"
	prNumber := 1
	htmlURL := "https://github.com/mzumi/gh-review/pull/" + strconv.Itoa(prNumber)

	ctx := context.Background()
	client := createAssignmentPullRequestListMockGithubClient(ctx, loginName, prNumber, htmlURL, repoName, repositoryURL)
	repositories, _ := RepositoryListByReviewRequest(ctx, client)

	if len(repositories) != 0 {
		t.Errorf("Unexpected output: %v", len(repositories))
	}
}

func TestAssignmentReviewee(t *testing.T) {
	repoName := "test repo"
	loginName := "test user"
	repositoryURL := "https://github.com/mzumi/gh-review"

	userLoginOld := userLogin
	defer func() { userLogin = userLoginOld }()
	userLogin = loginName
	prNumber := 1
	htmlURL := repositoryURL + "/pull/" + strconv.Itoa(prNumber)

	ctx := context.Background()
	client := createAssignmentPullRequestListMockGithubClient(ctx, loginName, prNumber, htmlURL, repoName, repositoryURL)
	repositories, _ := RepositoryListByReviewRequest(ctx, client)

	if len(repositories) != 1 {
		t.Errorf("Unexpected output: %v", len(repositories))
	}

	repo := repositories[0]
	if repo.GetName() != repoName {
		t.Errorf("Unexpected repository name: %v", repo.GetName())
	}

	if repo.GetHTMLURL() != repositoryURL {
		t.Errorf("Unexpected repository name: %v", repo.GetHTMLURL())
	}

	pr := repo.PullRequestList[0]
	if pr.GetNumber() != prNumber {
		t.Errorf("Unexpected pull request number: %v", pr.GetNumber())
	}

	if pr.GetHTMLURL() != htmlURL {
		t.Errorf("Unexpected pull request url: %v", pr.GetHTMLURL())
	}
}

func createNoRepositoryListMockGithubClient(ctx context.Context) *mockGithubClient {
	return &mockGithubClient{
		MockRepositoryListByOrg: func(ctx context.Context, org string, opt *github.RepositoryListByOrgOptions) ([]*github.Repository, *github.Response, error) {
			var repos []*github.Repository
			res := &github.Response{}
			return repos, res, nil
		},
		MockPullRequestList: func(ctx context.Context, owner string, repo string, opt *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error) {
			return nil, nil, nil
		},
		MockPullRequestListReviewers: func(ctx context.Context, owner, repo string, number int, opt *github.ListOptions) (*github.Reviewers, *github.Response, error) {
			return nil, nil, nil
		},
	}
}

func createNoPullRequestListMockGithubClient(ctx context.Context) *mockGithubClient {
	repoName := "test repo"
	return &mockGithubClient{
		MockRepositoryListByOrg: func(ctx context.Context, org string, opt *github.RepositoryListByOrgOptions) ([]*github.Repository, *github.Response, error) {
			repos := []*github.Repository{&github.Repository{Name: &repoName}}
			res := &github.Response{}
			return repos, res, nil
		},
		MockPullRequestList: func(ctx context.Context, owner string, repo string, opt *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error) {
			var prs []*github.PullRequest
			return prs, nil, nil
		},
		MockPullRequestListReviewers: func(ctx context.Context, owner, repo string, number int, opt *github.ListOptions) (*github.Reviewers, *github.Response, error) {
			return nil, nil, nil
		},
	}
}

func createAssignmentPullRequestListMockGithubClient(ctx context.Context, loginName string, prNumber int, htmlURL string, repoName string, repositoryURL string) *mockGithubClient {
	return &mockGithubClient{
		MockRepositoryListByOrg: func(ctx context.Context, org string, opt *github.RepositoryListByOrgOptions) ([]*github.Repository, *github.Response, error) {
			repos := []*github.Repository{&github.Repository{Name: &repoName, HTMLURL: &repositoryURL}}
			res := &github.Response{}

			return repos, res, nil
		},
		MockPullRequestList: func(ctx context.Context, owner string, repo string, opt *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error) {
			pr := github.PullRequest{Number: &prNumber, HTMLURL: &htmlURL}

			pullRequests := []*github.PullRequest{&pr}

			return pullRequests, nil, nil
		},
		MockPullRequestListReviewers: func(ctx context.Context, owner, repo string, number int, opt *github.ListOptions) (*github.Reviewers, *github.Response, error) {
			reviewers := &github.Reviewers{Users: []*github.User{&github.User{Login: &loginName}}}
			return reviewers, nil, nil
		},
	}
}

type mockGithubClient struct {
	*github.Client
	MockRepositoryListByOrg      func(ctx context.Context, org string, opt *github.RepositoryListByOrgOptions) ([]*github.Repository, *github.Response, error)
	MockPullRequestList          func(ctx context.Context, owner string, repo string, opt *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
	MockPullRequestListReviewers func(ctx context.Context, owner, repo string, number int, opt *github.ListOptions) (*github.Reviewers, *github.Response, error)
}

func (client *mockGithubClient) RepositoryListByOrg(ctx context.Context, org string, opt *github.RepositoryListByOrgOptions) ([]*github.Repository, *github.Response, error) {
	return client.MockRepositoryListByOrg(ctx, org, opt)
}

func (client *mockGithubClient) PullRequestList(ctx context.Context, owner string, repo string, opt *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error) {
	return client.MockPullRequestList(ctx, owner, repo, opt)
}

func (client *mockGithubClient) PullRequestListReviewers(ctx context.Context, owner, repo string, number int, opt *github.ListOptions) (*github.Reviewers, *github.Response, error) {
	return client.MockPullRequestListReviewers(ctx, owner, repo, number, opt)
}
