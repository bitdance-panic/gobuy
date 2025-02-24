package mtl

import (
	"log"
	"net"
	"net/http"

	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Registry *prometheus.Registry // 全局 Prometheus 注册表

func InitMetric(serviceName, metricsPort, registryAddr string) {
	// 初始化 prometheus 注册表
	Registry = prometheus.NewRegistry()
	Registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	// 注册到 Consul
	r, _ := consul.NewConsulRegister(registryAddr)
	addr, _ := net.ResolveTCPAddr("tcp", metricsPort)
	registryInfo := &registry.Info{
		ServiceName: serviceName,
		Addr:        addr,
		Weight:      1,
		Tags:        map[string]string{"metrics": "prometheus"},
	}

	_ = r.Register(registryInfo)
	server.RegisterShutdownHook(func() {
		r.Deregister(registryInfo)
	})

	// 暴露指标端点
	http.Handle("/metrics", promhttp.HandlerFor(Registry, promhttp.HandlerOpts{}))
	go func() {
		log.Printf("[Metrics] Serving metrics at %s/metrics", metricsPort)
		if err := http.ListenAndServe(metricsPort, nil); err != nil {
			panic(err)
		}
	}()
}
