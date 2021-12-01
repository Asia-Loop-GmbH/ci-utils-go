package gitcommand

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

type gitRepo struct {
	URL   string
	Owner string
	Repo  string
}

func getGitRepo() *gitRepo {
	output, err := exec.Command("git", "remote", "-v").CombinedOutput()
	if err != nil {
		panic(string(output))
	}

	r := regexp.MustCompile("https://(.*)/(.*)/(.*).git")
	matches := r.FindStringSubmatch(string(output))

	return &gitRepo{
		URL:   fmt.Sprintf("https://%s:%s@%s/%s/%s.git", os.Getenv("AL_GITHUB_USER"), os.Getenv("AL_GITHUB_TOKEN"), matches[1], matches[2], matches[3]),
		Owner: matches[2],
		Repo:  matches[3],
	}
}
