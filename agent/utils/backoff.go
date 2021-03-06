// Copyright 2014-2015 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package utils

import (
	"math"
	mathrand "math/rand"
	"time"
)

type Backoff interface {
	Reset()
	Duration() time.Duration
}

type SimpleBackoff struct {
	current        time.Duration
	start          time.Duration
	max            time.Duration
	jitterMultiple float64
	multiple       float64
}

// NewSimpleBackoff creates a Backoff which ranges from min to max increasing by
// multiple each time.
// It also adds (and yes, the jitter is always added, never
// subtracted) a random amount of jitter up to jitterMultiple percent (that is,
// jitterMultiple = 0.0 is no jitter, 0.15 is 15% added jitter). The total time
// may exceed "max" when accounting for jitter, such that the absolute max is
// max + max * jiterMultiple
func NewSimpleBackoff(min, max time.Duration, jitterMultiple, multiple float64) *SimpleBackoff {
	return &SimpleBackoff{
		start:          min,
		current:        min,
		max:            max,
		jitterMultiple: jitterMultiple,
		multiple:       multiple,
	}
}

func (sb *SimpleBackoff) Duration() time.Duration {
	ret := sb.current
	sb.current = time.Duration(math.Min(float64(sb.max.Nanoseconds()), float64(float64(sb.current.Nanoseconds())*sb.multiple)))

	return AddJitter(ret, time.Duration(int64(float64(ret)*sb.jitterMultiple)))

	return ret
}

func (sb *SimpleBackoff) Reset() {
	sb.current = sb.start
}

// AddJitter adds an amount of jitter between 0 and the given jitter to the
// given duration
func AddJitter(duration time.Duration, jitter time.Duration) time.Duration {
	var randJitter int64
	if jitter.Nanoseconds() == 0 {
		randJitter = 0
	} else {
		randJitter = mathrand.Int63n(jitter.Nanoseconds())
	}
	return time.Duration(duration.Nanoseconds() + randJitter)
}
