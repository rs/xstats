# HTTP Handler Metrics

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/rs/xmetrics) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/rs/xmetrics/master/LICENSE) [![Build Status](https://travis-ci.org/rs/xmetrics.svg?branch=master)](https://travis-ci.org/rs/xmetrics)

xmetric is a generic client for service instrumentation.

## Supported Clients

- [StatsD](https://github.com/b/statsd_spec)
- [DogStatsD](http://docs.datadoghq.com/guides/dogstatsd/#datagram-format)

## Install

    go get github.com/rs/xmetric

## Usage

```go
// Defines interval between flushes to statsd server
flushInterval := 5 * time.Second

// Global tags sent with all metrics (only with supported clients like datadog's)
tags := []string{"role:my-service"}

// Connection to the statsd server
statsdWriter, err := net.Dial("udp", "127.0.0.1:8126")
if err != nil {
    log.Fatal(err)
}

// Create the metric client
m := xmetric.New(statsd.New(statsdWriter, flushInterval), tags)

// Send some observations
m.Count("requests", 1, "tag")
m.Timing("something", 5*time.Millisecond, "tag")
```

Integration with [github.com/rs/xhandler](https://github.com/rs/xhandler):

```go
var xh xhandler.CtxHandler

// Here is your handler
xh = xhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    // Get the xmetric request client from the context. You can safely assume it will
    // be always there, if the handler is removed, xmetric.FromContext will return a nopClient
    m := xmetric.FromContext(ctx)

    // Count something
    m.Count("requests", 1, "route:index")
})

// Install the metric handler with statsd backend client and some env tags
flushInterval := 5 * time.Second
tags := []string{"role:my-service"}
statsdWriter, err := net.Dial("udp", "127.0.0.1:8126")
if err != nil {
    log.Fatal(err)
}
xh := xmetric.NewHandler(statsd.New(statsdWriter, flushInterval), tags, xh)

// Root context
var h http.Handler
ctx := context.Background()
h = xhandler.CtxHandler(ctx, lh)
http.Handle("/", h)

if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Fatal(err)
}
```

## Licenses

All source code is licensed under the [MIT License](https://raw.github.com/rs/xmetrics/master/LICENSE).
