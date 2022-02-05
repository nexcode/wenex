# wenex

[![GoDoc](https://godoc.org/github.com/nexcode/wenex?status.svg)](https://pkg.go.dev/github.com/nexcode/wenex)
[![Go Report Card](https://goreportcard.com/badge/github.com/nexcode/wenex)](https://goreportcard.com/report/github.com/nexcode/wenex)

Simple and fast web framework for Go

## Table of Contents

* [Build Status](https://travis-ci.org/nexcode/wenex)
* [GoDoc Reference](https://godoc.org/github.com/nexcode/wenex)
* [Requirements](#requirements)
* [Quick Start](#quick-start)
  * [Download and Install](#download-and-install)
  * [Simple Example](#simple-example)
* [Starting the webserver](#starting-the-webserver)
* [Configuration options](#configuration-options)
* [Work with config file](#work-with-config-file)
* [Routing configuration](#routing-configuration)
* [Work with logger](#work-with-logger)

## Requirements

    Go >= 1.17

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

	wnx, err := wenex.New("simpleapp", config, nil)
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

## Starting the webserver

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

## Configuration options

* `server.gzip.enable` - enables and disables gzip
* `server.gzip.level` - gzip compression level
* `server.http.listen` - port that will listen to http traffic
* `server.https.listen` - port that will listen to TLS (https) traffic
* `server.https.stringCert.cert` - string containing certificate
* `server.https.stringCert.key` - string containing private key
* `server.https.loadCert.cert` - file containing certificate
* `server.https.loadCert.key` - file containing private key
* `server.https.autoCert.hosts` - array of domains
* `server.https.autoCert.dirCache` - cache directory
* `server.timeout.read` - connection read timeout
* `server.timeout.write` - connection write timeout
* `server.timeout.idle` - connection idle timeout
* `logger.defaultName` - log filename for empty logger
* `logger.namePrefix` - prefix that will be added to all saved log files.
	For example, if you use `log/` prefix, then all logs files will be in `log/` folder
* `logger.useFlag` - sets the output flags for the logger. The flag bits are Ldate, Ltime, and so on
* `logger.usePrefix` - string that will be added at the beginning of each message

## Work with config file

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

## Routing configuration

For the routing declaration in wenex two methods are used:
* `wnx.Router.StrictRoute(pattern, methods)` - tied to the end of pattern
* `wnx.Router.WeakRoute(pattern, methods)` - not tied to the end of pattern

Wenex supports the following special constructs in the pattern:
* `*` - a sequence of any characters, including the empty string
* `:name` - value of this path element will be available as a value of a get-variable with the same name

Routing declaration returns a method, that allows you to specify multiple handlers:  
`wnx.Router.StrictRoute(pattern, methods).Chain(handler1, handler2, handlerN)`

Matching examples:
```go
wnx.Router.StrictRoute("/*/:var/test/", "HEAD", "GET").Chain(...)
// matching requests:
// /sefsef/aaa/test/
// /zzz/qwe/test/

wnx.Router.WeakRoute("/*/:var/test/", "HEAD", "GET").Chain(...)
// matching requests:
// /sefsef/aaa/test/
// /zzz/zxc/test/rrr/
// /zzz/gg/test/ppp/fff
```

Chains can run completely sequentially, or you can call the next chain before the first one has completed.  
For this, the `Next()` method is used.  
An example of this behavior is given in the section [Simple Example](#simple-example).

## Work with logger

Winx creates files with logs dynamically.  
It use `log.filePrefix` fo path prefix fo all logs files.  
For example:
```go
wnx.Logger("file1").Print("some data...")
wnx.Logger("folder2/file2").Print("some data...")

// default log file:
wnx.Logger("").Print("some data...")
```
You can customize the logger in accordance with the std `log.logger` api:
```go
wnx.Logger("").SetPrefix("prefix")
```
