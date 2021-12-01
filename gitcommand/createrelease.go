package gitcommand

import (
	"asialoop.de/ci-utils-go/maven"
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/v41/github"
	"log"
	"os"
	"os/exec"
	"time"
)

type gitTransport struct{}

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
	client := getGithubClient()
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
