package maven

import (
	"encoding/xml"
	"io/ioutil"
)

type POM struct {
	XMLName xml.Name `xml:"project"`
	Version string   `xml:"version"`
	Build   POMBuild `xml:"build"`
}

type POMBuild struct {
	FinalName string `xml:"finalName"`
}

func Read() *POM {
	filePath := "pom.xml"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	var pom POM
	xml.Unmarshal(data, &pom)
	return &pom
}
