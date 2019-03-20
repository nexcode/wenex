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
  * [Configuration options](#configuration-options)
  * [Work with config file](#work-with-config-file)

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
`simpleapp` - name of the application (a `simpleapp.conf` file will be created in the working directory)  
`config` - configuration options

#### Configuration options

* `server.http.listen` - port that will listen to http traffic
* `server.https.listen` - port that will listen to TLS (https) traffic
* `server.https.crt` - TLS certificate
* `server.https.key` - TLS private key
* `server.timeout.read` - connection read timeout
* `server.timeout.write` - connection write timeout
* `server.timeout.idle` - connection idle timeout
* `log.filePrefix` - prefix that will be added to all saved log files.
	For example, if you use `log/` prefix, then all logs files will be in `log/` folder

#### Work with config file

If you run wenex with this config:
```go
config := wenex.DefaultConfig()

wnx, err := wenex.New("simpleapp", config)
if err != nil {
	panic(err)
}

// Some code and wnx.Run()
```

A config file (`simpleapp.conf`) appears in the working directory.  
From the `config` variable, only missing values will be written to the file.  
Overwriting existing values will not occur.
```json
{
    "log": {
        "filePrefix": "log/"
    },
    "server": {
        "http": {
            "listen": ":http"
        },
        "timeout": {
            "idle": "30s",
            "read": "30s",
            "write": "30s"
        }
    }
}
```

You can add any parameters directly to the file or use api:
```go
wnx.Config.Set("key1.key2.keyN", 1000)
err := wnx.Config.Save()
```

After this, the config file will look like this:
```json
{
    "key1": {
        "key2": {
            "keyN": 1000
        }
    },
    "log": {
        "filePrefix": "log/"
    },
    "server": {
        "http": {
            "listen": ":http"
        },
        "timeout": {
            "idle": "30s",
            "read": "30s",
            "write": "30s"
        }
    }
}
```

You can get the value of the parameters by api:
```go
valueF64, err := wnx.Config.Float64("key1.key2.keyN")
// Or use it (panic on type error):
// value := wnx.Config.MustFloat64("key1.key2.keyN")

valueStr, err := wnx.Config.String("server.http.listen")
// Or use it (panic on type error):
// value := wnx.Config.MustString("server.http.listen")

// You can get the value as an interface{}
valueInterface := wnx.Config.Get("key")
```
