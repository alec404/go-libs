package restyx

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/semconv/v1.20.0/httpconv"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"resty.dev/v3"
)

// TraceClient 在 Resty 客户端上安装 OTel 中间件。
func TraceClient(c *resty.Client, cfg *config) {
	tracer := cfg.TracerProvider.Tracer(
		cfg.TracerName,
		trace.WithInstrumentationVersion(SemVersion()),
	)

	c.AddRequestMiddleware(spanRequestMiddleware(tracer, cfg))
	c.AddResponseMiddleware(spanResponseMiddleware(cfg))
	c.OnError(onError(cfg))
}

func spanRequestMiddleware(tracer trace.Tracer, cfg *config) resty.RequestMiddleware {
	return func(c *resty.Client, r *resty.Request) error {
		if cfg.Skipper(r) {
			return nil
		}

		ctx, _ := tracer.Start(r.Context(), r.Method, cfg.SpanStartOptions...)

		cfg.Propagators.Inject(ctx, propagation.HeaderCarrier(r.Header))
		r.SetContext(ctx)
		return nil
	}
}

func spanResponseMiddleware(cfg *config) resty.ResponseMiddleware {
	return func(_ *resty.Client, res *resty.Response) error {
		if cfg.Skipper(res.Request) {
			return nil
		}
		span := trace.SpanFromContext(res.Request.Context())
		span.SetAttributes(httpconv.ClientResponse(res.RawResponse)...)

		// Setting request attributes here since res.Request.RawRequest is nil
		span.SetName(cfg.SpanNameFormatter("", res.Request))
		span = setRequestAttributes(span, cfg, res.Request)

		status := res.StatusCode()
		span.SetAttributes(attribute.Int("http.status_code", status))
		if status >= http.StatusInternalServerError {
			span.SetStatus(codes.Error, http.StatusText(status))
		}

		span.End()
		return nil
	}
}

func onError(cfg *config) resty.ErrorHook {
	return func(req *resty.Request, err error) {
		if !cfg.Skipper(req) {
			span := trace.SpanFromContext(req.Context())
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.SetName(cfg.SpanNameFormatter("", req))
			span = setRequestAttributes(span, cfg, req)
			span.End()
		}

	}
}

func setRequestAttributes(span trace.Span, _ *config, req *resty.Request) trace.Span {
	span.SetAttributes(httpconv.ClientRequest(req.RawRequest)...)
	span.SetAttributes(attribute.String("http.path", req.RawRequest.URL.Path))

	return span
}
