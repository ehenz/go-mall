package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// JaegerTrace 链路追踪的起点
// 为每个HTTP请求生成一个链路追踪的tracer
func JaegerTrace() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans:           true,
				LocalAgentHostPort: "106.13.214.17:6831",
			},
			ServiceName: "goods-web",
		}
		// TODO nacos
		tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
		if err != nil {
			panic(err)
		}
		defer closer.Close()

		// 生成起始的span
		startSpan := tracer.StartSpan(c.Request.URL.Path)
		defer startSpan.Finish()

		// 放入上下文
		c.Set("tracer", tracer)
		c.Set("start_span", startSpan)
		c.Next()
	}
}
