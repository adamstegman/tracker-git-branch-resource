package test

import (
	"github.com/adamstegman/go-tracker"
)

type FakeTrackerProjectClient struct {
	StoriesByQuery      map[tracker.StoriesQuery][]tracker.Story
	ActivitiesByStoryID map[int][]tracker.Activity
}

func NewFakeTrackerProjectClient() *FakeTrackerProjectClient {
	return &FakeTrackerProjectClient{
		StoriesByQuery:      make(map[tracker.StoriesQuery][]tracker.Story),
		ActivitiesByStoryID: make(map[int][]tracker.Activity),
	}
}

func (c *FakeTrackerProjectClient) AddStoriesForQuery(stories []tracker.Story, query tracker.StoriesQuery) {
	c.StoriesByQuery[query] = stories
}

func (c *FakeTrackerProjectClient) AddActivityForStoryID(activities []tracker.Activity, storyID int) {
	c.ActivitiesByStoryID[storyID] = activities
}

func (c *FakeTrackerProjectClient) Stories(query tracker.StoriesQuery) ([]tracker.Story, error) {
	return c.StoriesByQuery[query], nil
}

func (c *FakeTrackerProjectClient) StoryActivity(storyID int, query tracker.ActivityQuery) ([]tracker.Activity, error) {
	return c.ActivitiesByStoryID[storyID], nil
}
