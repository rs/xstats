package xstats_test

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/rs/xhandler"
	"github.com/rs/xstats"
	"github.com/rs/xstats/dogstatsd"
	"golang.org/x/net/context"
)

func ExampleNewHandler() {
	c := xhandler.Chain{}

	// Install the metric handler with dogstatsd backend client and some env tags
	flushInterval := 5 * time.Second
	tags := []string{"role:my-service"}
	statsdWriter, err := net.Dial("udp", "127.0.0.1:8126")
	if err != nil {
		log.Fatal(err)
	}
	c.UseC(xstats.NewHandler(dogstatsd.New(statsdWriter, flushInterval), tags))

	// Here is your handler
	h := c.Handler(xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		// Get the xstats request's instance from the context. You can safely assume it will
		// be always there, if the handler is removed, xstats.FromContext will return a nop
		// instance.
		m := xstats.FromContext(ctx)

		// Count something
		m.Count("requests", 1, "route:index")
	}))

	http.Handle("/", h)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
