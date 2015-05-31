package out_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"github.com/adamstegman/tracker-git-branch-resource"
	"github.com/adamstegman/tracker-git-branch-resource/out"
)

var _ = Describe("Out", func() {
	var tmpdir string

	var outCmd *exec.Cmd

	BeforeEach(func() {
		var err error

		tmpdir, err = ioutil.TempDir("", "out-tmp")
		Ω(err).ShouldNot(HaveOccurred())
		err = os.MkdirAll(tmpdir, 0755)
		Ω(err).ShouldNot(HaveOccurred())

		outCmd = exec.Command(outPath, tmpdir)
	})

	AfterEach(func() {
		os.RemoveAll(tmpdir)
	})

	Context("when executed", func() {
		var request out.OutRequest
		var response out.OutResponse

		BeforeEach(func() {
			request = out.OutRequest{
				Source: resource.Source{
					Token:     "abc",
					ProjectID: "1234",
				},
				Params: struct{}{},
			}
			response = out.OutResponse{}
		})

		It("outputs the current time", func() {
			session := runCommand(outCmd, request)

			err := json.Unmarshal(session.Out.Contents(), &response)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(response.Version.Time).Should(BeTemporally("~", time.Now(), time.Second))
		})
	})
})

func runCommand(outCmd *exec.Cmd, request out.OutRequest) *Session {
	stdin, err := outCmd.StdinPipe()
	Ω(err).ShouldNot(HaveOccurred())

	session, err := Start(outCmd, GinkgoWriter, GinkgoWriter)
	Ω(err).ShouldNot(HaveOccurred())
	err = json.NewEncoder(stdin).Encode(request)
	Ω(err).ShouldNot(HaveOccurred())
	Eventually(session).Should(Exit(0))

	return session
}
