package indexwait

import (
	"testing"
	"time"

	"github.com/dsoprea/go-logging"
)

const (
	testModuleName = "github.com/dsoprea/go-exif"
)

func TestNewModuleVersionFromRaw(t *testing.T) {
	mv, err := NewModuleVersionFromRaw("v0.0.0-20200520191204-1a12aec48f90")
	log.PanicIf(err)

	s := mv.String()
	if s != "ModuleVersion<MODULE-VERSION=[v0.0.0] COMMIT-TIME=[2020-05-20 19:12:04 +0000 UTC] COMMIT-REV-PREFIX=[1a12aec48f90]>" {
		t.Fatalf("String not correct: [%s]", s)
	}
}

func TestModuleVersion_InnerString(t *testing.T) {
	mv := ModuleVersion{}
	s := mv.InnerString()
	if s != "MODULE-VERSION=[] COMMIT-TIME=[0001-01-01 00:00:00 +0000 UTC] COMMIT-REV-PREFIX=[]" {
		t.Fatalf("InnerString() does not return the correct string: [%s]", s)
	}
}

func TestModuleVersion_String(t *testing.T) {
	mv := ModuleVersion{}
	s := mv.String()
	if s != "ModuleVersion<MODULE-VERSION=[] COMMIT-TIME=[0001-01-01 00:00:00 +0000 UTC] COMMIT-REV-PREFIX=[]>" {
		t.Fatalf("String() does not return the correct string: [%s]", s)
	}
}

func TestModuleVersion_DumpFields(t *testing.T) {
	mv := ModuleVersion{}
	mv.DumpFields()
}

func TestModuleVersionPhrase_Parse(t *testing.T) {
	mvp := ModuleVersionPhrase("v0.0.0-20200520191204-1a12aec48f90")

	mv, err := mvp.Parse()
	log.PanicIf(err)

	s := mv.String()
	if s != "ModuleVersion<MODULE-VERSION=[v0.0.0] COMMIT-TIME=[2020-05-20 19:12:04 +0000 UTC] COMMIT-REV-PREFIX=[1a12aec48f90]>" {
		t.Fatalf("String not correct: [%s]", s)
	}
}

func TestCachedModuleInfo_String(t *testing.T) {
	cmi := CachedModuleInfo{
		VersionPhrase: "v0.0.0-20200520191204-1a12aec48f90",
		Timestamp:     time.Time{},
	}

	s := cmi.String()
	if s != "CachedModuleInfo<MODULE-VERSION=[v0.0.0] COMMIT-TIME=[2020-05-20 19:12:04 +0000 UTC] COMMIT-REV-PREFIX=[1a12aec48f90] REPORT-TIME=[0001-01-01 00:00:00 +0000 UTC]>" {
		t.Fatalf("String not correct: [%s]", s)
	}
}

func TestCachedModuleInfo_DumpFields(t *testing.T) {
	cmi := CachedModuleInfo{
		VersionPhrase: "v0.0.0-20200520191204-1a12aec48f90",
		Timestamp:     time.Time{},
	}

	cmi.DumpFields()
}

func TestNewProxyClient(t *testing.T) {
	pc := NewProxyClient("")
	if pc.client == nil {
		t.Fatalf("'client' field not set.")
	}
}

func TestProxyClient_FetchModuleInfo(t *testing.T) {
	pc := NewProxyClient("")

	cmi, err := pc.FetchModuleInfo(testModuleName)
	log.PanicIf(err)

	s := cmi.String()
	if s != "CachedModuleInfo<MODULE-VERSION=[v0.0.0] COMMIT-TIME=[2020-05-20 19:12:04 +0000 UTC] COMMIT-REV-PREFIX=[1a12aec48f90] REPORT-TIME=[2020-05-20 19:12:04 +0000 UTC]>" {
		t.Fatalf("String not correct: [%s]", s)
	}
}
