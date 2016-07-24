# httpreq

The package provides an easy way to "unmarshal" query string data into a struct. Without reflect.

[![GoDoc](https://godoc.org/github.com/creack/httpreq?status.svg)](https://godoc.org/github.com/creack/httpreq) [![Build Status](https://travis-ci.org/creack/httpreq.svg)](https://travis-ci.org/creack/httpreq) [![Coverage Status](https://coveralls.io/repos/github/creack/httpreq/badge.svg?branch=master)](https://coveralls.io/github/creack/httpreq?branch=master)

# Example

## Literal

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/creack/httpreq"
)

// Req is the request query struct.
type Req struct {
	Fields    []string
	Limit     int
	Page      int
	Timestamp time.Time
}

func handler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data := &Req{}
	if err := (httpreq.ParsingMap{
		{Field: "limit", Fct: httpreq.ToInt, Dest: &data.Limit},
		{Field: "page", Fct: httpreq.ToInt, Dest: &data.Page},
		{Field: "fields", Fct: httpreq.ToCommaList, Dest: &data.Fields},
		{Field: "timestamp", Fct: httpreq.ToTSTime, Dest: &data.Timestamp},
	}.Parse(req.Form)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_ = json.NewEncoder(w).Encode(data)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
	// curl 'http://localhost:8080?timestamp=1437743020&limit=10&page=1&fields=a,b,c'
}
```

## Chained

```go
package main

import (
        "encoding/json"
        "log"
        "net/http"
        "time"

        "github.com/creack/httpreq"
)

// Req is the request query struct.
type Req struct {
        Fields    []string
        Limit     int
        Page      int
        Timestamp time.Time
}

func handler(w http.ResponseWriter, req *http.Request) {
        if err := req.ParseForm(); err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
        }

        data := &Req{}
        if err := httpreq.NewParsingMap().
                Add("limit", httpreq.ToInt, &data.Limit).
                Add("page", httpreq.ToInt, &data.Page).
                Add("fields", httpreq.ToCommaList, &data.Fields).
                Add("timestamp", httpreq.ToTSTime, &data.Timestamp).
                Parse(req.Form); err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
        }

        _ = json.NewEncoder(w).Encode(data)
}

func main() {
        http.HandleFunc("/", handler)
        log.Fatal(http.ListenAndServe(":8080", nil))
        // curl 'http://localhost:8080?timestamp=1437743020&limit=10&page=1&fields=a,b,c'
}
```

# Benchmarks

## Single CPU

```
BenchmarkRawLiteral           5000000        410 ns/op       96 B/op        2 allocs/op
BenchmarkRawAdd               1000000       1094 ns/op      384 B/op        5 allocs/op
BenchmarkRawJSONUnmarshal      500000       3038 ns/op      416 B/op       11 allocs/op
```

## 8 CPUs

```
BenchmarkRawPLiteral-8        5000000        299 ns/op       96 B/op        2 allocs/op
BenchmarkRawPAdd-8            2000000        766 ns/op      384 B/op        5 allocs/op
BenchmarkRawPJSONUnmarshal-8  1000000       1861 ns/op      416 B/op       11 allocs/op
```
