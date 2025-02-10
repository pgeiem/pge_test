// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//Copied from https://github.com/rickb777/date/blob/v2/timespan/timespan_test.go

package parser

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

const zero time.Duration = 0

const minusOneNano time.Duration = -1

var t0327 = time.Date(2015, 3, 27, 0, 0, 0, 0, time.UTC)
var t0328 = time.Date(2015, 3, 28, 0, 0, 0, 0, time.UTC)
var t0329 = time.Date(2015, 3, 29, 0, 0, 0, 0, time.UTC) // n.b. clocks go forward (UK)
var t0330 = time.Date(2015, 3, 30, 0, 0, 0, 0, time.UTC)

func isEq(t *testing.T, i int, a, b interface{}, msg ...interface{}) {
	t.Helper()
	if a != b {
		sa := make([]string, len(msg))
		for i, m := range msg {
			sa[i] = fmt.Sprintf(", %v", m)
		}
		t.Errorf("%d: %+v is not equal to %+v%s", i, a, b, strings.Join(sa, ""))
	}
}

func TestZeroTimeSpan(t *testing.T) {
	ts := ZeroTimeSpan(t0327)
	isEq(t, 0, ts.Mark(), t0327)
	isEq(t, 0, ts.Duration(), zero)
	isEq(t, 0, ts.End(), t0327)
}

func TestNewTimeSpan(t *testing.T) {
	ts1 := BetweenTimes(t0327, t0327)
	isEq(t, 0, ts1.Mark(), t0327)
	isEq(t, 0, ts1.Duration(), zero)
	isEq(t, 0, ts1.IsEmpty(), true)
	isEq(t, 0, ts1.End(), t0327)

	ts2 := BetweenTimes(t0327, t0328)
	isEq(t, 0, ts2.Mark(), t0327)
	isEq(t, 0, ts2.Duration(), time.Hour*24)
	isEq(t, 0, ts2.IsEmpty(), false)
	isEq(t, 0, ts2.End(), t0328)

	ts3 := BetweenTimes(t0329, t0327)
	isEq(t, 0, ts3.Mark(), t0327)
	isEq(t, 0, ts3.Duration(), time.Hour*48)
	isEq(t, 0, ts3.IsEmpty(), false)
	isEq(t, 0, ts3.End(), t0329)
}

func TestTSEnd(t *testing.T) {
	ts1 := TimeSpan{t0328, time.Hour * 24}
	isEq(t, 0, ts1.Start(), t0328)
	isEq(t, 0, ts1.End(), t0329)

	// not normalised, deliberately
	ts2 := TimeSpan{t0328, -time.Hour * 24}
	isEq(t, 0, ts2.Start(), t0327)
	isEq(t, 0, ts2.End(), t0328)
}

func TestTSShiftBy(t *testing.T) {
	ts1 := BetweenTimes(t0327, t0328).ShiftBy(time.Hour * 24)
	isEq(t, 0, ts1.Mark(), t0328)
	isEq(t, 0, ts1.Duration(), time.Hour*24)
	isEq(t, 0, ts1.End(), t0329)

	ts2 := BetweenTimes(t0328, t0329).ShiftBy(-time.Hour * 24)
	isEq(t, 0, ts2.Mark(), t0327)
	isEq(t, 0, ts2.Duration(), time.Hour*24)
	isEq(t, 0, ts2.End(), t0328)
}

func TestTSExtendBy(t *testing.T) {
	ts1 := BetweenTimes(t0327, t0328).ExtendBy(time.Hour * 24)
	isEq(t, 0, ts1.Mark(), t0327)
	isEq(t, 0, ts1.Duration(), time.Hour*48)
	isEq(t, 0, ts1.End(), t0329)

	ts2 := BetweenTimes(t0328, t0329).ExtendBy(-time.Hour * 48)
	isEq(t, 0, ts2.Mark(), t0327)
	isEq(t, 0, ts2.Duration(), time.Hour*24)
	isEq(t, 0, ts2.End(), t0328)
}

func TestTSExtendWithoutWrapping(t *testing.T) {
	ts1 := BetweenTimes(t0327, t0328).ExtendWithoutWrapping(time.Hour * 24)
	isEq(t, 0, ts1.Mark(), t0327)
	isEq(t, 0, ts1.Duration(), time.Hour*48)
	isEq(t, 0, ts1.End(), t0329)

	ts2 := BetweenTimes(t0328, t0329).ExtendWithoutWrapping(-time.Hour * 48)
	isEq(t, 0, ts2.Mark(), t0328)
	isEq(t, 0, ts2.Duration(), zero)
	isEq(t, 0, ts2.End(), t0328)
}

func TestTSString(t *testing.T) {
	s := BetweenTimes(t0327, t0328).String()
	isEq(t, 0, s, "24h0m0s from 2015-03-27 00:00:00 to 2015-03-28 00:00:00")
}

func TestTSEqual(t *testing.T) {
	// use Berlin, which is UTC+1/+2
	berlin, _ := time.LoadLocation("Europe/Berlin")
	t0 := time.Date(2015, 2, 20, 10, 13, 25, 0, time.UTC)
	t1 := t0.Add(time.Hour)
	z0 := ZeroTimeSpan(t0)
	ts1 := z0.ExtendBy(time.Hour)

	cases := []struct {
		a, b TimeSpan
	}{
		{a: z0, b: BetweenTimes(t0, t0)},
		{a: z0, b: z0.In(berlin)},
		{a: ts1, b: ts1},
		{a: ts1, b: BetweenTimes(t0, t1)},
		{a: ts1, b: ts1.In(berlin)},
		{a: ts1, b: ZeroTimeSpan(t1).ExtendBy(-time.Hour)},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.a), func(t *testing.T) {
			if !c.a.Equal(c.b) {
				t.Errorf("%d: %v is not equal to %v", i, c.a, c.b)
			}
		})
	}
}

func TestTSNotEqual(t *testing.T) {
	t0 := time.Date(2015, 2, 20, 10, 13, 25, 0, time.UTC)
	t1 := t0.Add(time.Hour)

	cases := []struct {
		a, b TimeSpan
	}{
		{a: ZeroTimeSpan(t0), b: TimeSpanOf(t0, time.Hour)},
		{a: ZeroTimeSpan(t0), b: ZeroTimeSpan(t1)},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.a), func(t *testing.T) {
			if c.a.Equal(c.b) {
				t.Errorf("%d: %v is not equal to %v", i, c.a, c.b)
			}
		})
	}
}

func TestTSContains(t *testing.T) {
	ts := BetweenTimes(t0327, t0329)
	isEq(t, 0, ts.Contains(t0327.Add(minusOneNano)), false)
	isEq(t, 0, ts.Contains(t0327), true)
	isEq(t, 0, ts.Contains(t0328), true)
	isEq(t, 0, ts.Contains(t0329.Add(minusOneNano)), true)
	isEq(t, 0, ts.Contains(t0329), false)
}

func TestTSIn(t *testing.T) {
	ts := ZeroTimeSpan(t0327).In(time.FixedZone("Test", 7200))
	isEq(t, 0, ts.Mark().Equal(t0327), true)
	isEq(t, 0, ts.Duration(), zero)
	isEq(t, 0, ts.End().Equal(t0327), true)
}

func TestTSMerge1(t *testing.T) {
	ts1 := BetweenTimes(t0327, t0328)
	ts2 := BetweenTimes(t0327, t0330)
	m1 := ts1.Merge(ts2)
	m2 := ts2.Merge(ts1)
	isEq(t, 0, m1.Mark(), t0327)
	isEq(t, 0, m1.End(), t0330)
	isEq(t, 0, m1, m2)
}

func TestTSMerge2(t *testing.T) {
	ts1 := BetweenTimes(t0328, t0329)
	ts2 := BetweenTimes(t0327, t0330)
	m1 := ts1.Merge(ts2)
	m2 := ts2.Merge(ts1)
	isEq(t, 0, m1.Mark(), t0327)
	isEq(t, 0, m1.End(), t0330)
	isEq(t, 0, m1, m2)
}

func TestTSMerge3(t *testing.T) {
	ts1 := BetweenTimes(t0329, t0330)
	ts2 := BetweenTimes(t0327, t0330)
	m1 := ts1.Merge(ts2)
	m2 := ts2.Merge(ts1)
	isEq(t, 0, m1.Mark(), t0327)
	isEq(t, 0, m1.End(), t0330)
	isEq(t, 0, m1, m2)
}

func TestTSMergeOverlapping(t *testing.T) {
	ts1 := BetweenTimes(t0327, t0329)
	ts2 := BetweenTimes(t0328, t0330)
	m1 := ts1.Merge(ts2)
	m2 := ts2.Merge(ts1)
	isEq(t, 0, m1.Mark(), t0327)
	isEq(t, 0, m1.End(), t0330)
	isEq(t, 0, m1, m2)
}

func TestTSMergeNonOverlapping(t *testing.T) {
	ts1 := BetweenTimes(t0327, t0328)
	ts2 := BetweenTimes(t0329, t0330)
	m1 := ts1.Merge(ts2)
	m2 := ts2.Merge(ts1)
	isEq(t, 0, m1.Mark(), t0327)
	isEq(t, 0, m1.End(), t0330)
	isEq(t, 0, m1, m2)
}
