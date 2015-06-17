package check

import (
	"fmt"
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
	repository      resource.Repository
	stories         []tracker.Story
}

func NewTrackerGitBranchCheck(
	startingVersion resource.Version,
	repository resource.Repository,
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

	err := c.repository.Fetch()
	if err != nil {
		return []resource.Version{}, fmt.Errorf("Could not fetch: %s", err)
	}
	remoteBranches, err := c.repository.RemoteBranches()
	if err != nil {
		return []resource.Version{}, fmt.Errorf("Could not list remote branches: %s", err)
	}

	if c.startingVersion.StoryID == "" {
		versions, err = c.latestStoryBranchRef(remoteBranches)
		if err != nil {
			return []resource.Version{}, fmt.Errorf("Could not find latest story branch ref from remote branches %v: %s", remoteBranches, err)
		}
	} else {
		versions, err = c.storyBranchRefsSinceStartingVersion(remoteBranches)
		if err != nil {
			return []resource.Version{}, fmt.Errorf("Could not find new story branch refs from remote branches %v: %s", remoteBranches, err)
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
				timestamp, err := c.repository.RefCommitTimestamp(branch)
				if err != nil {
					return []resource.Version{}, fmt.Errorf("Could not get ref commit timestamp for %s: %s", branch, err)
				}

				if timestamp > latestTime {
					ref, err := c.repository.LatestRef(branch)
					if err != nil {
						return []resource.Version{}, fmt.Errorf("Could not get latest SHA for %s: %s", branch, err)
					}
					versions = []resource.Version{{StoryID: strconv.Itoa(story.ID), Ref: ref, Timestamp: timestamp}}
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
				refs, err := c.repository.RefsSinceTimestamp(branch, c.startingVersion.Timestamp)
				if err != nil {
					return []resource.Version{}, fmt.Errorf("Could not get refs since time %d for %s: %s", c.startingVersion.Timestamp, branch, err)
				}

				// Collect versions for later sorting
				for _, ref := range refs {
					ref = strings.Trim(ref, "\"")
					if ref != "" && ref != c.startingVersion.Ref {
						timestamp, err := c.repository.RefCommitTimestamp(ref)
						if err != nil {
							return []resource.Version{}, fmt.Errorf("Could not get timestamp for %s: %s", ref, err)
						}
						versions = append(versions, resource.Version{StoryID: strconv.Itoa(story.ID), Ref: ref, Timestamp: timestamp})
					}
				}

				break
			}
		}
	}
	return sortVersionsByTimestamp(versions), nil
}

func isStoryBranch(branch string, story tracker.Story) bool {
	return strings.Contains(branch, strconv.Itoa(story.ID))
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
