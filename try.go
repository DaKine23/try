package try

import "errors"

// MaxRetries is the maximum number of retries before bailing.
var MaxRetries = 10

var errMaxRetriesReached = errors.New("exceeded retry limit")

var globalTry *Try = NewTryPtr()

// Func represents functions that can be retried.
type Func func(attempt int) (retry bool, err error)

// Do keeps trying the function until the second argument
// returns false, or no error is returned.
func Do(fn Func) error {
	if globalTry.MaxRetries != MaxRetries {
		globalTry.MaxRetries = MaxRetries
	}
	return globalTry.Do(fn)
}

// IsMaxRetries checks whether the error is due to hitting the
// maximum number of retries or not.
func IsMaxRetries(err error) bool {
	return err == errMaxRetriesReached
}

//Try capsules try in a struct
type Try struct {
	MaxRetries int
}

//NewTryPtr creates a new Try struct and returns its pointer with default MaxRetries of 10
func NewTryPtr() *Try {
	return &Try{10}
}

// Do keeps trying the function until the second argument
// returns false, or no error is returned.
func (t *Try) Do(fn Func) error {
	var err error
	var cont bool
	attempt := 1
	for {
		cont, err = fn(attempt)
		if !cont || err == nil {
			break
		}
		attempt++
		if attempt > t.MaxRetries {
			return errMaxRetriesReached
		}
	}
	return err
}

// IsMaxRetries checks whether the error is due to hitting the
// maximum number of retries or not.
func (t *Try) IsMaxRetries(err error) bool {
	return IsMaxRetries(err)
}
