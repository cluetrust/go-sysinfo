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
	"fmt"
	"strconv"
	"strings"

	"github.com/elastic/go-sysinfo/types"
	"golang.org/x/sys/unix"
)

func OperatingSystem() (*types.OSInfo, error) {
	return getOSInfo()
}

func arrayToString(x [257]byte) string {
	data := make([]byte, 0, len(x))
	for _, v := range x {
		if v == 0 {
			break
		}
		data = append(data, byte(v))
	}
	return string(data)
}

func getOSInfo() (*types.OSInfo, error) {
	var uname unix.Utsname
	if err := unix.Uname(&uname); err != nil {
		return nil, fmt.Errorf("os_info: %w", err)
	}

   // fmt.Println(string(uname.Sysname[:])) => SunOS
   // fmt.Println(string(uname.Nodename[:])) => hostname
   // fmt.Println(string(uname.Machine[:])) => architecture (base)
   release_string := arrayToString(uname.Release)
   release := strings.Split(release_string, ".")
   major, err := strconv.ParseInt(release[0], 10, 64)
   if err != nil {
		return nil, fmt.Errorf("os_info:major: %w", err)
   }
   minor, err := strconv.ParseInt(release[1], 10, 64)
   if err != nil {
		return nil, fmt.Errorf("os_info:minor: %w", err)
   }
   build := arrayToString(uname.Version)

	return &types.OSInfo{
		Type:     "solaris",
		Family:   "illumos",
		Platform: "smartos",
		Name:     release_string,
		Version:  release_string,
		Major:    int(major),
		Minor:    int(minor),
		Patch:    0,
		Build:    build,
	}, nil
}
