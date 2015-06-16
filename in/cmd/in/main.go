package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adamstegman/tracker-git-branch-resource"
	"github.com/adamstegman/tracker-git-branch-resource/in"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <target directory>\n", os.Args[0])
		os.Exit(1)
	}
	targetDir := os.Args[1]

	var request in.InRequest
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse input: %s", err)
		os.Exit(1)
	}

	repository := resource.NewRepository(request.Source.Repo, targetDir)
	err := repository.Clone()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not clone repo %s: %s", request.Source.Repo, err)
		os.Exit(1)
	}
	err = repository.Fetch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not fetch repo %s: %s", request.Source.Repo, err)
		os.Exit(1)
	}
	err = repository.CheckoutRef(request.Version.Ref)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not checkout %s#%s: %s", request.Source.Repo, request.Version.Ref, err)
		os.Exit(1)
	}

	metadata := []resource.MetadataPair{}

	response := in.InResponse{Version: request.Version, Metadata: metadata}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		fmt.Fprintf(os.Stderr, "Could not print response: %s", err)
		os.Exit(1)
	}
}
