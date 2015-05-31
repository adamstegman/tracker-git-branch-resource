package out

import (
	"time"

	"github.com/adamstegman/tracker-git-branch-resource"
)

type OutRequest struct {
	Source resource.Source `json:"source"`
	Params struct{}        `json:"params"`
}

type Version struct {
	Time time.Time `json:time`
}

type OutResponse struct {
	Version Version `json:"version"`
}
