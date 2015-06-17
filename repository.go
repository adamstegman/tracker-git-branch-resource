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
	var err error
	if r.Dir == "" {
		r.Dir, err = ioutil.TempDir("", "tracker-git-branch-resource-repository")
		if err != nil {
			return fmt.Errorf("Could not create temporary directory: %s", err)
		}
	}
	cloneCmd := exec.Command("git", "clone", r.source, r.Dir)
	err = cloneCmd.Run()
	if err != nil {
		return fmt.Errorf("Could not clone repository %s into directory %s: %s", r.source, r.Dir, err)
	}
	return nil
}

func (r Repository) Fetch() error {
	fetchCmd := exec.Command("git", "fetch", "origin")
	fetchCmd.Dir = r.Dir
	err := fetchCmd.Run()
	if err != nil {
		return fmt.Errorf("Could not fetch origin: %s", err)
	}
	return nil
}

func (r Repository) CheckoutRef(ref string) error {
	checkoutCmd := exec.Command("git", "checkout", ref)
	checkoutCmd.Dir = r.Dir
	err := checkoutCmd.Run()
	if err != nil {
		return fmt.Errorf("Could not checkout %s: %s", ref, err)
	}
	return nil
}

func (r Repository) RemoteBranches() ([]string, error) {
	branchesOutput, err := r.cmdOutput("git", "branch", "-r")
	if err != nil {
		return []string{}, fmt.Errorf("Could not list remote branches: %s", err)
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
		return "", fmt.Errorf("Could not show author name for %s: %s", ref, err)
	}
	return strings.Trim(nameOutput, "\""), nil
}

func (r Repository) RefAuthorDate(ref string) (string, error) {
	timeOutput, err := r.cmdOutput("git", "show", "-s", "--format=\"%ai\"", ref)
	if err != nil {
		return "", fmt.Errorf("Could not show author date for %s: %s", ref, err)
	}
	return strings.Trim(timeOutput, "\""), nil
}

func (r Repository) RefCommitName(ref string) (string, error) {
	nameOutput, err := r.cmdOutput("git", "show", "-s", "--format=\"%cn\"", ref)
	if err != nil {
		return "", fmt.Errorf("Could not show committer name for %s: %s", ref, err)
	}
	return strings.Trim(nameOutput, "\""), nil
}

func (r Repository) RefCommitDate(ref string) (string, error) {
	timeOutput, err := r.cmdOutput("git", "show", "-s", "--format=\"%ci\"", ref)
	if err != nil {
		return "", fmt.Errorf("Could not show committer date for %s: %s", ref, err)
	}
	return strings.Trim(timeOutput, "\""), nil
}

func (r Repository) RefCommitTimestamp(ref string) (int64, error) {
	timeOutput, err := r.cmdOutput("git", "show", "-s", "--format=\"%ct\"", ref)
	if err != nil {
		return 0, fmt.Errorf("Could not show committer timestamp for %s: %s", ref, err)
	}
	timeString := strings.Trim(timeOutput, "\"")
	timestamp, err := strconv.ParseInt(timeString, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Could not parse committer timestamp (%s) for %s: %s", timeString, ref, err)
	}
	return timestamp, nil
}

func (r Repository) RefMessage(ref string) (string, error) {
	msgOutput, err := r.cmdOutput("git", "show", "-s", "--format=\"%B\"", ref)
	if err != nil {
		return "", fmt.Errorf("Could not show message for %s: %s", ref, err)
	}
	return strings.Trim(msgOutput, "\""), nil
}

func (r Repository) LatestRef(branch string) (string, error) {
	refOutput, err := r.cmdOutput("git", "show", "-s", "--format=\"%H\"", branch)
	if err != nil {
		return "", fmt.Errorf("Could not show SHA for %s: %s", branch, err)
	}
	return strings.Trim(refOutput, "\""), nil
}

func (r Repository) RefsSinceTimestamp(branch string, timestamp int64) ([]string, error) {
	refsOutput, err := r.cmdOutput("git", "log", fmt.Sprintf("--since=%d", timestamp), "--format=\"%H\"", branch)
	if err != nil {
		return []string{}, fmt.Errorf("Could not list refs since %d for %s: %s", timestamp, branch, err)
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
		return "", fmt.Errorf("Could not run command %s with args %v: %s", name, args, err)
	}
	return strings.TrimSpace(outputBytes.String()), nil
}
