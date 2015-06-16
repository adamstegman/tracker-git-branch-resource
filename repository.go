package resource

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

type Repository struct {
	Dir    string
	source string
}

func NewRepository(source string, dir string) Repository {
	return Repository{
		Dir:    dir,
		source: source,
	}
}

func (r Repository) Clone() error {
	if r.Dir == "" {
		var err error
		r.Dir, err = ioutil.TempDir("", "tracker-git-branch-resource-repository")
		if err != nil {
			return err
		}
	}
	cloneCmd := exec.Command("git", "clone", r.source, r.Dir)
	return cloneCmd.Run()
}

func (r Repository) Fetch() error {
	fetchCmd := exec.Command("git", "fetch", "origin")
	fetchCmd.Dir = r.Dir
	return fetchCmd.Run()
}

func (r Repository) CheckoutRef(ref string) error {
	checkoutCmd := exec.Command("git", "checkout", ref)
	checkoutCmd.Dir = r.Dir
	return checkoutCmd.Run()
}

func (r Repository) RemoteBranches() ([]string, error) {
	branchesOutput, err := r.cmdOutput("git", "branch", "-r")
	if err != nil {
		return []string{}, err
	}
	branches := strings.Split(branchesOutput, "\n")

	trimmedBranches := []string{}
	for _, branch := range branches {
		if !strings.Contains(branch, "origin/HEAD ->") {
			trimmedBranches = append(trimmedBranches, strings.TrimSpace(branch))
		}
	}
	return trimmedBranches, nil
}

func (r Repository) RefAuthorName(ref string) (string, error) {
	nameOutput, err := r.cmdOutput("git", "show", "-s", "--format=\"%an\"", ref)
	if err != nil {
		return "", err
	}
	return strings.Trim(nameOutput, "\""), nil
}

func (r Repository) RefAuthorDate(ref string) (string, error) {
	timeOutput, err := r.cmdOutput("git", "show", "-s", "--format=\"%ai\"", ref)
	if err != nil {
		return "", err
	}
	return strings.Trim(timeOutput, "\""), nil
}

func (r Repository) RefCommitName(ref string) (string, error) {
	nameOutput, err := r.cmdOutput("git", "show", "-s", "--format=\"%cn\"", ref)
	if err != nil {
		return "", err
	}
	return strings.Trim(nameOutput, "\""), nil
}

func (r Repository) RefCommitDate(ref string) (string, error) {
	timeOutput, err := r.cmdOutput("git", "show", "-s", "--format=\"%ci\"", ref)
	if err != nil {
		return "", err
	}
	return strings.Trim(timeOutput, "\""), nil
}

func (r Repository) RefCommitTimestamp(ref string) (int64, error) {
	timeOutput, err := r.cmdOutput("git", "show", "-s", "--format=\"%ct\"", ref)
	if err != nil {
		return 0, err
	}
	timeString := strings.Trim(timeOutput, "\"")
	timestamp, err := strconv.ParseInt(timeString, 10, 64)
	if err != nil {
		return 0, err
	}
	return timestamp, nil
}

func (r Repository) RefMessage(ref string) (string, error) {
	msgOutput, err := r.cmdOutput("git", "show", "-s", "--format=\"%B\"", ref)
	if err != nil {
		return "", err
	}
	return strings.Trim(msgOutput, "\""), nil
}

func (r Repository) LatestRef(branch string) (string, error) {
	refOutput, err := r.cmdOutput("git", "show", "-s", "--format=\"%H\"", branch)
	if err != nil {
		return "", err
	}
	return strings.Trim(refOutput, "\""), nil
}

func (r Repository) RefsSinceTimestamp(branch string, timestamp int64) ([]string, error) {
	refsOutput, err := r.cmdOutput("git", "log", fmt.Sprintf("--since=%d", timestamp), "--format=\"%H\"", branch)
	if err != nil {
		return []string{}, err
	}
	return strings.Split(refsOutput, "\n"), nil
}

func (r Repository) cmdOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = r.Dir
	var outputBytes bytes.Buffer
	cmd.Stdout = &outputBytes
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(outputBytes.String()), nil
}
