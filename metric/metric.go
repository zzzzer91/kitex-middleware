package metric

import (
	"context"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/stats"
	monitor "github.com/zzzzer91/prometheus-monitor"
)

// Metrics
const (
	metricKitexClientThroughput = "kitex_client_throughput"
	metricKitexClientLatencyUs  = "kitex_client_latency_us"
	metricKitexServerThroughput = "kitex_server_throughput"
	metricKitexServerLatencyUs  = "kitex_server_latency_us"
)

// Labels
const (
	labelKeyCaller = "caller"
	labelKeyCallee = "callee"
	labelKeyMethod = "method"
	labelKeyStatus = "status"
	labelKeyRetry  = "retry"

	// status
	statusSucceed = "succeed"
	statusError   = "error"

	unknownLabelValue = "unknown"
)

type clientTracer struct {
	m monitor.Monitor
}

// Start record the beginning of an RPC invocation.
func (c *clientTracer) Start(ctx context.Context) context.Context {
	return ctx
}

// NewClientTracer provide tracer for client call, addr and path is the scrape_configs for prometheus server.
func NewClientTracer(m monitor.Monitor) stats.Tracer {
	m.AddMetric(&monitor.Metric{
		Type:        monitor.Counter,
		Name:        metricKitexClientThroughput,
		Description: "Total number of RPCs completed by the client, regardless of success or failure.",
		Labels:      []string{labelKeyCaller, labelKeyCallee, labelKeyMethod, labelKeyStatus, labelKeyRetry},
	})
	m.AddMetric(&monitor.Metric{
		Type:        monitor.Histogram,
		Name:        metricKitexClientLatencyUs,
		Description: "Latency (microseconds) of the RPC until it is finished.",
		Buckets:     []float64{5000, 10000, 25000, 50000, 100000, 250000, 500000, 1000000},
		Labels:      []string{labelKeyCaller, labelKeyCallee, labelKeyMethod, labelKeyStatus, labelKeyRetry},
	})
	return &clientTracer{
		m: m,
	}
}

// Finish record after receiving the response of server.
func (c *clientTracer) Finish(ctx context.Context) {
	ri := rpcinfo.GetRPCInfo(ctx)
	if ri.Stats().Level() == stats.LevelDisabled {
		return
	}
	rpcStart := ri.Stats().GetEvent(stats.RPCStart)
	rpcFinish := ri.Stats().GetEvent(stats.RPCFinish)
	cost := rpcFinish.Time().Sub(rpcStart.Time())

	labelValues := genLabelValues(ri)
	c.m.GetMetric(metricKitexClientThroughput).Inc(labelValues)
	c.m.GetMetric(metricKitexClientLatencyUs).Observe(labelValues, float64(cost.Microseconds()))
}

type serverTracer struct {
	m monitor.Monitor
}

// NewServerTracer provides tracer for server access, addr and path is the scrape_configs for prometheus server.
func NewServerTracer(m monitor.Monitor) stats.Tracer {
	m.AddMetric(&monitor.Metric{
		Type:        monitor.Counter,
		Name:        metricKitexServerThroughput,
		Description: "Total number of RPCs completed by the server, regardless of success or failure.",
		Labels:      []string{labelKeyCaller, labelKeyCallee, labelKeyMethod, labelKeyStatus, labelKeyRetry},
	})
	m.AddMetric(&monitor.Metric{
		Type:        monitor.Histogram,
		Name:        metricKitexServerLatencyUs,
		Description: "Latency (microseconds) of RPC that had been application-level handled by the server.",
		Buckets:     []float64{5000, 10000, 25000, 50000, 100000, 250000, 500000, 1000000},
		Labels:      []string{labelKeyCaller, labelKeyCallee, labelKeyMethod, labelKeyStatus, labelKeyRetry},
	})
	return &serverTracer{
		m: m,
	}
}

// Start record the beginning of server handling request from client.
func (c *serverTracer) Start(ctx context.Context) context.Context {
	return ctx
}

// Finish record the ending of server handling request from client.
func (c *serverTracer) Finish(ctx context.Context) {
	ri := rpcinfo.GetRPCInfo(ctx)
	if ri.Stats().Level() == stats.LevelDisabled {
		return
	}

	rpcStart := ri.Stats().GetEvent(stats.RPCStart)
	rpcFinish := ri.Stats().GetEvent(stats.RPCFinish)
	cost := rpcFinish.Time().Sub(rpcStart.Time())

	labelValues := genLabelValues(ri)
	c.m.GetMetric(metricKitexServerThroughput).Inc(labelValues)
	c.m.GetMetric(metricKitexServerLatencyUs).Observe(labelValues, float64(cost.Microseconds()))
}

// genLabelValues make labels values.
func genLabelValues(ri rpcinfo.RPCInfo) []string {
	var (
		res    []string
		caller = ri.From()
		callee = ri.To()
	)

	res = append(res, defaultValIfEmpty(caller.ServiceName(), unknownLabelValue))
	res = append(res, defaultValIfEmpty(callee.ServiceName(), unknownLabelValue))
	res = append(res, defaultValIfEmpty(callee.Method(), unknownLabelValue))
	if ri.Stats().Error() != nil {
		res = append(res, statusError)
	} else {
		res = append(res, statusSucceed)
	}
	if retriedCnt, ok := callee.Tag(rpcinfo.RetryTag); ok {
		res = append(res, retriedCnt)
	} else {
		res = append(res, "0")
	}

	return res
}

func defaultValIfEmpty(val, def string) string {
	if val == "" {
		return def
	}
	return val
}
