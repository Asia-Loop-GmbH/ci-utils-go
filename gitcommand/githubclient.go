package gitcommand

import (
	"fmt"
	"github.com/google/go-github/v41/github"
	"net/http"
	"os"
)

func (trans *gitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("token %s", os.Getenv("AL_GITHUB_TOKEN")))
	return http.DefaultTransport.RoundTrip(req)
}

func getGithubClient() *github.Client {
	return github.NewClient(getGithubHttpClient())
}

func getGithubHttpClient() *http.Client {
	return &http.Client{
		Transport: &gitTransport{},
	}
}
