package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/adamstegman/go-tracker"

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

	trackerCheck := check.NewTrackerGitBranchCheck(projectClient)
	stories, err := trackerCheck.StoriesFinishedAfterStory(request.Version.StoryID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not fetch stories: %s", err)
		os.Exit(1)
	}

	response := []resource.Version{}
	for _, story := range stories {
		response = append(response, resource.Version{StoryID: story.ID})
	}
	err = json.NewEncoder(os.Stdout).Encode(response)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not print response: %s", err)
		os.Exit(1)
	}
}
