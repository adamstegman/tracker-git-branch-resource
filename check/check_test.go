package check_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
	"github.com/xoebus/go-tracker"

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
		sourceRepo, err := filepath.Abs("..")
		Expect(err).NotTo(HaveOccurred())
		request = &check.Request{
			Source: resource.Source{
				Token:      "trackerToken",
				ProjectID:  "123456",
				TrackerURL: server.URL(),
				Repo:       sourceRepo,
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

		Eventually(session, 10).Should(gexec.Exit(0))

		err = json.Unmarshal(session.Out.Contents(), &response)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when no known version is given", func() {
		BeforeEach(func() {
			request.Version = resource.Version{}
		})

		Context("and story branches are found", func() {
			BeforeEach(func() {
				finishedStories := []tracker.Story{{ID: 1234}}
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories", "date_format=millis&with_state=finished"),
						ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, finishedStories),
					),
				)
				deliveredStories := []tracker.Story{{ID: 9999}}
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories", "date_format=millis&with_state=delivered"),
						ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, deliveredStories),
					),
				)
			})

			It("finds the latest ref out of the finished or delivered stories", func() {
				Expect(response).To(Equal([]resource.Version{
					{StoryID: "9999", Ref: "42f809095d489e446713cf20fdc3d30e5faaa4c9", Timestamp: 1433829600},
				}))
			})
		})

		Context("and no story branches are found", func() {
			BeforeEach(func() {
				finishedStories := []tracker.Story{{ID: 0000}}
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories", "date_format=millis&with_state=finished"),
						ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, finishedStories),
					),
				)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories", "date_format=millis&with_state=delivered"),
						ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, []tracker.Story{}),
					),
				)
			})

			It("returns an empty list", func() {
				Expect(response).To(Equal([]resource.Version{}))
			})
		})

		Context("and no finished or delivered stories are found", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories", "date_format=millis&with_state=finished"),
						ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
						ghttp.RespondWithJSONEncoded(http.StatusOK, []tracker.Story{}),
					),
				)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories", "date_format=millis&with_state=delivered"),
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

	Context("when a version is given", func() {
		BeforeEach(func() {
			request.Version = resource.Version{StoryID: "5454", Ref: "d6e5a26bc1e0b39b74f7aceb5ef651cb729cc5d0", Timestamp: 1433800800}
		})

		BeforeEach(func() {
			finishedStories := []tracker.Story{{ID: 5454}, {ID: 1234}}
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories", "date_format=millis&with_state=finished"),
					ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
					ghttp.RespondWithJSONEncoded(http.StatusOK, finishedStories),
				),
			)
			deliveredStories := []tracker.Story{{ID: 0000}, {ID: 9999}}
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories", "date_format=millis&with_state=delivered"),
					ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
					ghttp.RespondWithJSONEncoded(http.StatusOK, deliveredStories),
				),
			)
		})

		It("returns all refs in all finished and delivered story branches after the given ref, in chronological order", func() {
			Expect(response).To(Equal([]resource.Version{
				{StoryID: "9999", Ref: "1ad88c443704b2531d471b99c21489f3c4deb974", Timestamp: 1433818800},
				{StoryID: "1234", Ref: "98bc2acea806e1a507d70ae3ce21cee3cd2d6c38", Timestamp: 1433822400},
				{StoryID: "9999", Ref: "42f809095d489e446713cf20fdc3d30e5faaa4c9", Timestamp: 1433829600},
			}))
		})
	})

	Context("when a deleted version is given", func() {
		BeforeEach(func() {
			request.Version = resource.Version{StoryID: "1000", Ref: "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0", Timestamp: 1433800801}
		})

		BeforeEach(func() {
			finishedStories := []tracker.Story{{ID: 5454}, {ID: 1234}}
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories", "date_format=millis&with_state=finished"),
					ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
					ghttp.RespondWithJSONEncoded(http.StatusOK, finishedStories),
				),
			)
			deliveredStories := []tracker.Story{{ID: 1000}, {ID: 9999}}
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/projects/123456/stories", "date_format=millis&with_state=delivered"),
					ghttp.VerifyHeaderKV("X-Trackertoken", "trackerToken"),
					ghttp.RespondWithJSONEncoded(http.StatusOK, deliveredStories),
				),
			)
		})

		It("returns all refs in all finished and delivered story branches after the given ref, in chronological order", func() {
			Expect(response).To(Equal([]resource.Version{
				{StoryID: "9999", Ref: "1ad88c443704b2531d471b99c21489f3c4deb974", Timestamp: 1433818800},
				{StoryID: "1234", Ref: "98bc2acea806e1a507d70ae3ce21cee3cd2d6c38", Timestamp: 1433822400},
				{StoryID: "9999", Ref: "42f809095d489e446713cf20fdc3d30e5faaa4c9", Timestamp: 1433829600},
			}))
		})
	})
})
