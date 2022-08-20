# Gin Request Id
[Gin](https://github.com/gin-gonic/gin) middleware to add a missing `X-Request-ID` header.

## How to use
    import (
        "github.com/Cehir/ginrequestid"
    )
    
    cfg := ginrequestid.DefaultCfg()

    router := gin.New()
	router.Use(
        ginrequestid.RequestID(cfg)
    )

## Benchmark 

### Setup
    go install golang.org/x/perf/cmd/benchstat@latest

### Run
    go test -bench=. -benchmem -benchtime=2s -count 5 ./... | tee bench.txt
    benchstat -sort=name bench.txt

Results on an Apple Silicon M1 Max

    name                         time/op
    RequestIDBothDisabled-10      508ns ± 1%
    RequestIDWithBothOutPuts-10  1.31µs ± 1%
    RequestIDWithGinCtx-10       1.13µs ± 2%
    RequestIDWithReqHeader-10    1.13µs ± 1%
    
    name                         alloc/op
    RequestIDBothDisabled-10       783B ± 0%
    RequestIDWithBothOutPuts-10  1.82kB ± 0%
    RequestIDWithGinCtx-10       1.45kB ± 0%
    RequestIDWithReqHeader-10    1.49kB ± 0%
    
    name                         allocs/op
    RequestIDBothDisabled-10       10.0 ± 0%
    RequestIDWithBothOutPuts-10    19.0 ± 0%
    RequestIDWithGinCtx-10         16.0 ± 0%
    RequestIDWithReqHeader-10      16.0 ± 0%
