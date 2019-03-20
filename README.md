# wenex

[![Build Status](https://api.travis-ci.com/nexcode/wenex.svg?branch=master)](https://travis-ci.com/nexcode/wenex)
[![GoDoc](https://godoc.org/github.com/nexcode/wenex?status.svg)](https://godoc.org/github.com/nexcode/wenex)
[![Go Report Card](https://goreportcard.com/badge/github.com/nexcode/wenex)](https://goreportcard.com/report/github.com/nexcode/wenex)

Simple and fast web framework for Go

## Table of Contents

* [Build Status](https://travis-ci.org/nexcode/wenex)
* [GoDoc Reference](https://godoc.org/github.com/nexcode/wenex)
* [Requirements](#requirements)
* [Quick Start](#quick-start)
  * [Download and Install](#download-and-install)
  * [Simple Example](#simple-example)
* [Documentation](#documentation)
  * [Starting the webserver](#starting-the-webserver)

## Requirements

    Go >= 1.8

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

func first(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello,")
	wenex.GetRun(r.Context()).Next()
	io.WriteString(w, "!")
}

func second(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, " World")
}

func main() {
	config := wenex.DefaultConfig()
	config["server.http.listen"] = ":8080"

	wnx, err := wenex.New("simpleapp", config)
	if err != nil {
		panic(err)
	}

	if err = wnx.Router.StrictRoute("/", "HEAD", "GET").Chain(first, second); err != nil {
		panic(err)
	}

	wnx.Logger("info").Print("running application...")

	if err = wnx.Run(); err != nil {
		panic(err)
	}
}
```

Open your browser and visit `http://localhost:8080`

## Documentation

#### Starting the webserver

In its simplest form, a webserver can be started like this:

```go
config := wenex.DefaultConfig()
config["server.http.listen"] = ":8080"

wnx, err := wenex.New("simpleapp", config)
if err != nil {
	panic(err)
}

// define routing and something else...

if err = wnx.Run(); err != nil {
	panic(err)
}
```

In this simple example:
`server.http.listen` - port that will listen to the webserver
`simpleapp` - name of the application (a *simpleapp.conf* file will be created in the working directory)
`config` - configuration parameter map
