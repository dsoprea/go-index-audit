package indexwait

import (
	"strings"
	"time"
)

type Vcs interface {
	// Name returns the name of the VCS method. Must be in the family of names
	// returned by go-git (though it will be compared case-insensitively).
	Name() string

	// GetHeadCommit returns the current revision of the packaging/module.
	GetHeadCommit(packagePath string) (revision string, timestamp time.Time, err error)
}

var (
	vcsRegistry = make(map[string]Vcs)
)

func registerVcs(vcs Vcs) {
	name := vcs.Name()
	name = strings.ToLower(name)

	vcsRegistry[name] = vcs
}

func GetVcs(name string) (vcs Vcs, found bool) {
	name = strings.ToLower(name)

	vcs, found = vcsRegistry[name]
	return vcs, found
}
