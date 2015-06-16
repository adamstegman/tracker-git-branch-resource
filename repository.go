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

func NewRepository(source string) Repository {
	return Repository{source: source}
}

func (r Repository) Clone() error {
	var err error
	r.Dir, err = ioutil.TempDir("", "tracker-git-branch-resource-repository")
	if err != nil {
		return err
	}
	cloneCmd := exec.Command("git", "clone", r.source, r.Dir)
	return cloneCmd.Run()
}

func (r Repository) FetchBranches() error {
	fetchCmd := exec.Command("git", "fetch", "origin")
	fetchCmd.Dir = r.Dir
	return fetchCmd.Run()
}

func (r Repository) RemoteBranches() ([]string, error) {
	branchesOutput, err := r.cmdOutput("git", "branch", "-r")
	if err != nil {
		return []string{}, err
	}
	branches := strings.Split(branchesOutput, "\n")

	trimmedBranches := []string{}
	for _, branch := range branches {
		trimmedBranches = append(trimmedBranches, strings.TrimSpace(branch))
	}
	return trimmedBranches, nil
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
