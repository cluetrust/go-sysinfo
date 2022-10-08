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

//go:build solaris

package solaris

import (
        "context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/prometheus/procfs"

	"github.com/elastic/go-sysinfo/internal/registry"
	"github.com/elastic/go-sysinfo/providers/shared"
	"github.com/elastic/go-sysinfo/types"

	"github.com/siebenmann/go-kstat"
)

func init() {
	registry.Register(newSolarisSystem(""))
}

type solarisSystem struct {
	procFS procFS
}

func newSolarisSystem(hostFS string) solarisSystem {
	mountPoint := filepath.Join(hostFS, procfs.DefaultMountPoint)
	fs, _ := procfs.NewFS(mountPoint)
	return solarisSystem{
		procFS: procFS{FS: fs, mountPoint: mountPoint},
	}
}

func (s solarisSystem) Host() (types.Host, error) {
	return newHost()
}

type host struct {
	info types.HostInfo
}

func (h *host) Info() types.HostInfo {
	return h.info
}

func (h *host) Memory() (*types.HostMemoryInfo, error) {
	var mem types.HostMemoryInfo

	tok, err := kstat.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to get total physical memory: %w", err)
	}
	defer tok.Close()

	phys, err := tok.GetNamed("unix", 0, "system_pages", "physmem")
	if err != nil {
		return nil, fmt.Errorf("failed to get total physical memory: %w", err)
	}
	mem.Total = phys.UintVal

	return &mem, nil
}

// NetworkCounters reports data from /proc/net on linux
func (h *host) NetworkCounters() (*types.NetworkCountersInfo, error) {
	/*
		snmpRaw, err := ioutil.ReadFile(h.procFS.path("net/snmp"))
		if err != nil {
			return nil, err
		}
		snmp, err := getNetSnmpStats(snmpRaw)
		if err != nil {
			return nil, err
		}

		netstatRaw, err := ioutil.ReadFile(h.procFS.path("net/netstat"))
		if err != nil {
			return nil, err
		}
		netstat, err := getNetstatStats(netstatRaw)
		if err != nil {
			return nil, err
		}

		return &types.NetworkCountersInfo{SNMP: snmp, Netstat: netstat}, nil
	*/
	return &types.NetworkCountersInfo{}, nil
}

func (h *host) CPUTime() (types.CPUTimes, error) {
	/*
		stat, err := h.procFS.NewStat()
		if err != nil {
			return types.CPUTimes{}, err
		}

		return types.CPUTimes{
			User:    time.Duration(stat.CPUTotal.User * float64(time.Second)),
			System:  time.Duration(stat.CPUTotal.System * float64(time.Second)),
			Idle:    time.Duration(stat.CPUTotal.Idle * float64(time.Second)),
			IOWait:  time.Duration(stat.CPUTotal.Iowait * float64(time.Second)),
			IRQ:     time.Duration(stat.CPUTotal.IRQ * float64(time.Second)),
			Nice:    time.Duration(stat.CPUTotal.Nice * float64(time.Second)),
			SoftIRQ: time.Duration(stat.CPUTotal.SoftIRQ * float64(time.Second)),
			Steal:   time.Duration(stat.CPUTotal.Steal * float64(time.Second)),
		}, nil
	*/
	return types.CPUTimes{}, nil
}

func (h *host) FQDNWithContext(ctx context.Context) (string, error) {
	return shared.FQDNWithContext(ctx)
}

func (h *host) FQDN() (string, error) {
	return h.FQDNWithContext(context.Background())
}

func newHost() (*host, error) {
	h := &host{} // TODO: fix
	r := &reader{}
	r.architecture(h)
	r.bootTime(h)
	r.hostname(h)
	r.network(h)
	r.kernelVersion(h)
	r.os(h)
	r.time(h)
	r.uniqueID(h)
	return h, r.Err()
}

type reader struct {
	errs []error
}

func (r *reader) addErr(err error) bool {
	if err != nil {
		if !errors.Is(err, types.ErrNotImplemented) {
			r.errs = append(r.errs, err)
		}
		return true
	}
	return false
}

func (r *reader) Err() error {
	if len(r.errs) > 0 {
		return errors.Join(r.errs...)
	}
	return nil
}

func (r *reader) architecture(h *host) {
	v, err := Architecture()
	if r.addErr(err) {
		return
	}
	h.info.Architecture = v
}

func (r *reader) bootTime(h *host) {
	v, err := bootTime()
	if r.addErr(err) {
		return
	}
	h.info.BootTime = v
}

func (r *reader) hostname(h *host) {
	v, err := os.Hostname()
	if r.addErr(err) {
		return
	}
	h.info.Hostname = v
}

func (r *reader) network(h *host) {
	ips, macs, err := shared.Network()
	if r.addErr(err) {
		return
	}
	h.info.IPs = ips
	h.info.MACs = macs
}

func (r *reader) kernelVersion(h *host) {
	v, err := KernelVersion()
	if r.addErr(err) {
		return
	}
	h.info.KernelVersion = v
}

func (r *reader) os(h *host) {
	v, err := OperatingSystem()
	if r.addErr(err) {
		return
	}
	h.info.OS = v
}

func (r *reader) time(h *host) {
	h.info.Timezone, h.info.TimezoneOffsetSec = time.Now().Zone()
}

func (r *reader) uniqueID(h *host) {
	v, err := MachineID()
	if r.addErr(err) {
		return
	}
	h.info.UniqueID = v
}

type procFS struct {
	procfs.FS
	mountPoint string
}

func (fs *procFS) path(p ...string) string {
	elem := append([]string{fs.mountPoint}, p...)
	return filepath.Join(elem...)
}
