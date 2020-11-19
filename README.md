# brouter
[![Go](https://github.com/antlabs/brouter/workflows/Go/badge.svg)](https://github.com/antlabs/brouter/actions)
[![codecov](https://codecov.io/gh/antlabs/brouter/branch/master/graph/badge.svg)](https://codecov.io/gh/antlabs/brouter)

项目开始时只是一个新的尝试，看能否性能比httprouter。
## demo
```go
package main

import (
    "fmt"
    "net/http"
    "log"

    "github.com/antlabs/brouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ brouter.Params) {
    fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps brouter.Params) {
    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main() {
    router := brouter.New()
    router.GET("/", Index)
    router.GET("/hello/:name", Hello)

    log.Fatal(http.ListenAndServe(":8080", router))
}
```

## benchmark

## benchmark
