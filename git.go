package indexwait

import (
	"os"
	"path"
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

	repositoryPath := packagePath

	// Search from current directory up to find the repository root.
	for repositoryPath != "/" && repositoryPath != "." {
		metaPath := path.Join(repositoryPath, ".git")

		f, err := os.Open(metaPath)
		f.Close()

		if err == nil {
			break
		} else if os.IsNotExist(err) != true {
			log.Panic(err)
		}

		repositoryPath = path.Dir(repositoryPath)
	}

	gr, err := git.PlainOpen(repositoryPath)
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
