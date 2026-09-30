package restyx

import (
	"net/url"

	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"resty.dev/v3"
)

// defaultSkipper provides default behaviour, which won't skip span creation.
func defaultSkipper(*resty.Request) bool {
	return false
}

func defaultSpanNameFormatter(_ string, req *resty.Request) string {
	if req.URL != "" {
		if u, err := url.Parse(req.URL); err == nil {
			return req.Method + " " + u.Path
		}
	}
	return req.Method
}

// config holds configurable parameters for the client and middleware.
// Modify此结构体即可扩展更多可选项。
type config struct {
	TracerProvider    oteltrace.TracerProvider
	Propagators       propagation.TextMapPropagator
	Logger            log.Logger
	TracerName        string
	SpanNameFormatter func(string, *resty.Request) string
	SpanStartOptions  []oteltrace.SpanStartOption
	Skipper           func(*resty.Request) bool
}

func newConfig(options ...Option) *config {
	cfg := &config{
		Propagators:       otel.GetTextMapPropagator(),
		TracerProvider:    otel.GetTracerProvider(),
		Logger:            nil,
		TracerName:        "restyx",
		SpanNameFormatter: defaultSpanNameFormatter,
		Skipper:           defaultSkipper,
	}

	defaultOpts := []Option{
		WithSpanOptions(oteltrace.WithSpanKind(oteltrace.SpanKindClient)),
	}

	options = append(defaultOpts, options...)

	for _, opt := range options {
		opt.apply(cfg)
	}

	return cfg
}

// Option applies a configuration value.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

// WithPropagators specifies propagators to use for extracting
// information from the HTTP requests. If none are specified, global
// ones will be used.
func WithPropagators(propagators propagation.TextMapPropagator) Option {
	return optionFunc(func(cfg *config) {
		if propagators != nil {
			cfg.Propagators = propagators
		}
	})
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
// If none is specified, the global provider is used.
func WithTracerProvider(provider oteltrace.TracerProvider) Option {
	return optionFunc(func(cfg *config) {
		if provider != nil {
			cfg.TracerProvider = provider
		}
	})
}
func WithLogger(l log.Logger) Option {
	return optionFunc(func(c *config) {
		if l != nil {
			c.Logger = l
		}
	})
}

// WithSkipper specifies a skipper function to determine if the middleware
// should not create a span for a determined request. If not specified,
// a span will always be created.
func WithSkipper(skipper func(r *resty.Request) bool) Option {
	return optionFunc(func(c *config) {
		if skipper != nil {
			c.Skipper = skipper
		}
	})
}

// WithSpanOptions configures an additional set of
// trace.SpanOptions, which are applied to each new span.
func WithSpanOptions(opts ...trace.SpanStartOption) Option {
	return optionFunc(func(c *config) {
		c.SpanStartOptions = append(c.SpanStartOptions, opts...)
	})
}

// WithSpanNameFormatter takes a function that will be called on every
// request and the returned string will become the Span Name.
func WithSpanNameFormatter(f func(operation string, r *resty.Request) string) Option {
	return optionFunc(func(c *config) {
		c.SpanNameFormatter = f
	})
}

// WithTracerName sets the name of the tracer used to create spans. The
func WithTracerName(name string) Option {
	return optionFunc(func(c *config) {
		c.TracerName = name
	})
}
