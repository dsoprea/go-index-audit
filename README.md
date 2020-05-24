[![Build Status](https://travis-ci.org/dsoprea/go-index-audit.svg?branch=master)](https://travis-ci.org/dsoprea/go-index-audit)
[![Coverage Status](https://coveralls.io/repos/github/dsoprea/go-index-audit/badge.svg?branch=master)](https://coveralls.io/github/dsoprea/go-index-audit?branch=master)
[![GoDoc](https://godoc.org/github.com/dsoprea/go-index-audit?status.svg)](https://godoc.org/github.com/dsoprea/go-index-audit)
[![Go Report Card](https://goreportcard.com/badge/github.com/dsoprea/go-index-audit)](https://goreportcard.com/report/github.com/dsoprea/go-index-audit)

# Overview

This tool is a workflow optimization for Go module development.

When you are pushing changes to one module and will then need to bump its
revision in the dependencies of another module, you will typically need to wait
between five and forty minutes for the [Go Proxy](https://proxy.golang.org) to
represent the change and for "go get -u" to then see it and be able to update
dependencies with it.

The primary purpoe of this tool is to block until the version reported by the
Index is equal to or newer than your local repository. This also works for
private Proxy implementations/instances.


# Usage

```
$ go run command/go-wait-for-index/main.go github.com/dsoprea/go-exif
$
```

Get more information by printing verbosity:

```
$ go run command/go-wait-for-index/main.go -v github.com/dsoprea/go-exif
2020/05/24 02:15:56 main.main: [DEBUG]  Package path: [/home/doprea/development/go/src/github.com/dsoprea/go-exif]
2020/05/24 02:15:56 main.main: [DEBUG]  Current commit: REVISION=[1a12aec48f9023351e9c08b9bb7dd381e47233e9] TIMESTAMP=[2020-05-20 15:12:04 -0400 -0400]
2020/05/24 02:15:57 main.main: [INFO]  Index matches local revision: INDEX=[1a12aec48f90] LOCAL=[1a12aec48f9023351e9c08b9bb7dd381e47233e9]
2020/05/24 02:15:57 main.main: [DEBUG]  Wait time: [281.851456ms]
```

See the command-line help for more information.
