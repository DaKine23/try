package try_test

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"
	"try"

	"github.com/cheekybits/is"
)

func TestTimesTryExample(t *testing.T) {
	tr := try.NewTimesTryPtr([]time.Duration{1, 1, 1, 1, 1})

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

func TestTimesTryExamplePanic(t *testing.T) {
	SomeFunction := func() (string, error) {
		panic("something went badly wrong")
	}

	tr := try.NewTimesTryPtr([]time.Duration{1, 1, 1, 1, 1})

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

func TestTimesTryDoSuccessful(t *testing.T) {
	is := is.New(t)
	callCount := 0

	tr := try.NewTimesTryPtr([]time.Duration{1, 1, 1, 1, 1})

	err := tr.Do(func(attempt int) (bool, error) {
		callCount++
		return attempt < 5, nil
	})
	is.NoErr(err)
	is.Equal(callCount, 1)
}

func TestTimesTryDoFailed(t *testing.T) {
	is := is.New(t)
	theErr := errors.New("something went wrong")

	tr := try.NewTimesTryPtr([]time.Duration{1, 1, 1, 1, 1})

	callCount := 0
	err := tr.Do(func(attempt int) (bool, error) {
		callCount++
		return attempt < 5, theErr
	})
	is.Equal(err, theErr)
	is.Equal(callCount, 5)
}

func TestTimesTryPanics(t *testing.T) {
	is := is.New(t)
	theErr := errors.New("something went wrong")

	tr := try.NewTimesTryPtr([]time.Duration{1, 1, 1, 1, 1})

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

func TestReTimesTryLimit(t *testing.T) {
	is := is.New(t)

	tr := try.NewTimesTryPtr([]time.Duration{1, 1, 1, 1, 1})

	err := tr.Do(func(attempt int) (bool, error) {
		return true, errors.New("nope")
	})
	is.OK(err)
	is.Equal(tr.IsMaxRetries(err), true)
}

func TestNewDurationSlice(t *testing.T) {
	type args struct {
		ntf        try.NextTimeFunc
		startValue time.Duration
		length     int
	}
	tests := []struct {
		name    string
		args    args
		want    []time.Duration
		wantErr bool
	}{
		{
			name: "Fixed Difference 1 to 5",
			args: args{
				ntf:        try.FixedDifference(1000 * time.Millisecond),
				startValue: 1000 * time.Millisecond,
				length:     5,
			},
			want:    []time.Duration{1000 * time.Millisecond, 2000 * time.Millisecond, 3000 * time.Millisecond, 4000 * time.Millisecond, 5000 * time.Millisecond},
			wantErr: false,
		}, {
			name: "Fixed Difference 1 to 1",
			args: args{
				ntf:        try.FixedDifference(0),
				startValue: 1000 * time.Millisecond,
				length:     5,
			},
			want:    []time.Duration{1000 * time.Millisecond, 1000 * time.Millisecond, 1000 * time.Millisecond, 1000 * time.Millisecond, 1000 * time.Millisecond},
			wantErr: false,
		}, {
			name: "Exponential Backoff",
			args: args{
				ntf:        try.ExponentialBackoff(time.Second, 2),
				startValue: 1000 * time.Millisecond,
				length:     5,
			},
			want:    []time.Duration{1000 * time.Millisecond, 2000 * time.Millisecond, 4000 * time.Millisecond, 8000 * time.Millisecond, 16000 * time.Millisecond},
			wantErr: false,
		}, {
			name: "error when calculating",
			args: args{
				ntf:        func(last time.Duration, index int) (time.Duration, error) { return 0, errors.New("ERROR!") },
				startValue: 1000 * time.Millisecond,
				length:     5,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := try.NewDurationSlice(tt.args.ntf, tt.args.startValue, tt.args.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDurationSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDurationSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
