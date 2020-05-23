package indexwait

import (
	"testing"
)

func TestRegisterVcs(t *testing.T) {
	backupVcs := vcsRegistry
	vcsRegistry = make(map[string]Vcs)

	defer func() {
		vcsRegistry = backupVcs
	}()

	g := new(Git)
	registerVcs(g)

	if len(vcsRegistry) != 1 {
		t.Fatalf("Expected one VCS to be registered: (%d)", len(vcsRegistry))
	}
}

func TestGetVcs_Hit(t *testing.T) {
	backupVcs := vcsRegistry
	vcsRegistry = make(map[string]Vcs)

	defer func() {
		vcsRegistry = backupVcs
	}()

	g := new(Git)
	registerVcs(g)

	vcs, found := GetVcs("git")
	if found != true {
		t.Fatalf("Could not find Git as a registered VCS.")
	} else if _, ok := vcs.(*Git); ok != true {
		t.Fatalf("VCS getter did not return a Git struct.")
	}
}

func TestGetVcs_Miss(t *testing.T) {
	backupVcs := vcsRegistry
	vcsRegistry = make(map[string]Vcs)

	defer func() {
		vcsRegistry = backupVcs
	}()

	_, found := GetVcs("invalid")
	if found != false {
		t.Fatalf("Expected a miss for an invalid VCS.")
	}
}
