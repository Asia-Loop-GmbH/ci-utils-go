package main

import (
	"os"
	"testing"
)

func TestDownloadReleaseArtifact(t *testing.T) {

	os.Args = []string{
		"ci-utils-go",
		"git",
		"download-release",
		"-owner=Asia-Loop-GmbH",
		"-repo=asia-loop",
		"-releaseId=54474200",
		"-asset=asia-loop-admin.jar",
	}

	main()
}
