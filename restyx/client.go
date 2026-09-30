package restyx

import (
	"github.com/alec404/go-libs/restyx/net"
	"go.opentelemetry.io/otel"
	"resty.dev/v3"
)

// NewClient 创建注入了 OTel 与 Logger 的 Resty 客户端。
func NewClient(opts ...Option) *resty.Client {
	cfg := newConfig(opts...)

	if cfg.TracerProvider == nil {
		cfg.TracerProvider = otel.GetTracerProvider()
	}

	c := resty.New()

	if cfg.Logger != nil {
		c.SetLogger(NewAdapter(cfg.Logger))
		c.SetDebug(true)
		c.SetDebugLogFormatter(resty.DebugLogJSONFormatter)
	}

	// 不使用 CookieJar
	c.SetCookieJar(nil)

	// 安装 OTel
	TraceClient(c, cfg)

	c.AddRetryConditions(func(res *resty.Response, err error) bool {
		if net.IsConnectionReset(err) || net.IsHTTP2ConnectionForceClosed(err) || net.IsBrokenPipe(err) || net.IsProbableEOF(err) {
			return true
		}
		return false
	})

	return c
}
