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
			request.Version = resource.Version{}
		})

		Context("and finished stories are found", func() {
			BeforeEach(func() {
				stories := []tracker.Story{{ID: 9999}, {ID: 1234}}
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories", "date_format=millis&with_state=finished"),
						ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, stories),
					),
				)
				activities9999 := []tracker.Activity{
					{Highlight: "finished", OccurredAt: 1999999999999},
				}
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories/9999/activity", "date_format=millis"),
						ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, activities9999),
					),
				)
				activities1234 := []tracker.Activity{
					{Highlight: "finished", OccurredAt: 1000000000000},
				}
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories/1234/activity", "date_format=millis"),
						ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, activities1234),
					),
				)
			})

			It("returns the last finished story", func() {
				Expect(response).To(Equal([]resource.Version{{StoryID: 9999}}))
			})
		})

		Context("and no finished stories are found", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories", "date_format=millis&with_state=finished"),
						ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, []tracker.Story{}),
					),
				)
			})

			It("returns an empty list", func() {
				Expect(response).To(Equal([]resource.Version{}))
			})
		})
	})

	Context("when a story ID is given", func() {
		BeforeEach(func() {
			request.Version = resource.Version{StoryID: 1234}
			stories := []tracker.Story{{ID: 9999}, {ID: 5454}, {ID: 1234}}
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories", "date_format=millis&with_state=finished"),
					ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
					ghttp.RespondWithJSONEncoded(http.StatusOK, stories),
				),
			)
			activities1234 := []tracker.Activity{
				{Highlight: "finished", OccurredAt: 1000000000000},
			}
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories/1234/activity", "date_format=millis"),
					ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
					ghttp.RespondWithJSONEncoded(http.StatusOK, activities1234),
				),
			)
			activities9999 := []tracker.Activity{
				{Highlight: "finished", OccurredAt: 1999999999999},
			}
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories/9999/activity", "date_format=millis&occurred_after=1000000000000"),
					ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
					ghttp.RespondWithJSONEncoded(http.StatusOK, activities9999),
				),
			)
			activities5454 := []tracker.Activity{
				{Highlight: "finished", OccurredAt: 1545454545454},
			}
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories/5454/activity", "date_format=millis&occurred_after=1000000000000"),
					ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
					ghttp.RespondWithJSONEncoded(http.StatusOK, activities5454),
				),
			)
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories/1234/activity", "date_format=millis&occurred_after=1000000000000"),
					ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
					ghttp.RespondWithJSONEncoded(http.StatusOK, []tracker.Activity{}),
				),
			)
		})

		It("returns all stories finished after the given story", func() {
			Expect(response).To(Equal([]resource.Version{
				{StoryID: 9999},
				{StoryID: 5454},
			}))
		})
	})
})
