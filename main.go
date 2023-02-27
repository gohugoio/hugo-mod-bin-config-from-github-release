package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/bep/githubreleasedownloader"
)

var (
	repo            string
	owner           string
	tag             string
	binName         string
	targetDirectory string
)

func main() {
	flag.StringVar(&repo, "repo", "", "The repository name")
	flag.StringVar(&owner, "owner", "", "The repository owner")
	flag.StringVar(&tag, "tag", "", "The tag to download")
	flag.StringVar(&binName, "bin-name", "", "The name of the binary")
	flag.StringVar(&targetDirectory, "target-directory", "", "The target directory")

	flag.Parse()

	if repo == "" || owner == "" || tag == "" || targetDirectory == "" || binName == "" {
		flag.PrintDefaults()
		return
	}
	client, err := githubreleasedownloader.New()
	if err != nil {
		log.Fatal(err)
	}
	release, err := client.GetRelease(owner, repo, tag)
	if err != nil {
		log.Fatal(err)
	}

	targetFilename := filepath.Join(targetDirectory, "hugo.json")

	f, err := os.Create(targetFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	m := module{
		Type: "bin",
		Params: map[string]interface{}{
			"bin":            binName,
			"github_release": release,
		},
	}

	// Encode m as JSON to f.
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(config{Module: m}); err != nil {
		log.Fatal(err)
	}

}

type config struct {
	Module module `json:"module"`
}

type module struct {
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params,omitempty"`
}
