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
		stories := []tracker.Story{{ID: 9999}, {ID: 1234}}
		fakeTrackerProjectClient.AddStoriesForQuery(stories, tracker.StoriesQuery{State: tracker.StoryStateFinished})
		activities9999 := []tracker.Activity{
			{Highlight: "finished", OccurredAt: 1999999999999},
		}
		fakeTrackerProjectClient.AddActivityForStoryID(activities9999, 9999)
		activities1234 := []tracker.Activity{
			{Highlight: "finished", OccurredAt: 1000000000000},
		}
		fakeTrackerProjectClient.AddActivityForStoryID(activities1234, 1234)

		trackerCheck = check.NewTrackerGitBranchCheck(fakeTrackerProjectClient)
	})

	Context("when no known story is given", func() {
		It("returns the last finished story", func() {
			stories, err := trackerCheck.StoriesFinishedAfterStory(0)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(stories)).To(Equal(1))
			Expect(stories[0].ID).To(Equal(9999))
		})
	})

	Context("when a story ID is given", func() {
		XIt("returns all stories finished after the given story", func() {
		})
	})
})
