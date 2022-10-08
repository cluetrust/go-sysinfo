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
//go:build (amd64 && cgo)
// +build amd64,cgo

package solaris

// #include <zone.h>
import "C"

import (
	"fmt"
//	"unsafe"
)

// MachineID returns the Zone UUID also accessible via
func MachineID() (string, error) {
	return getHostUUID()
}

func getHostUUID() (string, error) {
	// zoneid_t is an id_t, is an int type
	
//	var zoneidC C.zoneid_t
	var uuid [C.ZONENAME_MAX]C.char

	zoneidC := C.getzoneid()
	ret := C.getzonenamebyid(zoneidC, &uuid[0], C.ZONENAME_MAX)

	if ret != 0 {
		return "", fmt.Errorf("gethostuuid failed with %v", ret)
	}


	return C.GoString(&uuid[0]), nil
}
