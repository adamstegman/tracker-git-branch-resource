package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		fmt.Fprintf(os.Stderr, "Could not parse input: %s\n", err)
		os.Exit(1)
	}

	targetDir, err := ioutil.TempDir("", "tracker-git-branch-resource-check")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create temporary directory: %s\n", err)
		os.Exit(1)
	}
	var keyFile string
	if request.Source.PrivateKey != "" {
		keyFile, err = resource.CreateKeyFile(request.Source.PrivateKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not create keyfile: %s", err)
			os.Exit(1)
		}
		defer os.Remove(keyFile)
	}
	repository := resource.NewRepository(request.Source.Repo, targetDir, keyFile)
	err = repository.Clone()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not clone repo %s: %s\n", request.Source.Repo, err)
		os.Exit(1)
	}

	if request.Source.TrackerURL != "" {
		tracker.DefaultURL = request.Source.TrackerURL
	}
	trackerToken := request.Source.Token
	stories := []tracker.Story{}
	for _, projectID := range request.Source.Projects {
		trackerProjectID, err := strconv.Atoi(projectID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid Tracker project ID %s: %s\n", projectID, err)
			os.Exit(1)
		}
		projectClient := tracker.NewClient(trackerToken).InProject(trackerProjectID)

		finishedQuery := tracker.StoriesQuery{State: tracker.StoryStateFinished}
		finishedStories, err := projectClient.Stories(finishedQuery)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not fetch finished stories: %s\n", err)
			os.Exit(1)
		}
		stories = append(stories, finishedStories...)
		deliveredQuery := tracker.StoriesQuery{State: tracker.StoryStateDelivered}
		deliveredStories, err := projectClient.Stories(deliveredQuery)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not fetch delivered stories: %s\n", err)
			os.Exit(1)
		}
		stories = append(stories, deliveredStories...)
	}

	trackerGitBranchCheck := check.NewTrackerGitBranchCheck(request.Version, repository, stories)
	versions, err := trackerGitBranchCheck.NewVersions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not find versions: %s\n", err)
		os.Exit(1)
	}

	err = json.NewEncoder(os.Stdout).Encode(versions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not print response: %s\n", err)
		os.Exit(1)
	}
}
