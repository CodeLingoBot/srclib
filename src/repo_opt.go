package src

import (
	"log"
	"os"
	"path/filepath"

	"sourcegraph.com/sourcegraph/go-flags"
)

var (
	localRepo    *Repo
	localRepoErr error

	// CacheLocalRepo controls whether openLocalRepo caches the
	// results of OpenRepo the first time it runs and returns the same
	// repo for all subsequent calls (even if you call os.Chdir, for
	// example).
	CacheLocalRepo = true
)

// openLocalRepo opens the VCS repository in or above the current
// directory.
func openLocalRepo() (*Repo, error) {
	if !CacheLocalRepo {
		return OpenRepo(".")
	}

	// Only try to open the current-dir repo once (we'd get the same result each
	// time, since we never modify it).
	if localRepo == nil && localRepoErr == nil {
		localRepo, localRepoErr = OpenRepo(".")
	}
	return localRepo, localRepoErr
}

func setDefaultRepoURIOpt(c *flags.Command) {
	openLocalRepo()
	if localRepo != nil {
		if localRepo.CloneURL != "" {
			SetOptionDefaultValue(c.Group, "repo", localRepo.URI())
		}
	}
}

func setDefaultCommitIDOpt(c *flags.Command) {
	openLocalRepo()
	if localRepo != nil {
		if localRepo.CommitID != "" {
			SetOptionDefaultValue(c.Group, "commit", localRepo.CommitID)
		}
	}
}

func setDefaultRepoSubdirOpt(c *flags.Command) {
	openLocalRepo()
	if localRepo != nil {
		absDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		subdir, _ := filepath.Rel(localRepo.RootDir, absDir)
		SetOptionDefaultValue(c.Group, "subdir", subdir)
	}
}
