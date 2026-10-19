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

	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/stretchr/testify/assert"

	"github.com/containerd/containerd/v2/pkg/oci"
)

func TestBuildContainerSpecMetricsCPU(t *testing.T) {
	period := uint64(100000)
	shares := uint64(1024)
	quota := int64(50000)
	spec := &oci.Spec{
		Linux: &specs.Linux{
			Resources: &specs.LinuxResources{
				CPU: &specs.LinuxCPU{
					Period: &period,
					Shares: &shares,
					Quota:  &quota,
				},
			},
		},
	}

	metrics := buildContainerSpecMetrics(spec, nil, 1)

	got := map[string]uint64{}
	for _, m := range metrics {
		got[m.Name] = m.Value.Value
	}

	assert.Equal(t, period, got[containerSpecCPUPeriod.Name])
	assert.Equal(t, shares, got[containerSpecCPUShares.Name])
	assert.Equal(t, uint64(quota), got[containerSpecCPUQuota.Name])
}
