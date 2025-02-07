package try_test

import (
	"errors"
	"fmt"
	"log"
	"testing"

	"try"

	"github.com/cheekybits/is"
)

func TestTryExample(t *testing.T) {
	try.MaxRetries = 20
	SomeFunction := func() (string, error) {
		return "", nil
	}

	err := try.Do(func(attempt int) (bool, error) {
		var err error
		_, err = SomeFunction()
		return attempt < 5, err // try 5 times
	})
	if err != nil {
		log.Fatalln("error:", err)
	}
}

func TestTryExamplePanic(t *testing.T) {
	SomeFunction := func() (string, error) {
		panic("something went badly wrong")
	}

	err := try.Do(func(attempt int) (retry bool, err error) {
		retry = attempt < 5 // try 5 times
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(fmt.Sprintf("panic: %v", r))
			}
		}()
		_, err = SomeFunction()
		return
	})
	if err != nil {
		//log.Fatalln("error:", err)
	}
}

func TestTryDoSuccessful(t *testing.T) {
	is := is.New(t)
	callCount := 0
	err := try.Do(func(attempt int) (bool, error) {
		callCount++
		return attempt < 5, nil
	})
	is.NoErr(err)
	is.Equal(callCount, 1)
}

func TestTryDoFailed(t *testing.T) {
	is := is.New(t)
	theErr := errors.New("something went wrong")
	callCount := 0
	err := try.Do(func(attempt int) (bool, error) {
		callCount++
		return attempt < 5, theErr
	})
	is.Equal(err, theErr)
	is.Equal(callCount, 5)
}

func TestTryPanics(t *testing.T) {
	is := is.New(t)
	theErr := errors.New("something went wrong")
	callCount := 0
	err := try.Do(func(attempt int) (retry bool, err error) {
		retry = attempt < 5
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(fmt.Sprintf("panic: %v", r))
			}
		}()
		callCount++
		if attempt > 2 {
			panic("I don't like three")
		}
		err = theErr
		return
	})
	is.Equal(err.Error(), "panic: I don't like three")
	is.Equal(callCount, 5)
}

func TestRetryLimit(t *testing.T) {
	is := is.New(t)
	err := try.Do(func(attempt int) (bool, error) {
		return true, errors.New("nope")
	})
	is.OK(err)
	is.Equal(try.IsMaxRetries(err), true)
}

func TestTryStructExample(t *testing.T) {
	tr := try.NewTryPtr()
	tr.MaxRetries = 20

	SomeFunction := func() (string, error) {
		return "", nil
	}

	err := tr.Do(func(attempt int) (bool, error) {
		var err error
		_, err = SomeFunction()
		return attempt < 5, err // try 5 times
	})
	if err != nil {
		log.Fatalln("error:", err)
	}
}

func TestTryStructExamplePanic(t *testing.T) {
	SomeFunction := func() (string, error) {
		panic("something went badly wrong")
	}

	tr := try.NewTryPtr()
	tr.MaxRetries = 20

	err := tr.Do(func(attempt int) (retry bool, err error) {
		retry = attempt < 5 // try 5 times
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(fmt.Sprintf("panic: %v", r))
			}
		}()
		_, err = SomeFunction()
		return
	})
	if err != nil {
		//log.Fatalln("error:", err)
	}
}

func TestTryStructDoSuccessful(t *testing.T) {
	is := is.New(t)
	callCount := 0

	tr := try.NewTryPtr()
	tr.MaxRetries = 20

	err := tr.Do(func(attempt int) (bool, error) {
		callCount++
		return attempt < 5, nil
	})
	is.NoErr(err)
	is.Equal(callCount, 1)
}

func TestTryStructDoFailed(t *testing.T) {
	is := is.New(t)
	theErr := errors.New("something went wrong")

	tr := try.NewTryPtr()
	tr.MaxRetries = 20

	callCount := 0
	err := tr.Do(func(attempt int) (bool, error) {
		callCount++
		return attempt < 5, theErr
	})
	is.Equal(err, theErr)
	is.Equal(callCount, 5)
}

func TestTryStructPanics(t *testing.T) {
	is := is.New(t)
	theErr := errors.New("something went wrong")

	tr := try.NewTryPtr()
	tr.MaxRetries = 20

	callCount := 0
	err := tr.Do(func(attempt int) (retry bool, err error) {
		retry = attempt < 5
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(fmt.Sprintf("panic: %v", r))
			}
		}()
		callCount++
		if attempt > 2 {
			panic("I don't like three")
		}
		err = theErr
		return
	})
	is.Equal(err.Error(), "panic: I don't like three")
	is.Equal(callCount, 5)
}

func TestRetryStructLimit(t *testing.T) {
	is := is.New(t)

	tr := try.NewTryPtr()
	tr.MaxRetries = 20

	err := tr.Do(func(attempt int) (bool, error) {
		return true, errors.New("nope")
	})
	is.OK(err)
	is.Equal(tr.IsMaxRetries(err), true)
}
