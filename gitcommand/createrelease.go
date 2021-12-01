package gitcommand

import (
	"asialoop.de/ci-utils-go/maven"
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/v41/github"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"time"
)

type gitTransport struct{}

func (trans *gitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("token %s", os.Getenv("AL_GITHUB_TOKEN")))
	return http.DefaultTransport.RoundTrip(req)
}

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

func tagDate() string {
	return time.Now().Format("20060102T150405")
}

// currently support only maven
func createRelease() {
	log.Println("read version and artifact information from pom.xml")
	flag.Parse()
	pom := maven.Read()
	tag := fmt.Sprintf("v%s-%s", pom.Version, tagDate())
	artifactName := fmt.Sprintf("%s.jar", pom.Build.FinalName)
	artifact := fmt.Sprintf("target/%s", artifactName)
	log.Printf("version '%s' will be created with artifact '%s'", tag, artifact)

	repo := getGitRepo()

	log.Printf("create git tag '%s'", tag)
	out, err := exec.Command("git", "tag", tag).CombinedOutput()
	if err != nil {
		panic(string(out))
	}
	out, err = exec.Command("git", "push", repo.URL, tag).CombinedOutput()
	if err != nil {
		panic(string(out))
	}
	log.Printf("git tag '%s' created and pushed to origin", tag)

	log.Printf("create release '%s'", tag)
	httpClient := http.Client{
		Transport: &gitTransport{},
	}
	client := github.NewClient(&httpClient)
	release, _, err := client.Repositories.CreateRelease(context.TODO(), repo.Owner, repo.Repo, &github.RepositoryRelease{
		TagName: &tag,
		Name:    &tag,
	})
	if err != nil {
		panic(err)
	}
	log.Printf("release '%s' created at %s", tag, *release.HTMLURL)

	log.Printf("upload '%s'", artifact)
	file, err := os.Open(artifact)
	if err != nil {
		panic(err)
	}
	_, _, err = client.Repositories.UploadReleaseAsset(context.TODO(), repo.Owner, repo.Repo, *release.ID, &github.UploadOptions{
		Name: artifactName,
	}, file)
	if err != nil {
		panic(err)
	}
	log.Printf("artifact uploaded")
}
