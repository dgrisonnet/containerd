/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package server

import (
	"testing"

	cg1 "github.com/containerd/cgroups/v3/cgroup1/stats"
	cg2 "github.com/containerd/cgroups/v3/cgroup2/stats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	runtime "k8s.io/cri-api/pkg/apis/runtime/v1"
)

func TestExtractPressureMetrics(t *testing.T) {
	labels := []string{"pod", "ns", "ctr", "id", "image"}

	// Values in microseconds, chosen to divide evenly into whole seconds
	// so the integer-seconds conversion is unambiguous.
	stats := cgroupMetrics{
		v2: &cg2.Metrics{
			CPU: &cg2.CPUStat{
				PSI: &cg2.PSIStats{
					Full: &cg2.PSIData{Total: 1_000_000}, // 1s stalled
					Some: &cg2.PSIData{Total: 2_000_000}, // 2s waiting
				},
			},
			Memory: &cg2.MemoryStat{
				PSI: &cg2.PSIStats{
					Full: &cg2.PSIData{Total: 3_000_000}, // 3s stalled
					Some: &cg2.PSIData{Total: 4_000_000}, // 4s waiting
				},
			},
			Io: &cg2.IOStat{
				PSI: &cg2.PSIStats{
					Full: &cg2.PSIData{Total: 5_000_000}, // 5s stalled
					Some: &cg2.PSIData{Total: 6_000_000}, // 6s waiting
				},
			},
		},
	}

	c := &criService{}
	metrics, err := c.extractPressureMetrics(stats, labels, 42)
	require.NoError(t, err)

	got := make(map[string]uint64, len(metrics))
	for _, m := range metrics {
		assert.Equal(t, runtime.MetricType_COUNTER, m.MetricType, "metric %s must be a counter", m.Name)
		assert.Equal(t, int64(42), m.Timestamp)
		assert.Equal(t, labels, m.LabelValues)
		got[m.Name] = m.Value.Value
	}

	want := map[string]uint64{
		"container_pressure_cpu_stalled_seconds_total":    1,
		"container_pressure_cpu_waiting_seconds_total":    2,
		"container_pressure_memory_stalled_seconds_total": 3,
		"container_pressure_memory_waiting_seconds_total": 4,
		"container_pressure_io_stalled_seconds_total":     5,
		"container_pressure_io_waiting_seconds_total":     6,
	}
	assert.Equal(t, want, got)
}

func TestExtractPressureMetricsCgroupV1(t *testing.T) {
	// cgroup v1 has no PSI, so no pressure metrics are emitted.
	c := &criService{}
	metrics, err := c.extractPressureMetrics(cgroupMetrics{v1: &cg1.Metrics{}}, []string{"pod", "ns", "ctr", "id", "image"}, 42)
	require.NoError(t, err)
	assert.Empty(t, metrics)
}
