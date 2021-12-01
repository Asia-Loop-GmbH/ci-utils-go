package gitcommand

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/v41/github"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
)

var owner = flag.String("owner", "", "Github Organisation")
var repo = flag.String("repo", "", "Github repository")
var releaseId = flag.Int64("releaseId", 0, "Github release ID")
var assetName = flag.String("asset", "", "Github release asset name")

func parseParams() {
	err := flag.CommandLine.Parse(os.Args[3:])
	if err != nil {
		panic(err)
	}
}

func downloadRelease() {
	parseParams()
	client := getGithubClient()
	release, _, err := client.Repositories.GetRelease(context.TODO(), *owner, *repo, *releaseId)
	if err != nil {
		panic(err)
	}
	assetId := int64(0)
	for _, asset := range release.Assets {
		if asset.GetName() == *assetName {
			assetId = asset.GetID()
		}
	}

	if assetId == 0 {
		log.Panicf("could not find asset '%s' from release '%d'", *assetName, *releaseId)
	}

	download, _, err := client.Repositories.DownloadReleaseAsset(context.TODO(), *owner, *repo, assetId, getGithubHttpClient())
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := download.Close(); err != nil {
			log.Println("could not close connection to github")
			os.Exit(1)
		}
	}()

	content, err := ioutil.ReadAll(download)
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(*assetName, content, fs.ModePerm); err != nil {
		panic(err)
	}
}

func listReleases() {
	parseParams()

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
