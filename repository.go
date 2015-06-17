package resource

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Repository struct {
	dir    string
	source string
}

func NewRepository(source string, dir string) Repository {
	return Repository{
		dir:    dir,
		source: source,
	}
}

func (r Repository) Clone() error {
	err := r.runCmd("git", "clone", r.source, r.dir)
	if err != nil {
		return fmt.Errorf("Could not clone repository %s into directory %s: %s", r.source, r.dir, err)
	}
	return nil
}

func (r Repository) Fetch() error {
	err := r.runCmd("git", "fetch", "origin")
	if err != nil {
		return fmt.Errorf("Could not fetch origin: %s", err)
	}
	return nil
}

func (r Repository) CheckoutRef(ref string) error {
	err := r.runCmd("git", "checkout", ref)
	if err != nil {
		return fmt.Errorf("Could not checkout %s: %s", ref, err)
	}
	return nil
}

func (r Repository) RemoteBranches() ([]string, error) {
	branchesOutput, err := r.runCmdOutput("git", "branch", "-r")
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
	nameOutput, err := r.runCmdOutput("git", "show", "-s", "--format=\"%an\"", ref)
	if err != nil {
		return "", fmt.Errorf("Could not show author name for %s: %s", ref, err)
	}
	return strings.Trim(nameOutput, "\""), nil
}

func (r Repository) RefAuthorDate(ref string) (string, error) {
	timeOutput, err := r.runCmdOutput("git", "show", "-s", "--format=\"%ai\"", ref)
	if err != nil {
		return "", fmt.Errorf("Could not show author date for %s: %s", ref, err)
	}
	return strings.Trim(timeOutput, "\""), nil
}

func (r Repository) RefCommitName(ref string) (string, error) {
	nameOutput, err := r.runCmdOutput("git", "show", "-s", "--format=\"%cn\"", ref)
	if err != nil {
		return "", fmt.Errorf("Could not show committer name for %s: %s", ref, err)
	}
	return strings.Trim(nameOutput, "\""), nil
}

func (r Repository) RefCommitDate(ref string) (string, error) {
	timeOutput, err := r.runCmdOutput("git", "show", "-s", "--format=\"%ci\"", ref)
	if err != nil {
		return "", fmt.Errorf("Could not show committer date for %s: %s", ref, err)
	}
	return strings.Trim(timeOutput, "\""), nil
}

func (r Repository) RefCommitTimestamp(ref string) (int64, error) {
	timeOutput, err := r.runCmdOutput("git", "show", "-s", "--format=\"%ct\"", ref)
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
	msgOutput, err := r.runCmdOutput("git", "show", "-s", "--format=\"%B\"", ref)
	if err != nil {
		return "", fmt.Errorf("Could not show message for %s: %s", ref, err)
	}
	return strings.Trim(msgOutput, "\""), nil
}

func (r Repository) LatestRef(branch string) (string, error) {
	refOutput, err := r.runCmdOutput("git", "show", "-s", "--format=\"%H\"", branch)
	if err != nil {
		return "", fmt.Errorf("Could not show SHA for %s: %s", branch, err)
	}
	return strings.Trim(refOutput, "\""), nil
}

func (r Repository) RefsSinceTimestamp(branch string, timestamp int64) ([]string, error) {
	refsOutput, err := r.runCmdOutput("git", "log", fmt.Sprintf("--since=%d", timestamp), "--format=\"%H\"", branch)
	if err != nil {
		return []string{}, fmt.Errorf("Could not list refs since %d for %s: %s", timestamp, branch, err)
	}
	return strings.Split(refsOutput, "\n"), nil
}

func (r Repository) runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = r.dir
	var errBytes bytes.Buffer
	cmd.Stderr = &errBytes
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s %v in %s failed: %s\n[STDERR]\n%s", name, args, r.dir, err, errBytes.String())
	}
	return nil
}

func (r Repository) runCmdOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = r.dir
	var outputBytes bytes.Buffer
	cmd.Stdout = &outputBytes
	var errBytes bytes.Buffer
	cmd.Stderr = &errBytes
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%s %v in %s failed: %s\n[STDERR]\n%s", name, args, r.dir, err, errBytes.String())
	}
	return strings.TrimSpace(outputBytes.String()), nil
}
