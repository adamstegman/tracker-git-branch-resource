package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adamstegman/tracker-git-branch-resource"
	"github.com/adamstegman/tracker-git-branch-resource/in"
)

func main() {
	var request in.InRequest
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		fatal("reading request from stdin", err)
	}

	response := in.InResponse{
		Version: resource.Version{
			StoryID: request.Version.StoryID,
		},
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		fatal("writing response", err)
	}
}

func fatal(message string, err error) {
	fmt.Fprintf(os.Stderr, "error %s: %s\n", message, err)
	os.Exit(1)
}
