package out

import "github.com/adamstegman/tracker-git-branch-resource"

type OutRequest struct {
	Source resource.Source `json:"source"`
	Params Params          `json:"params"`
}

type Params struct {
	Repos []string `json:"repos"`
}

type OutResponse struct {
	Version resource.Version `json:"version"`
}
