package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/adamstegman/tracker-git-branch-resource/out"
)

func main() {
	if len(os.Args) < 2 {
		sayf("usage: %s <sources directory>\n", os.Args[0])
		os.Exit(1)
	}

	json.NewEncoder(os.Stdout).Encode(out.OutResponse{
		Version: out.Version{
			Time: time.Now(),
		},
	})
}

func sayf(message string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, message, args...)
}
