package indexwait

import (
	"os"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestGit_Name(t *testing.T) {
	g := new(Git)
	name := g.Name()

	if name != "git" {
		t.Fatalf("Git VCS name not correct: [%s]", name)
	}
}

func TestGit_GetHeadCommit(t *testing.T) {
	originalWd, err := os.Getwd()
	log.PanicIf(err)

	// Get our own source path.

	p := Package{}

	packagePath, _, err := p.GetPackagePath("github.com/dsoprea/go-index-audit")
	log.PanicIf(err)

	// Change path.

	err = os.Chdir(packagePath)
	log.PanicIf(err)

	defer os.Chdir(originalWd)

	// Get revision.

	g := new(Git)

	revision, timestamp, err := g.GetHeadCommit(packagePath)
	log.PanicIf(err)

	// Check all that we can really check.

	if timestamp.IsZero() != false {
		t.Fatalf("Expected project commit timestamp to be non-zero: [%s]", timestamp)
	} else if len(revision) != 40 {
		t.Fatalf("Revision does not look correct: [%s]", revision)
	}
}
