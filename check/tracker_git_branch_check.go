package check

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/xoebus/go-tracker"

	"github.com/adamstegman/tracker-git-branch-resource"
)

type TrackerGitBranchCheck interface {
	NewVersions() ([]resource.Version, error)
}

type trackerGitBranchCheck struct {
	startingVersion resource.Version
	repository      string
	stories         []tracker.Story
}

func NewTrackerGitBranchCheck(
	startingVersion resource.Version,
	repository string,
	stories []tracker.Story,
) trackerGitBranchCheck {
	return trackerGitBranchCheck{
		startingVersion: startingVersion,
		repository:      repository,
		stories:         stories,
	}
}

func (c trackerGitBranchCheck) NewVersions() ([]resource.Version, error) {
	versions := []resource.Version{}

	err := c.fetchBranches()
	if err != nil {
		return []resource.Version{}, err
	}
	remoteBranches, err := c.remoteBranches()
	if err != nil {
		return []resource.Version{}, err
	}

	if c.startingVersion.StoryID == 0 {
		versions, err = c.latestStoryBranchRef(remoteBranches)
		if err != nil {
			return []resource.Version{}, err
		}
	} else {
		versions, err = c.storyBranchRefsSinceStartingVersion(remoteBranches)
		if err != nil {
			return []resource.Version{}, err
		}
	}

	return versions, nil
}

func (c trackerGitBranchCheck) latestStoryBranchRef(remoteBranches []string) ([]resource.Version, error) {
	var latestTime int64
	versions := []resource.Version{}
	for _, story := range c.stories {
		for _, branch := range remoteBranches {
			if isStoryBranch(branch, story) {
				timestamp, err := c.refCommitTimestamp(branch)
				if err != nil {
					return []resource.Version{}, err
				}

				if timestamp > latestTime {
					ref, err := c.latestRef(branch)
					if err != nil {
						return []resource.Version{}, err
					}
					versions = []resource.Version{{StoryID: story.ID, Ref: ref, Timestamp: timestamp}}
					latestTime = timestamp
				}

				break
			}
		}
	}
	return versions, nil
}

func (c trackerGitBranchCheck) storyBranchRefsSinceStartingVersion(remoteBranches []string) ([]resource.Version, error) {
	versions := []resource.Version{}
	for _, story := range c.stories {
		for _, branch := range remoteBranches {
			if isStoryBranch(branch, story) {
				refs, err := c.refsSinceTimestamp(branch, c.startingVersion.Timestamp)
				if err != nil {
					return []resource.Version{}, err
				}

				// Collect versions for later sorting
				for _, ref := range refs {
					ref = strings.Trim(ref, "\"")
					if ref != "" && ref != c.startingVersion.Ref {
						timestamp, err := c.refCommitTimestamp(ref)
						if err != nil {
							return []resource.Version{}, err
						}
						versions = append(versions, resource.Version{StoryID: story.ID, Ref: ref, Timestamp: timestamp})
					}
				}

				break
			}
		}
	}
	return sortVersionsByTimestamp(versions), nil
}

func (c trackerGitBranchCheck) fetchBranches() error {
	fetchCmd := exec.Command("git", "fetch", "origin")
	fetchCmd.Dir = c.repository
	return fetchCmd.Run()
}

func (c trackerGitBranchCheck) remoteBranches() ([]string, error) {
	branchCmd := exec.Command("git", "branch", "-r")
	branchCmd.Dir = c.repository
	var branchBytes bytes.Buffer
	branchCmd.Stdout = &branchBytes
	err := branchCmd.Run()
	if err != nil {
		return []string{}, err
	}
	branches := strings.Split(branchBytes.String(), "\n")

	trimmedBranches := []string{}
	for _, branch := range branches {
		trimmedBranches = append(trimmedBranches, strings.TrimSpace(branch))
	}
	return trimmedBranches, nil
}

func isStoryBranch(branch string, story tracker.Story) bool {
	return strings.Contains(branch, strconv.Itoa(story.ID))
}

func (c trackerGitBranchCheck) refCommitTimestamp(ref string) (int64, error) {
	timeCmd := exec.Command("git", "show", "-s", "--format=\"%ct\"", ref)
	timeCmd.Dir = c.repository
	var timeBytes bytes.Buffer
	timeCmd.Stdout = &timeBytes
	err := timeCmd.Run()
	if err != nil {
		return 0, err
	}
	timeString := strings.Trim(strings.TrimSpace(timeBytes.String()), "\"")
	timestamp, err := strconv.ParseInt(timeString, 10, 64)
	if err != nil {
		return 0, err
	}
	return timestamp, nil
}

func (c trackerGitBranchCheck) latestRef(branch string) (string, error) {
	refCmd := exec.Command("git", "show", "-s", "--format=\"%H\"", branch)
	refCmd.Dir = c.repository
	var refBytes bytes.Buffer
	refCmd.Stdout = &refBytes
	err := refCmd.Run()
	if err != nil {
		return "", err
	}
	return strings.Trim(strings.TrimSpace(refBytes.String()), "\""), nil
}

func (c trackerGitBranchCheck) refsSinceTimestamp(branch string, timestamp int64) ([]string, error) {
	refsCmd := exec.Command("git", "log", fmt.Sprintf("--since=%d", timestamp), "--format=\"%H\"", branch)
	refsCmd.Dir = c.repository
	var refsBytes bytes.Buffer
	refsCmd.Stdout = &refsBytes
	err := refsCmd.Run()
	if err != nil {
		return []string{}, err
	}
	return strings.Split(refsBytes.String(), "\n"), nil
}

func sortVersionsByTimestamp(versions []resource.Version) []resource.Version {
	// Sort refs by timestamp from oldest to newest
	sortedVersions := []resource.Version{}
	for _, version := range versions {
		var postIndex int
		for _, v := range sortedVersions {
			if v.Timestamp < version.Timestamp {
				postIndex = postIndex + 1
			} else {
				break
			}
		}
		if postIndex == 0 {
			sortedVersions = append([]resource.Version{version}, sortedVersions...)
		} else if postIndex == len(sortedVersions) {
			sortedVersions = append(sortedVersions, version)
		} else {
			sortedVersions = append(sortedVersions[0:postIndex], append([]resource.Version{version}, sortedVersions[postIndex:]...)...)
		}
	}
	return sortedVersions
}
