package main

import (
	"os"
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/testing"
)

func TestMain(t *testing.T) {
	// This test will block until the index is at least as up to date as the
	// local repository.

	ritesting.RedirectTty()
	defer ritesting.RestoreAndDumpTty()

	ritesting.EnableMarshaledExits()
	defer ritesting.DisableMarshaledExits()

	originalArgs := os.Args

	originalArguments := arguments
	arguments = new(parameters)

	defer func() {
		ritesting.DisableMarshaledExits()

		if errRaw := recover(); errRaw != nil {
			ritesting.RestoreAndDumpTty()

			err := errRaw.(error)
			log.Panic(err)
		} else {
			ritesting.RestoreTty()
		}

		os.Args = originalArgs
		arguments = originalArguments
	}()

	os.Args = []string{
		"",
		"github.com/dsoprea/go-index-audit",
	}

	main()
}
