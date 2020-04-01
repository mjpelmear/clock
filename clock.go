// This is based on https://raw.githubusercontent.com/smartystreets/clock/master/clock.go
// but with some different implementation choices

package clock

import "time"

type NowCallbackFn func(int) time.Time
type SleepCallbackFn func(time.Duration)
type Clock struct {
	index         int
	nowCallback   NowCallbackFn
	sleepCallback SleepCallbackFn
	napDurations  []time.Duration
}

func NewClock() *Clock {
	return NewClockAt(time.Now())
}

func NewClockAt(instant time.Time) *Clock {
	return NewClockFromCallbacks(NewDefaultSleeperCallbacks(instant))
}

func NewClockFromSlice(instants ...time.Time) *Clock {
	return NewClockFromCallbacks(slicesNowFn(instants...), nil)
}

func slicesNowFn(instants ...time.Time) NowCallbackFn {
	return func(i int) time.Time {
		if i >= len(instants) {
			panic("ran out of moments frozen in time")
		} else {
			return instants[i]
		}
	}
}

func NewClockFromCallbacks(nowFn NowCallbackFn, sleepFn SleepCallbackFn) *Clock {
	return &Clock{
		nowCallback:   nowFn,
		sleepCallback: sleepFn,
		napDurations:  []time.Duration{},
	}
}

func NewDefaultSleeperCallbacks(start time.Time) (NowCallbackFn, SleepCallbackFn) {
	sleeperNow := start

	nowFn := func(i int) time.Time {
		return sleeperNow
	}

	sleepFn := func(d time.Duration) {
		sleeperNow = sleeperNow.Add(d)
	}

	return nowFn, sleepFn
}

func (c *Clock) Now() time.Time {
	if c == nil {
		return time.Now()
	}

	defer func() {
		c.index++
	}()

	return c.nowCallback(c.index)
}

func (c *Clock) Sleep(duration time.Duration) {
	if c == nil {
		time.Sleep(duration)
	} else if c.sleepCallback != nil {
		c.sleepCallback(duration)
		c.napDurations = append(c.napDurations, duration)
	}
}

func (c *Clock) TimeSince(before time.Time) time.Duration {
	return c.Now().Sub(before)
}

func (c *Clock) NapDurations() []time.Duration {
	return c.napDurations
}

func (c *Clock) ClearNapDurations() {
	c.napDurations = []time.Duration{}
}
