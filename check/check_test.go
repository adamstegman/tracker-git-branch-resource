package check_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"

	"github.com/adamstegman/go-tracker"

	"github.com/adamstegman/tracker-git-branch-resource"
	"github.com/adamstegman/tracker-git-branch-resource/check"
)

var _ = Describe("check", func() {
	var (
		server   *ghttp.Server
		request  *check.Request
		response []resource.Version
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		stories := []tracker.Story{{ID: 9999}, {ID: 1234}}
		storiesJson, err := json.Marshal(stories)
		Expect(err).NotTo(HaveOccurred())
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories", "date_format=millis&with_state=finished"),
				ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
				ghttp.RespondWith(http.StatusOK, storiesJson),
			),
		)
		activities9999 := []tracker.Activity{
			{Highlight: "finished", OccurredAt: 1999999999999},
		}
		activities9999Json, err := json.Marshal(activities9999)
		Expect(err).NotTo(HaveOccurred())
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories/9999/activity", "date_format=millis"),
				ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
				ghttp.RespondWith(http.StatusOK, activities9999Json),
			),
		)
		activities1234 := []tracker.Activity{
			{Highlight: "finished", OccurredAt: 1000000000000},
		}
		activities1234Json, err := json.Marshal(activities1234)
		Expect(err).NotTo(HaveOccurred())
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories/1234/activity", "date_format=millis"),
				ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
				ghttp.RespondWith(http.StatusOK, activities1234Json),
			),
		)

		request = &check.Request{
			Source: resource.Source{
				Token:      "trackerToken",
				ProjectID:  "123456",
				TrackerURL: server.URL(),
			},
		}
	})
	AfterEach(func() {
		server.Close()
	})

	JustBeforeEach(func() {
		binPath, err := gexec.Build("github.com/adamstegman/tracker-git-branch-resource/check/cmd/check")
		Expect(err).NotTo(HaveOccurred())

		stdin := &bytes.Buffer{}
		err = json.NewEncoder(stdin).Encode(request)
		Expect(err).NotTo(HaveOccurred())

		cmd := exec.Command(binPath)
		cmd.Stdin = stdin

		session, err := gexec.Start(
			cmd,
			GinkgoWriter,
			GinkgoWriter,
		)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))

		err = json.Unmarshal(session.Out.Contents(), &response)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when no known story is given", func() {
		BeforeEach(func() {
			// FIXME: is this accurate, or is it null? or missing entirely?
			// use omitempty
			request.Version = resource.Version{}
		})

		It("returns the last finished story", func() {
			Expect(len(response)).To(Equal(1))
			Expect(response[0]).To(Equal(resource.Version{StoryID: 9999}))
		})
	})

	Context("when a story ID is given", func() {
		XIt("returns all stories finished after the given story", func() {
		})
	})
})
