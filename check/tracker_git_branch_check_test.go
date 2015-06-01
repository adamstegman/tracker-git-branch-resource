package check_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/adamstegman/go-tracker"

	"github.com/adamstegman/tracker-git-branch-resource/check"
	checktest "github.com/adamstegman/tracker-git-branch-resource/check/test"
)

var _ = Describe("trackerGitBranchCheck", func() {
	var (
		fakeTrackerProjectClient *checktest.FakeTrackerProjectClient
		trackerCheck             check.TrackerGitBranchCheck
	)

	BeforeEach(func() {
		fakeTrackerProjectClient = checktest.NewFakeTrackerProjectClient()
		stories := []tracker.Story{{ID: 9999}, {ID: 5454}, {ID: 1234}}
		fakeTrackerProjectClient.AddStoriesForQuery(stories, tracker.StoriesQuery{State: tracker.StoryStateFinished})
		activities9999 := []tracker.Activity{
			{Highlight: "finished", OccurredAt: 1999999999999},
		}
		fakeTrackerProjectClient.AddActivityForStoryIDAndQuery(activities9999, 9999, tracker.ActivityQuery{})
		fakeTrackerProjectClient.AddActivityForStoryIDAndQuery(activities9999, 9999, tracker.ActivityQuery{OccurredAfter: 1000000000000})
		activities5454 := []tracker.Activity{
			{Highlight: "finished", OccurredAt: 1545454545454},
		}
		fakeTrackerProjectClient.AddActivityForStoryIDAndQuery(activities5454, 5454, tracker.ActivityQuery{})
		fakeTrackerProjectClient.AddActivityForStoryIDAndQuery(activities5454, 5454, tracker.ActivityQuery{OccurredAfter: 1000000000000})
		activities1234 := []tracker.Activity{
			{Highlight: "finished", OccurredAt: 1000000000000},
		}
		fakeTrackerProjectClient.AddActivityForStoryIDAndQuery(activities1234, 1234, tracker.ActivityQuery{})
		fakeTrackerProjectClient.AddActivityForStoryIDAndQuery([]tracker.Activity{}, 1234, tracker.ActivityQuery{OccurredAfter: 1000000000000})

		trackerCheck = check.NewTrackerGitBranchCheck(fakeTrackerProjectClient)
	})

	Context("when no known story is given", func() {
		It("returns the last finished story", func() {
			stories, err := trackerCheck.StoriesFinishedAfterStory(0)
			Expect(err).NotTo(HaveOccurred())
			Expect(stories).To(Equal([]tracker.Story{{ID: 9999}}))
		})
	})

	Context("when a story ID is given", func() {
		It("returns all stories finished after the given story", func() {
			stories, err := trackerCheck.StoriesFinishedAfterStory(1234)
			Expect(err).NotTo(HaveOccurred())
			Expect(stories).To(Equal([]tracker.Story{{ID: 9999}, {ID: 5454}}))
		})
	})
})
