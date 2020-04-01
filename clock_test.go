package clock

import (
	"testing"
	"time"
)

func TestNewClock(t *testing.T) {
	c := NewClock()

	first := c.Now()
	time.Sleep(100 * time.Millisecond)
	second := c.Now()

	if first != second {
		t.Error("results from subsequent calls to Now() don't match")
		t.Errorf("first:  %s", first)
		t.Errorf("second: %s", second)
	}
}

func TestNewClockAt(t *testing.T) {
	now := time.Date(2020, 4, 1, 12, 12, 12, 0, time.UTC)

	var expected time.Time
	var actual time.Time
	T := func(testName string) {
		if !actual.Equal(expected) {
			t.Errorf("test failed: %s", testName)
			t.Errorf("expected: %s", expected)
			t.Errorf("actual:   %s", actual)
		}
	}

	c := NewClockAt(now)

	expected = now
	actual = c.Now()
	T("initial call to Now()")

	d := 10 * time.Second
	c.Sleep(d)
	expected = now.Add(d)
	actual = c.Now()
	T("check Now() after Sleep()")
}

func TestNewClockFromSlice(t *testing.T) {
	time1 := time.Date(2020, 4, 1, 12, 12, 12, 0, time.UTC)
	time2 := time.Date(2021, 4, 1, 12, 12, 12, 0, time.UTC)

	var expected time.Time
	var actual time.Time
	T := func(testName string) {
		if !actual.Equal(expected) {
			t.Errorf("test failed: %s", testName)
			t.Errorf("expected: %s", expected)
			t.Errorf("actual:   %s", actual)
		}
	}

	c := NewClockFromSlice(time1, time2)

	expected = time1
	actual = c.Now()
	T("first instant returned first")

	expected = time2
	actual = c.Now()
	T("second instant returned second")

	// We ran out of instants, so on the next call we should panic
	defer func() {
		if err := recover(); err != nil {
			t.Log("recovered from panic as expected")
		}
	}()

	// panics
	c.Now()

	t.Error("missing panic")
}

func TestNewClockFromCallbacks(t *testing.T) {
	fixedTime := time.Date(2020, 4, 1, 12, 12, 12, 0, time.UTC)
	nowCallCt := 0
	nowLastCalledWith := -1
	nowFn := func(i int) time.Time {
		nowCallCt++
		nowLastCalledWith = i
		return fixedTime
	}
	sleepCallCt := 0
	sleepLastCalledWith := -1 * time.Second
	sleepFn := func(d time.Duration) {
		sleepCallCt++
		sleepLastCalledWith = d
	}

	c := NewClockFromCallbacks(nowFn, sleepFn)

	actual := c.Now()
	if nowLastCalledWith != 0 {
		t.Error("1: failed to call nowFn")
	}
	if sleepCallCt != 0 {
		t.Error("1: unexpectedly called sleepFn")
	}
	if !actual.Equal(fixedTime) {
		t.Error("1: call to Now()")
	}

	actual = c.Now()
	if nowLastCalledWith != 1 { // counter increments correctly
		t.Error("2: failed to call nowFn")
	}
	if sleepCallCt != 0 {
		t.Error("2: unexpectedly called sleepFn")
	}
	if !actual.Equal(fixedTime) { // our function returns the same result
		t.Error("2: call to Now()")
	}

	c.Sleep(1 * time.Minute)
	actual = c.Now()
	if nowLastCalledWith != 2 { // counter increments correctly
		t.Error("3: failed to call nowFn")
	}
	if sleepCallCt != 1 || sleepLastCalledWith != 1*time.Minute {
		t.Error("3: unexpected result calling sleepFn")
	}
	if !actual.Equal(fixedTime) { // our function returns the same result
		t.Error("3: call to Now()")
	}
}

func TestClock_Now(t *testing.T) {
	tolerance := 50 * time.Millisecond

	// When c is nil we should actually return time.Now()
	var c *Clock
	now := time.Now()
	actual := c.Now()

	diff := actual.Sub(now)
	if diff < -tolerance || diff > tolerance {
		t.Error("didn't return real time.Now() value as expected")
		t.Errorf("expected: %s ± %d ms", now, tolerance/time.Millisecond)
		t.Errorf("actual:   %s", actual)
	}
}

func TestClock_Sleep(t *testing.T) {
	d := 500 * time.Millisecond
	tolerance := d / 10

	// When c is nil we should actually call time.Sleep()
	var c *Clock
	before := time.Now()
	c.Sleep(d)
	after := time.Now()

	diff := after.Sub(before) - d
	if diff < -tolerance || diff > tolerance {
		t.Error("didn't call time.Sleep() as expected")
		t.Errorf("expected sleep: %d ± %d ms", d/time.Millisecond, tolerance/time.Millisecond)
		t.Errorf("actual sleep:   %d ms", after.Sub(before)/time.Millisecond)
	}

}

func TestClock_TimeSince(t *testing.T) {
	d := 25 * time.Minute
	now := time.Date(2020, 4, 1, 12, 12, 12, 0, time.UTC)
	before := now.Add(-d)
	c := NewClockAt(now)

	actual := c.TimeSince(before)
	if actual != d {
		t.Errorf("expected: %s", d)
		t.Errorf("actual:   %s", actual)
	}
}
