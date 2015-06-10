package tracker_test

import (
	"encoding/json"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/xoebus/go-tracker"
)

var _ = Describe("Me", func() {
	It("has attributes", func() {
		var me tracker.Me
		reader := strings.NewReader(Fixture("me.json"))
		err := json.NewDecoder(reader).Decode(&me)
		Ω(err).ToNot(HaveOccurred())

		Ω(me.Username).To(Equal("vader"))
		Ω(me.Name).To(Equal("Darth Vader"))
		Ω(me.Initials).To(Equal("DV"))
		Ω(me.ID).To(Equal(101))
		Ω(me.Email).To(Equal("vader@deathstar.mil"))
	})
})

var _ = Describe("Story", func() {
	It("has attributes", func() {
		var stories []tracker.Story
		reader := strings.NewReader(Fixture("stories.json"))
		err := json.NewDecoder(reader).Decode(&stories)
		Ω(err).ToNot(HaveOccurred())
		story := stories[0]

		Ω(story.ID).Should(Equal(560))
		Ω(story.Name).Should(Equal("Tractor beam loses power intermittently"))
	})
})

var _ = Describe("Activity", func() {
	It("has attributes", func() {
		var activities []tracker.Activity
		reader := strings.NewReader(Fixture("activities.json"))
		err := json.NewDecoder(reader).Decode(&activities)
		Ω(err).ToNot(HaveOccurred())
		activity := activities[0]

		Ω(activity.GUID).Should(Equal("99_45"))
		Ω(activity.Message).Should(Equal("Darth Vader started this feature"))
	})
})
