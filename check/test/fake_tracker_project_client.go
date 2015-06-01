package test

import (
	"github.com/adamstegman/go-tracker"
)

type FakeTrackerProjectClient struct {
	StoriesByQuery           map[tracker.StoriesQuery][]tracker.Story
	ActivityQueriesByStoryID map[int]map[tracker.ActivityQuery][]tracker.Activity
}

func NewFakeTrackerProjectClient() *FakeTrackerProjectClient {
	return &FakeTrackerProjectClient{
		StoriesByQuery:           make(map[tracker.StoriesQuery][]tracker.Story),
		ActivityQueriesByStoryID: make(map[int]map[tracker.ActivityQuery][]tracker.Activity),
	}
}

func (c *FakeTrackerProjectClient) AddStoriesForQuery(stories []tracker.Story, query tracker.StoriesQuery) {
	c.StoriesByQuery[query] = stories
}

func (c *FakeTrackerProjectClient) AddActivityForStoryIDAndQuery(activities []tracker.Activity, storyID int, query tracker.ActivityQuery) {
	if c.ActivityQueriesByStoryID[storyID] == nil {
		c.ActivityQueriesByStoryID[storyID] = make(map[tracker.ActivityQuery][]tracker.Activity)
	}
	c.ActivityQueriesByStoryID[storyID][query] = activities
}

func (c *FakeTrackerProjectClient) Stories(query tracker.StoriesQuery) ([]tracker.Story, error) {
	return c.StoriesByQuery[query], nil
}

func (c *FakeTrackerProjectClient) StoryActivity(storyID int, query tracker.ActivityQuery) ([]tracker.Activity, error) {
	return c.ActivityQueriesByStoryID[storyID][query], nil
}
