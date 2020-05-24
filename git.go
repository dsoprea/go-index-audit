package indexwait

import (
	"time"

	"github.com/dsoprea/go-logging"
	"github.com/go-git/go-git/v5"
)

// Git provides a VCS implementation for Git.
type Git struct {
}

// Name returns the name of the VCS method.
func (*Git) Name() string {
	return "git"
}

// GetHeadCommit returns the current revision of the package/module.
func (*Git) GetHeadCommit(packagePath string) (revision string, timestamp time.Time, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	gr, err := git.PlainOpen(packagePath)
	log.PanicIf(err)

	lo := new(git.LogOptions)
	ci, err := gr.Log(lo)
	log.PanicIf(err)

	commit, err := ci.Next()
	log.PanicIf(err)

	revision = commit.Hash.String()
	timestamp = commit.Committer.When

	return revision, timestamp, nil
}

func init() {
	registerVcs(new(Git))
}
