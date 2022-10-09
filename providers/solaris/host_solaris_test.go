// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package solaris

import (
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
    "github.com/elastic/go-sysinfo/internal/registry"
    
)

var _ registry.HostProvider = solarisSystem{}

func TestHost(t *testing.T) {
    solaris := newSolarisSystem("")
	assert.NotEmpty(t, solaris)
    
	host, err := solaris.Host()
    if err != nil {
        t.Fatal(err)
    }
    info := host.Info()
	assert.Equal(t, "i86pc", info.Architecture)
	assert.Less(t, time.Unix(971128197, 0), info.BootTime)
	assert.Empty(t, info.Containerized)
	assert.NotEmpty(t, info.Hostname)
	assert.NotEmpty(t, info.IPs)
	assert.Equal(t, "5.11", info.KernelVersion)
	assert.NotEmpty(t, info.MACs)
	assert.NotEqual(t, "", info.Timezone)
	assert.NotEqual(t, "", info.UniqueID)
	assert.NotEmpty(t, info.OS)
	assert.Equal(t, "solaris", info.OS.Type)
}
