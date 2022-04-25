package initialize

import (
	"mall-srv/order-srv/global"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// InitJaegerTrace 生成tracer
func InitJaegerTrace() {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "106.13.214.17:6831",
		},
		ServiceName: "order-web",
	}
	// TODO nacos
	var err error
	global.Tracer, global.Closer, err = cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}

}
