package indexwait

import (
	"fmt"
	"strings"
	"time"

	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dsoprea/go-logging"
)

const (
	defaultProxyUrl = "https://proxy.golang.org"

	// moduleInfoTimestampFormat is the layout string that tells us how to
	// parse the timestamp reported by the Go Proxy [specification].
	//
	// v0.0.0-20200520191204-1a12aec48f90
	moduleInfoTimestampFormat = "20060102150405"
)

// ProxyClient knows how to hit the globalpublic Go Proxy endpoint.
type ProxyClient struct {
	client  *http.Client
	rootUrl string
}

// NewProxyClient returns a new ProxyClient struct.
func NewProxyClient(url string) *ProxyClient {
	client := new(http.Client)

	if url == "" {
		url = defaultProxyUrl
	}

	return &ProxyClient{
		client:  client,
		rootUrl: url,
	}
}

// ModuleVersion describes the components of the reported version.
type ModuleVersion struct {
	// Version is the currently advertised module version.
	Version string

	// Timestamp is the currently advertised committer timestamp.
	Timestamp time.Time

	// RevisionPrefix is the currently advertised revision (just a prefix of the
	// revision, actually).
	RevisionPrefix string
}

// NewModuleVersionFromRaw returns a ModuleVersion instance from the raw phrase.
func NewModuleVersionFromRaw(versionPhrase ModuleVersionPhrase) (mv ModuleVersion, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	// v0.0.0-20200520191204-1a12aec48f90
	parts := strings.Split(string(versionPhrase), "-")

	timestamp, err := time.Parse(moduleInfoTimestampFormat, parts[1])
	log.PanicIf(err)

	mv = ModuleVersion{
		Version:        parts[0],
		Timestamp:      timestamp,
		RevisionPrefix: parts[2],
	}

	return mv, nil
}

// InnerString returns a string that can be embedded in another.
func (mv ModuleVersion) InnerString() string {
	return fmt.Sprintf("MODULE-VERSION=[%s] COMMIT-TIME=[%s] COMMIT-REV-PREFIX=[%s]", mv.Version, mv.Timestamp, mv.RevisionPrefix)
}

// String returns a complete string for printing.
func (mv ModuleVersion) String() string {
	return fmt.Sprintf("ModuleVersion<%s>", mv.InnerString())
}

// DumpFields prints the current values of the fields.
func (mv ModuleVersion) DumpFields() {
	fmt.Printf("Version: [%s]\n", mv.Version)
	fmt.Printf("Commit Timestamp: [%s]\n", mv.Timestamp)
	fmt.Printf("Commit RevisionPrefix: [%s]\n", mv.RevisionPrefix)
}

// ModuleVersionPhrase is a string that knows how to parse itself into
// components.
type ModuleVersionPhrase string

// Parse parses the raw module-version string and returns a ModuleVersion.
func (mvp ModuleVersionPhrase) Parse() (mv ModuleVersion, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	mv, err = NewModuleVersionFromRaw(mvp)
	log.PanicIf(err)

	return mv, nil
}

// CachedModuleInfo is the information currently advertised for a module.
type CachedModuleInfo struct {
	VersionPhrase ModuleVersionPhrase `json:"Version"`
	Timestamp     time.Time           `json:"Time"`
}

// String returns a descriptive string.
func (cmi CachedModuleInfo) String() string {
	mv, err := cmi.VersionPhrase.Parse()
	log.PanicIf(err)

	return fmt.Sprintf("CachedModuleInfo<%s REPORT-TIME=[%s]>", mv.InnerString(), cmi.Timestamp)
}

// DumpFields prints the current values of the fields.
func (cmi CachedModuleInfo) DumpFields() {
	fmt.Printf("Advertised Timestamp: [%s]\n", cmi.Timestamp)

	mv, err := cmi.VersionPhrase.Parse()
	log.PanicIf(err)

	mv.DumpFields()
}

// FetchModuleInfo returns the current published version of the module.
func (pc *ProxyClient) FetchModuleInfo(packageName string) (cmi CachedModuleInfo, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	url := fmt.Sprintf("%s/%s/@latest", pc.rootUrl, packageName)

	response, err := pc.client.Get(url)
	if err != nil {
		fmt.Printf("Could not query proxy: [%s]\n", err.Error())
		log.Panic(err)
	}

	defer response.Body.Close()

	raw, err := ioutil.ReadAll(response.Body)
	log.PanicIf(err)

	b := bytes.NewBuffer(raw)
	jd := json.NewDecoder(b)

	err = jd.Decode(&cmi)
	if err != nil {
		fmt.Printf("Could not parse return: [%s]\n", err.Error())
		fmt.Printf("\n")
		fmt.Printf("OUTPUT:\n")
		fmt.Printf("\n")
		fmt.Printf("%s\n", b.String())

		log.Panic(err)
	}

	return cmi, nil
}
