package check

import (
	"github.com/adamstegman/go-tracker"
)

type TrackerProjectClient interface {
	Stories(tracker.StoriesQuery) ([]tracker.Story, error)
	StoryActivity(int, tracker.ActivityQuery) ([]tracker.Activity, error)
}

type TrackerGitBranchCheck interface {
	StoriesFinishedAfterStory(storyID int) ([]tracker.Story, error)
}

type trackerGitBranchCheck struct {
	projectClient TrackerProjectClient
}

func NewTrackerGitBranchCheck(projectClient TrackerProjectClient) trackerGitBranchCheck {
	return trackerGitBranchCheck{
		projectClient: projectClient,
	}
}

func (c trackerGitBranchCheck) StoriesFinishedAfterStory(storyID int) ([]tracker.Story, error) {
	finishedQuery := tracker.StoriesQuery{State: tracker.StoryStateFinished}
	finishedStories, err := c.projectClient.Stories(finishedQuery)
	if err != nil {
		return []tracker.Story{}, err
	}

	// TODO:
	// * get activity for given story to find beginning time
	// * get activity for each story
	//     * NOTE: use occurred_after parameter in activity request
	// * select stories finished after the given story was finished
	latestFinishedStory, err := c.findLatestFinishedStory(finishedStories)
	if err != nil {
		return []tracker.Story{}, err
	}
	return []tracker.Story{latestFinishedStory}, nil
}

func (c trackerGitBranchCheck) findLatestFinishedStory(stories []tracker.Story) (tracker.Story, error) {
	var (
		latestFinishedStoryMillis int64
		latestFinishedStory       tracker.Story
	)
	for _, story := range stories {
		activities, err := c.projectClient.StoryActivity(story.ID, tracker.ActivityQuery{})
		if err != nil {
			return tracker.Story{}, err
		}
		for _, activity := range activities {
			if activity.Highlight == "finished" && activity.OccurredAt > latestFinishedStoryMillis {
				latestFinishedStoryMillis = activity.OccurredAt
				latestFinishedStory = story
				break
			}
		}
	}
	return latestFinishedStory, nil
}
