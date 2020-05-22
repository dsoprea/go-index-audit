// Package main looks up the given repository's branch's HEAD revision and then
// queries the default global Go Proxy for that module in a loop until it
// reports the latest revision.
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dsoprea/go-logging"
	"github.com/jessevdk/go-flags"

	"github.com/dsoprea/go-index-audit"
)

var (
	mainLogger = log.NewLogger("main.main")
)

type parameters struct {
	TimeoutDuration time.Duration `short:"t" long:"timeout" description:"Alternative timeout duration" default:"1h"`
	PollInterval    time.Duration `short:"i" long:"poll-interval" description:"Alternative pol interval" default:"10s"`
	IsVerbose       bool          `short:"v" long:"verbose" description:"Print logging"`
	ProxyUrl        string        `short:"P" long:"proxy-url" description:"Non-default Proxy URL"`

	Positional struct {
		PackageName string `positional-arg-name:"package-name" description:"Package name"`
	} `positional-args:"yes" required:"yes"`
}

var (
	arguments = new(parameters)
)

func main() {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			log.PrintError(err)

			os.Exit(-2)
		}
	}()

	_, err := flags.Parse(arguments)
	if err != nil {
		os.Exit(-1)
	}

	if arguments.IsVerbose == true {
		cla := log.NewConsoleLogAdapter()
		log.AddAdapter("console", cla)

		scp := log.NewStaticConfigurationProvider()
		scp.SetLevelName(log.LevelNameDebug)

		log.LoadConfiguration(scp)
	}

	packageName := arguments.Positional.PackageName

	p := indexwait.Package{}

	packagePath, vcsName, err := p.GetPackagePath(packageName)
	log.PanicIf(err)

	mainLogger.Debugf(nil, "Package path: [%s]", packagePath)

	vcs, found := indexwait.GetVcs(vcsName)
	if found == false {
		fmt.Printf("VCS [%s] not currently supported.\n", vcsName)
		os.Exit(1)
	}

	currentCommitRevision, currentCommitTimestamp, err := vcs.GetHeadCommit(packagePath)
	log.PanicIf(err)

	mainLogger.Debugf(nil, "Current commit: REVISION=[%s] TIMESTAMP=[%s]", currentCommitRevision, currentCommitTimestamp)

	pc := indexwait.NewProxyClient(arguments.ProxyUrl)

	startTime := time.Now()
	timeoutAt := startTime.Add(arguments.TimeoutDuration)

	for {
		cmi, err := pc.FetchModuleInfo(packageName)
		log.PanicIf(err)

		mv, err := cmi.VersionPhrase.Parse()
		log.PanicIf(err)

		if strings.HasPrefix(currentCommitRevision, mv.RevisionPrefix) == true {
			mainLogger.Infof(nil, "Index matches local revision: INDEX=[%s] LOCAL=[%s]", mv.RevisionPrefix, currentCommitRevision)
			break
		} else if mv.Timestamp.Before(currentCommitTimestamp) != true {
			mainLogger.Infof(nil, "Index now reports a newer or equal revision to local:\n"+
				"INDEX-TIMESTAMP=[%s] INDEX-REVISION=[%s]\n"+
				"LOCAL-TIMESTAMP=[%s] LOCAL-REVISION=[%s]\n",
				mv.Timestamp, mv.RevisionPrefix,
				currentCommitTimestamp, currentCommitRevision)

			break
		}

		if time.Now().After(timeoutAt) == true {
			fmt.Printf("Module has not been updated. Timeout.\n")
			os.Exit(1)
		}

		time.Sleep(arguments.PollInterval)
	}

	stopTime := time.Now()
	duration := stopTime.Sub(startTime)

	mainLogger.Debugf(nil, "Wait time: [%s]", duration)

	// TODO(dustin): We should also be able to call commands with the result.
}
