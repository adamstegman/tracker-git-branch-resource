package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/xoebus/go-tracker"

	"github.com/adamstegman/tracker-git-branch-resource"
	"github.com/adamstegman/tracker-git-branch-resource/check"
)

func main() {
	var request check.Request
	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse input: %s", err)
		os.Exit(1)
	}

	repository := resource.NewRepository(request.Source.Repo)
	err = repository.Clone()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not clone repo %s: %s", request.Source.Repo, err)
		os.Exit(1)
	}

	if request.Source.TrackerURL != "" {
		tracker.DefaultURL = request.Source.TrackerURL
	}
	trackerToken := request.Source.Token
	trackerProjectID, err := strconv.Atoi(request.Source.ProjectID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid Tracker project ID %s: %s", request.Source.ProjectID, err)
		os.Exit(1)
	}
	projectClient := tracker.NewClient(trackerToken).InProject(trackerProjectID)

	finishedQuery := tracker.StoriesQuery{State: tracker.StoryStateFinished}
	finishedStories, err := projectClient.Stories(finishedQuery)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not fetch finished stories: %s", err)
		os.Exit(1)
	}
	deliveredQuery := tracker.StoriesQuery{State: tracker.StoryStateDelivered}
	deliveredStories, err := projectClient.Stories(deliveredQuery)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not fetch delivered stories: %s", err)
		os.Exit(1)
	}
	stories := append(finishedStories, deliveredStories...)

	trackerGitBranchCheck := check.NewTrackerGitBranchCheck(request.Version, repository, stories)
	versions, err := trackerGitBranchCheck.NewVersions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not find versions: %s", err)
		os.Exit(1)
	}

	err = json.NewEncoder(os.Stdout).Encode(versions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not print response: %s", err)
		os.Exit(1)
	}
}
