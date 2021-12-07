package try

import (
	"math"
	"time"
)

//TimeCalcTry Trywith a NextTimeFunc to calculate the next waiting time in case of error
type TimeCalcTry struct {
	MaxRetries int
	NextTimeFunc
}

//TimesTry Try with an array of waiting times per try
type TimesTry struct {
	Times []time.Duration
}

//NextTimeFunc is a function type to calculate the next wait time by index and last waiting time
type NextTimeFunc = func(last time.Duration, index int) (time.Duration, error)

//FixedDifference returns a NextTimeFunc that adds diff ms to get the next time to wait
func FixedDifference(diff time.Duration) NextTimeFunc {

	return func(last time.Duration, index int) (time.Duration, error) {
		return diff + last, nil
	}

}

//ExponentialBackoff follows wait_time = base * multiplier^n
func ExponentialBackoff(startValue time.Duration, multiplier int) NextTimeFunc {

	return func(last time.Duration, index int) (time.Duration, error) {
		return startValue * time.Duration((math.Pow(float64(multiplier), float64(index)))), nil
	}

}

//NewDurationSlice create a time.Duration slice once to use it with TimesTry
func NewDurationSlice(ntf NextTimeFunc, startValue time.Duration, length int) ([]time.Duration, error) {

	last := startValue
	res := make([]time.Duration, length)
	var err error = nil
	res[0] = startValue
	for i := 1; i < length; i++ {
		last, err = ntf(last, i)
		if err != nil {
			return nil, err
		}
		res[i] = last
	}
	return res, nil
}

//NewTimesTryPtr creates a new Try struct and returns its pointer
func NewTimesTryPtr(times []time.Duration) *TimesTry {

	return &TimesTry{times}

}

// Do keeps trying the function until the second argument
// returns false, or no error is returned.
func (t *TimesTry) Do(fn Func) error {
	var err error
	var cont bool
	attempt := 1
	for {
		cont, err = fn(attempt)
		if !cont || err == nil {
			break
		}
		time.Sleep(t.Times[attempt-1]) // wait
		attempt++
		if attempt > len(t.Times) {
			return errMaxRetriesReached
		}
	}
	return err
}

// IsMaxRetries checks whether the error is due to hitting the
// maximum number of retries or not.
func (t *TimesTry) IsMaxRetries(err error) bool {
	return IsMaxRetries(err)
}
