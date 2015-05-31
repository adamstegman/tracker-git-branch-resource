package check

import "github.com/adamstegman/tracker-git-branch-resource"

type Request struct {
	Source  resource.Source  `json:"source"`
	Version resource.Version `json:"version"`
}
