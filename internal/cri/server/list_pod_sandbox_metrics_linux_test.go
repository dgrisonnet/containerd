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
	runtime "k8s.io/cri-api/pkg/apis/runtime/v1"
)

func TestExtractHugetlbMetrics(t *testing.T) {
	baseLabels := []string{"pod", "ns", "ctr", "id", "img"}

	for _, tc := range []struct {
		name  string
		stats cgroupMetrics
		want  []*runtime.Metric
	}{
		{
			name: "cgroup v1",
			stats: cgroupMetrics{v1: &cg1.Metrics{
				Hugetlb: []*cg1.HugetlbStat{
					{Usage: 100, Max: 200, Pagesize: "2MB"},
				},
			}},
			want: []*runtime.Metric{
				{
					Name:        containerHugetlbUsageBytes.Name,
					MetricType:  runtime.MetricType_GAUGE,
					LabelValues: []string{"pod", "ns", "ctr", "id", "img", "2MB"},
					Value:       &runtime.UInt64Value{Value: 100},
				},
				{
					Name:        containerHugetlbMaxUsageBytes.Name,
					MetricType:  runtime.MetricType_GAUGE,
					LabelValues: []string{"pod", "ns", "ctr", "id", "img", "2MB"},
					Value:       &runtime.UInt64Value{Value: 200},
				},
			},
		},
		{
			name: "cgroup v2",
			stats: cgroupMetrics{v2: &cg2.Metrics{
				Hugetlb: []*cg2.HugeTlbStat{
					{Current: 300, Max: 400, Pagesize: "1GB"},
				},
			}},
			want: []*runtime.Metric{
				{
					Name:        containerHugetlbUsageBytes.Name,
					MetricType:  runtime.MetricType_GAUGE,
					LabelValues: []string{"pod", "ns", "ctr", "id", "img", "1GB"},
					Value:       &runtime.UInt64Value{Value: 300},
				},
				{
					Name:        containerHugetlbMaxUsageBytes.Name,
					MetricType:  runtime.MetricType_GAUGE,
					LabelValues: []string{"pod", "ns", "ctr", "id", "img", "1GB"},
					Value:       &runtime.UInt64Value{Value: 400},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := &criService{}
			got, err := c.extractHugetlbMetrics(tc.stats, baseLabels, 0)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
