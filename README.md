# brouter
[![Go](https://github.com/antlabs/brouter/workflows/Go/badge.svg)](https://github.com/antlabs/brouter/actions)
[![codecov](https://codecov.io/gh/antlabs/brouter/branch/main/graph/badge.svg)](https://codecov.io/gh/antlabs/brouter)

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

## 测试结果
* httprouter 1.3
* brouter 0.0.1
```
GithubAPI Routes: 180
GithubAPI2 Routes: 203
   BeegoMuxRouter: 97024 Bytes
   BoneRouter: 86368 Bytes
   ChiRouter: 64584 Bytes
   HttpRouter: 35360 Bytes
   BRouter: 51096 Bytes
   trie-mux: 121736 Bytes
   MuxRouter: 1373064 Bytes
   GoRouter1: 76528 Bytes
   GoRouter2: 84624 Bytes
#Static Routes: 157
   HttpRouter: 21680 Bytes
   BRouter: 35208 Bytes

goos: linux
goarch: amd64
pkg: test
BenchmarkBeegoMuxRouterWithGithubAPI-4   	   10000	    111713 ns/op	  116352 B/op	     900 allocs/op
BenchmarkBoneRouterWithGithubAPI-4       	     721	   1546779 ns/op	  562878 B/op	    6807 allocs/op
BenchmarkTrieMuxRouterWithGithubAPI-4    	   19060	     64217 ns/op	   57024 B/op	     468 allocs/op
BenchmarkBRouterWithGithubAPI-4          	   70467	     16971 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpRouterWithGithubAPI-4       	   45830	     26609 ns/op	   11744 B/op	     144 allocs/op
BenchmarkGoRouter1WithGithubAPI-4        	   24042	     49632 ns/op	   11920 B/op	     360 allocs/op
BenchmarkGoRouter2WithGithubAPI2-4       	   21088	     57408 ns/op	   13832 B/op	     406 allocs/op
BenchmarkChiRouterWithGithubAPI2-4       	    8469	    128427 ns/op	  106000 B/op	    1110 allocs/op
BenchmarkMuxRouterWithGithubAPI2-4       	     358	   3414899 ns/op	   59373 B/op	     992 allocs/op
BenchmarkHttpRouter_StaticAll-4          	  120054	      9994 ns/op	       0 B/op	       0 allocs/op
BenchmarkBRouter_StaticAll-4             	  111614	     10838 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	test	15.950s
```

* httprouter fe77dd05ab5a80f54110cccf1b7d8681c2648323(在测试这版本时候遇到过panic)
* brouter 0.0.1
```
GithubAPI Routes: 180
GithubAPI2 Routes: 203
   BeegoMuxRouter: 97232 Bytes
   BoneRouter: 88296 Bytes
   ChiRouter: 64592 Bytes
   HttpRouter: 33480 Bytes
   BRouter: 51096 Bytes
   trie-mux: 123448 Bytes
   MuxRouter: 1373064 Bytes
   GoRouter1: 76320 Bytes
   GoRouter2: 85040 Bytes
#Static Routes: 157
   HttpRouter: 21712 Bytes
   BRouter: 35200 Bytes

goos: linux
goarch: amd64
pkg: test
BenchmarkBeegoMuxRouterWithGithubAPI-4   	   10000	    111217 ns/op	  116352 B/op	     900 allocs/op
BenchmarkBoneRouterWithGithubAPI-4       	     780	   1554833 ns/op	  562871 B/op	    6807 allocs/op
BenchmarkTrieMuxRouterWithGithubAPI-4    	   19020	     62882 ns/op	   57024 B/op	     468 allocs/op
BenchmarkBRouterWithGithubAPI-4          	   69400	     17426 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpRouterWithGithubAPI-4       	   75178	     15907 ns/op	       0 B/op	       0 allocs/op
BenchmarkGoRouter1WithGithubAPI-4        	   24426	     49969 ns/op	   11920 B/op	     360 allocs/op
BenchmarkGoRouter2WithGithubAPI2-4       	   21087	     56849 ns/op	   13832 B/op	     406 allocs/op
BenchmarkChiRouterWithGithubAPI2-4       	    8490	    129326 ns/op	  105974 B/op	    1110 allocs/op
BenchmarkMuxRouterWithGithubAPI2-4       	     349	   3383843 ns/op	   59369 B/op	     992 allocs/op
BenchmarkHttpRouter_StaticAll-4          	  120390	     10094 ns/op	       0 B/op	       0 allocs/op
BenchmarkBRouter_StaticAll-4             	  113125	     10633 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	test	15.887s

```

## 测试代码位置
https://github.com/junelabs/brouter-benchmark

