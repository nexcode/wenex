# wenex

[![Build Status](https://api.travis-ci.org/nexcode/wenex.svg?branch=master)](https://travis-ci.org/nexcode/wenex)
[![GoDoc](https://godoc.org/github.com/nexcode/wenex?status.svg)](https://godoc.org/github.com/nexcode/wenex)

Simple and fast web framework for Go

## Table of Contents

* [Build Status](https://travis-ci.org/nexcode/wenex)
* [GoDoc Reference](https://godoc.org/github.com/nexcode/wenex)
* [Requirements](#requirements)
* [Quick Start](#quick-start)
  * [Download and Install](#download-and-install)
  * [Simple Example](#simple-example)

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
	config := wenex.GetDefaultConfig()
	config["server.http.listen"] = ":8080"

	wnx, err := wenex.New("simpleApp", config)
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
