# wenex

[![Build Status](https://api.travis-ci.org/nexcode/wenex.svg?branch=master)](https://travis-ci.org/nexcode/wenex)
[![GoDoc](https://godoc.org/github.com/nexcode/wenex?status.svg)](https://godoc.org/github.com/nexcode/wenex)

Simple and fast web framework for Go


## Table of Contents

* [Requirements](#requirements)
* [Quick Start](#quick-start)
  * [Download and Install](#download-and-install)
  * [Simple Example](#simple-example)

## Requirements

    Go >= 1.7

## Quick Start

#### Download and Install

    go get -u github.com/nexcode/wenex

#### Simple Example

```go
package main

import (
	"io"
	"net/http"

	"github.com/nexcode/wenex"
)

const appName = "simpleapp"

func first(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello,")
	wenex.GetRun(r.Context()).Next()
	io.WriteString(w, "!")
}

func second(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, " World")
}

func main() {
	config := wenex.GetDefaultConfig()
	config["log.filePrefix"] = "log/"

	wnx, err := wenex.New(appName, config)
	if err != nil {
		panic(err)
	}

	if err = wnx.Router.Route("HEAD, GET", "/").Chain(first, second); err != nil {
		panic(err)
	}

	wnx.Logger("info").Print("Running " + appName)

	if err = wnx.Run(); err != nil {
		panic(err)
	}
}
```

Open your browser and visit `http://localhost:3000`
