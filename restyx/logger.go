package restyx

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
)

// adapter 实现 resty.Logger 接口，将 Resty 日志转给 Kratos log.Logger。
// Resty 只关心 Errorf/Warnf/Debugf 三个方法。

type Adapter struct{ k log.Logger }

func NewAdapter(l log.Logger) *Adapter { return &Adapter{k: l} }

func (a *Adapter) Errorf(format string, v ...any) {
	_ = a.k.Log(log.LevelError, "module", "resty", log.DefaultMessageKey, fmt.Sprintf(format, v...))
}
func (a *Adapter) Warnf(format string, v ...any) {
	_ = a.k.Log(log.LevelWarn, "module", "resty", log.DefaultMessageKey, fmt.Sprintf(format, v...))
}
func (a *Adapter) Debugf(format string, v ...any) {
	_ = a.k.Log(log.LevelDebug, "module", "resty", log.DefaultMessageKey, fmt.Sprintf(format, v...))
}
