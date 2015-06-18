package in_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"

	"github.com/adamstegman/tracker-git-branch-resource"
	"github.com/adamstegman/tracker-git-branch-resource/in"
)

var _ = Describe("In", func() {
	var (
		tmpDir   string
		request  in.InRequest
		response in.InResponse
	)
	BeforeEach(func() {
		var err error
		tmpDir, err = ioutil.TempDir("", "tracker_resource_in")
		Expect(err).NotTo(HaveOccurred())
	})

	JustBeforeEach(func() {
		binPath, err := gexec.Build("github.com/adamstegman/tracker-git-branch-resource/in/cmd/in")
		Expect(err).NotTo(HaveOccurred())

		stdin := &bytes.Buffer{}
		err = json.NewEncoder(stdin).Encode(request)
		Expect(err).NotTo(HaveOccurred())

		cmd := exec.Command(binPath, tmpDir)
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

	AfterEach(func() {
		err := os.RemoveAll(tmpDir)
		Expect(err).NotTo(HaveOccurred())
	})

	BeforeEach(func() {
		sourceRepo, err := filepath.Abs("..")
		Expect(err).NotTo(HaveOccurred())
		request = in.InRequest{
			Source: resource.Source{
				Token:     "trackerToken",
				ProjectID: "123456",
				Repo:      sourceRepo,
			},
			Version: resource.Version{StoryID: "9999", Ref: "42f809095d489e446713cf20fdc3d30e5faaa4c9", Timestamp: "1433829600"},
		}
	})

	It("clones the ref in the given directory and outputs that version", func() {
		repository := resource.NewRepository("", tmpDir)
		ref, err := repository.LatestRef("HEAD")
		Expect(err).NotTo(HaveOccurred())
		Expect(ref).To(Equal(request.Version.Ref))

		Expect(response.Version).To(Equal(request.Version))
	})

	It("outputs metadata about the story and ref", func() {
		Expect(response.Metadata).To(Equal([]resource.MetadataPair{
			{Name: "commit", Value: "42f809095d489e446713cf20fdc3d30e5faaa4c9"},
			{Name: "author", Value: "Adam Stegman"},
			{Name: "author_date", Value: "2015-06-08 23:00:00 -0700"},
			{Name: "committer", Value: "Adam Stegman"},
			{Name: "committer_date", Value: "2015-06-08 23:00:00 -0700"},
			{Name: "message", Value: "Update\n"},
			{Name: "story_url", Value: "https://www.pivotaltracker.com/story/show/9999"},
		}))
	})

	Context("when the target directory does not exist", func() {
		BeforeEach(func() {
			err := os.RemoveAll(tmpDir)
			Expect(err).NotTo(HaveOccurred())
		})

		It("creates the target directory", func() {
			tdInfo, err := os.Stat(tmpDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(tdInfo.IsDir()).To(BeTrue())
		})
	})
})
