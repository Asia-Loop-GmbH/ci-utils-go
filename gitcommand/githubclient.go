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
	httpClient := http.Client{
		Transport: &gitTransport{},
	}
	return github.NewClient(&httpClient)
}
