package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var in = flag.String("in", "Godeps.json", "The Godeps.json file that you want to convert")

// Godeps describes what a package needs to be rebuilt reproducibly.
// It's the same information stored in file Godeps.
type Godeps struct {
	ImportPath string
	Deps       []Dependency
}

type Dependency struct {
	ImportPath string
	Comment    string `json:",omitempty"` // Description of commit, if present.
	Rev        string // VCS-specific commit ID.
}

func main() {
	flag.Parse()

	godep, err := os.Open(*in)
	if err != nil {
		log.Fatalf("error opening file: %s", err)
	}

	var parsed Godeps
	if err := json.NewDecoder(godep).Decode(&parsed); err != nil {
		log.Fatalf("error parsing json: %s", err)
	}

	template := "[[override]]\n  name = \"%s\"\n  version = \"%s\"\n\n"

	deps := make(map[string]string)
	for _, dep := range parsed.Deps {
		splitted := strings.SplitN(dep.ImportPath, "/", 4)
		repo := ""
		if len(splitted) >= 3 {
			repo = fmt.Sprintf("%s/%s/%s", splitted[0], splitted[1], splitted[2])
		} else {
			repo = fmt.Sprintf("%s/%s", splitted[0], splitted[1])
		}

		deps[repo] = dep.Rev
	}

	for repo, version := range deps {
		fmt.Printf(template, repo, version)
	}
}
