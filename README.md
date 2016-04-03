# XStats

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/rs/xstats) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/rs/xstats/master/LICENSE) [![Build Status](https://travis-ci.org/rs/xstats.svg?branch=master)](https://travis-ci.org/rs/xstats) [![Coverage](http://gocover.io/_badge/github.com/rs/xstats)](http://gocover.io/github.com/rs/xstats)

Package `xstats` is a generic client for service instrumentation.

`xstats` is inspired from Go-kit's [metrics](https://github.com/go-kit/kit/tree/master/metrics) package but it takes a slightly different path. Instead of having to create an instance for each metric, `xstats` use a single instance to log every metrics you want. This reduces the boiler plate when you have a lot a metrics in your app. It's also easier in term of dependency injection.

Talking about dependency injection, `xstats` comes with a [xhandler.Handler](https://github.com/rs/xhandler) integration so it can automatically inject the `xstats` client within the `net/context` of each request. Each request's `xstats` instance have its own tags storage ; This let you inject some per request contextual tags to be included with all observations sent within the lifespan of the request.

`xstats` is pluggable and comes with integration for `expvar`, `StatsD` and `DogStatsD`, the [Datadog](http://datadoghq.com) augmented version of StatsD with support for tags. More integration may come later (PR welcome).

## Supported Clients

- [StatsD](https://github.com/b/statsd_spec)
- [DogStatsD](http://docs.datadoghq.com/guides/dogstatsd/#datagram-format)
- [expvar](https://golang.org/pkg/expvar/)
- [prometheus](https://github.com/prometheus/client_golang)
- [telegraf](https://influxdata.com/blog/getting-started-with-sending-statsd-metrics-to-telegraf-influxdb)

## Install

    go get github.com/rs/xstats

## Usage

```go
// Defines interval between flushes to statsd server
flushInterval := 5 * time.Second

// Connection to the statsd server
statsdWriter, err := net.Dial("udp", "127.0.0.1:8126")
if err != nil {
    log.Fatal(err)
}

// Create the stats client
s := xstats.New(dogstatsd.New(statsdWriter, flushInterval))

// Global tags sent with all metrics (only with supported clients like datadog's)
s.AddTags("role:my-service", "dc:sv6")

// Send some observations
s.Count("requests", 1, "tag")
s.Timing("something", 5*time.Millisecond, "tag")
```

Integration with [github.com/rs/xhandler](https://github.com/rs/xhandler):

```go
var xh xhandler.HandlerC

// Here is your handler
xh = xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    // Get the xstats request's instance from the context. You can safely assume it will
    // be always there, if the handler is removed, xstats.FromContext will return a nop
    // instance.
    m := xstats.FromContext(ctx)

    // Count something
    m.Count("requests", 1, "route:index")
})

// Install the metric handler with dogstatsd backend client and some env tags
flushInterval := 5 * time.Second
tags := []string{"role:my-service"}
statsdWriter, err := net.Dial("udp", "127.0.0.1:8126")
if err != nil {
    log.Fatal(err)
}
xh = xstats.NewHandler(dogstatsd.New(statsdWriter, flushInterval), tags, xh)

// Root context
ctx := context.Background()
h := xhandler.New(ctx, xh)
http.Handle("/", h)

if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Fatal(err)
}
```

## Licenses

All source code is licensed under the [MIT License](https://raw.github.com/rs/xstats/master/LICENSE).
