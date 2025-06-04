package core

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/prometheus/common/expfmt"
)

type NodeExporterMetrics struct {
	LoadPercent int   `json:"load_percent"`
	DiskSize    int64 `json:"disk_size"`
	DiskUsed    int64 `json:"disk_used"`
	MemorySize  int64 `json:"memory_size"`
	MemoryUsed  int64 `json:"memory_used"`
}

func ParseNodeExporterMetrics(data []byte) (*NodeExporterMetrics, error) {

	metrics := &NodeExporterMetrics{}

	var parser expfmt.TextParser
	mfs, err := parser.TextToMetricFamilies(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("metric parsing failed: %v", err)
	}

	mf := mfs["node_cpu_seconds_total"]
	if mf == nil {
		return nil, fmt.Errorf("metric 'node_cpu_seconds_total' not found")
	}

	core_count := 0
	for _, m := range mf.Metric {
		for _, v := range m.GetLabel() {
			if v.GetName() == "cpu" {
				num, _ := strconv.ParseInt(v.GetValue(), 10, 64)
				if num > int64(core_count) {
					core_count = int(num)
				}
			}
		}
	}
	// offset by 1
	core_count++

	load := 0.0
	load15 := mfs["node_load15"]
	for _, m := range load15.Metric {
		load = m.GetGauge().GetValue()
	}
	metrics.LoadPercent = int(load / float64(core_count) * 100)

	// disk size
	mf = mfs["node_filesystem_size_bytes"]
	if mf == nil {
		return nil, fmt.Errorf("metric 'node_filesystem_size_bytes' not found")
	}
	for _, m := range mf.Metric {
		for _, v := range m.GetLabel() {
			if v.GetName() == "mountpoint" && v.GetValue() == "/" {
				metrics.DiskSize = int64(m.Gauge.GetValue())
				break
			}
		}
	}

	// disk usage
	mf = mfs["node_filesystem_free_bytes"]
	if mf == nil {
		return nil, fmt.Errorf("metric 'node_filesystem_free_bytes' not found")
	}
	for _, m := range mf.Metric {
		for _, v := range m.GetLabel() {
			if v.GetName() == "mountpoint" && v.GetValue() == "/" {
				metrics.DiskUsed = metrics.DiskSize - int64(m.Gauge.GetValue())
				break
			}
		}
	}

	// memory size
	mf = mfs["node_memory_MemTotal_bytes"]
	if mf == nil {
		return nil, fmt.Errorf("metric 'node_memory_MemTotal_bytes' not found")
	}
	for _, m := range mf.Metric {
		metrics.MemorySize = int64(m.Gauge.GetValue())
		break
	}

	// memory usage
	mf = mfs["node_memory_MemFree_bytes"]
	if mf == nil {
		return nil, fmt.Errorf("metric 'node_memory_MemFree_bytes' not found")
	}
	for _, m := range mf.Metric {
		metrics.MemoryUsed = metrics.MemorySize - int64(m.Gauge.GetValue())
		break
	}
	// subtract buffer
	mf = mfs["node_memory_Buffers_bytes"]
	if mf == nil {
		return nil, fmt.Errorf("metric 'node_memory_Buffers_bytes' not found")
	}
	for _, m := range mf.Metric {
		metrics.MemoryUsed -= int64(m.Gauge.GetValue())
		break
	}
	// subtract cache
	mf = mfs["node_memory_Cached_bytes"]
	if mf == nil {
		return nil, fmt.Errorf("metric 'node_memory_Cached_bytes' not found")
	}
	for _, m := range mf.Metric {
		metrics.MemoryUsed -= int64(m.Gauge.GetValue())
		break
	}

	return metrics, nil
}

func FetchMetrics(hostname string) (*NodeExporterMetrics, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:9100/metrics", hostname), nil)
	if err != nil {
		return nil, fmt.Errorf("request build failed: %v", err)
	}

	var client http.Client
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response-code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("body reading error: %v", err)
	}

	return ParseNodeExporterMetrics(data)
}
