package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/xoebus/go-tracker"

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
		fmt.Fprintf(os.Stderr, "Could not parse input: %s\n", err)
		os.Exit(1)
	}

	repository := resource.NewRepository(request.Source.Repo, targetDir, request.Source.PrivateKey)
	err := repository.Clone()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not clone repo %s: %s\n", request.Source.Repo, err)
		os.Exit(1)
	}
	err = repository.Fetch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not fetch repo %s: %s\n", request.Source.Repo, err)
		os.Exit(1)
	}
	err = repository.CheckoutRef(request.Version.Ref)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not checkout %s#%s: %s\n", request.Source.Repo, request.Version.Ref, err)
		os.Exit(1)
	}

	metadata, err := metadata(request, repository)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not fetch metadata for %s#%s: %s\n", request.Source.Repo, request.Version.Ref, err)
		os.Exit(1)
	}

	response := in.InResponse{Version: request.Version, Metadata: metadata}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		fmt.Fprintf(os.Stderr, "Could not print response: %s\n", err)
		os.Exit(1)
	}
}

func metadata(request in.InRequest, repository resource.Repository) ([]resource.MetadataPair, error) {
	authorName, err := repository.RefAuthorName(request.Version.Ref)
	if err != nil {
		return []resource.MetadataPair{}, err
	}
	authorDate, err := repository.RefAuthorDate(request.Version.Ref)
	if err != nil {
		return []resource.MetadataPair{}, err
	}
	committerName, err := repository.RefCommitName(request.Version.Ref)
	if err != nil {
		return []resource.MetadataPair{}, err
	}
	committerDate, err := repository.RefCommitDate(request.Version.Ref)
	if err != nil {
		return []resource.MetadataPair{}, err
	}
	message, err := repository.RefMessage(request.Version.Ref)
	if err != nil {
		return []resource.MetadataPair{}, err
	}
	trackerURL := request.Source.TrackerURL
	if trackerURL == "" {
		trackerURL = tracker.DefaultURL
	}
	storyURL := fmt.Sprintf("%s/story/show/%s", trackerURL, request.Version.StoryID)
	return []resource.MetadataPair{
		{Name: "commit", Value: request.Version.Ref},
		{Name: "author", Value: authorName},
		{Name: "author_date", Value: authorDate},
		{Name: "committer", Value: committerName},
		{Name: "committer_date", Value: committerDate},
		{Name: "message", Value: message},
		{Name: "story_url", Value: storyURL},
	}, nil
}
