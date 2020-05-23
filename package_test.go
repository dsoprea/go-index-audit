package indexwait

import (
	"strings"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestPackage_GetPackagePath(t *testing.T) {
	p := Package{}

	packageName := "github.com/dsoprea/go-index-audit"

	packagePath, vcsName, err := p.GetPackagePath(packageName)
	log.PanicIf(err)

	if strings.HasSuffix(packagePath, packageName) != true {
		t.Fatalf("Package-path does not appear to be correct: [%s]", packagePath)
	}

	if vcsName != "git" {
		t.Fatalf("Reported VCS not correct: [%s]", vcsName)
	}
}
