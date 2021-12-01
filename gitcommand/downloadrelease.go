package gitcommand

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/v41/github"
	"os"
)

var owner = flag.String("owner", "", "Github Organisation")
var repo = flag.String("repo", "", "Github repository")

func listReleases() {
	err := flag.CommandLine.Parse(os.Args[3:])
	if err != nil {
		panic(err)
	}

	client := getGithubClient()
	releases, _, err := client.Repositories.ListReleases(context.TODO(), *owner, *repo, &github.ListOptions{
		Page:    1,
		PerPage: 10,
	})
	if err != nil {
		panic(err)
	}
	for _, r := range releases {
		fmt.Printf("%s|%d\n", *r.Name, *r.ID)
	}
}
