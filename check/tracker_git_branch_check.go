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

	if storyID == 0 {
		latestFinishedStory, err := c.findLatestFinishedStory(finishedStories)
		if err != nil {
			return []tracker.Story{}, err
		}
		return []tracker.Story{latestFinishedStory}, nil
	} else {
		storiesFinishedAfterGivenStory, err := c.findStoriesFinishedAfterStory(storyID, finishedStories)
		if err != nil {
			return []tracker.Story{}, err
		}
		return storiesFinishedAfterGivenStory, nil
	}
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

func (c trackerGitBranchCheck) findStoriesFinishedAfterStory(storyID int, finishedStories []tracker.Story) ([]tracker.Story, error) {
	// get activity for given story to find beginning time
	var givenStoryDeliveredTime int64
	activities, err := c.projectClient.StoryActivity(storyID, tracker.ActivityQuery{})
	if err != nil {
		return []tracker.Story{}, err
	}
	for _, activity := range activities {
		if activity.Highlight == "finished" {
			givenStoryDeliveredTime = activity.OccurredAt
			break
		}
	}

	// select stories finished after the given story was finished
	afterGivenStoryQuery := tracker.ActivityQuery{OccurredAfter: givenStoryDeliveredTime}
	var storiesFinishedAfterGivenStory []tracker.Story
	for _, story := range finishedStories {
		activities, err := c.projectClient.StoryActivity(story.ID, afterGivenStoryQuery)
		if err != nil {
			return []tracker.Story{}, err
		}
		for _, activity := range activities {
			if activity.Highlight == "finished" {
				storiesFinishedAfterGivenStory = append(storiesFinishedAfterGivenStory, story)
				break
			}
		}
	}
	return storiesFinishedAfterGivenStory, nil
}
