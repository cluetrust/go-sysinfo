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
	"sync"
	"time"

	"github.com/siebenmann/go-kstat"
)

var (
	bootTimeValue time.Time  // Cached boot time.
	bootTimeLock  sync.Mutex // Lock that guards access to bootTime.
)

func bootTime() (time.Time, error) {
	bootTimeLock.Lock()
	defer bootTimeLock.Unlock()

	if !bootTimeValue.IsZero() {
		return bootTimeValue, nil
	}

	tok, err := kstat.Open()
	if err != nil {
		return time.Time{}, err
	}
	defer tok.Close()


	boot,err := tok.GetNamed("unix",0, "system_misc","boot_time")
	if err != nil {
		return time.Time{}, err
	}
	bootTimeValue = time.Unix(int64(boot.UintVal), 0)
	return bootTimeValue, nil
}
