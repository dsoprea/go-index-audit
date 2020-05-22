package indexwait

import (
	"strings"

	"go/build"

	"github.com/dsoprea/go-logging"
	"golang.org/x/tools/go/vcs"
)

// Package manages project package information.
type Package struct {
}

// GetPackagePath returns the giben package's
func (Package) GetPackagePath(packageName string) (packagePath, vcsName string, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	p, err := build.Default.Import(
		packageName,
		build.Default.GOPATH,
		build.FindOnly)

	log.PanicIf(err)

	packagePath = p.Dir

	rr, err := vcs.RepoRootForImportPath(packageName, false)
	log.PanicIf(err)

	return packagePath, strings.ToLower(rr.VCS.Name), nil
}
