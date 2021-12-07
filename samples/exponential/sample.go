package main

import (
	"errors"
	"fmt"
	"time"
	"try"
)

var initTime time.Time = time.Now()

func main() {

	// specify first waiting delay
	firstRetryTryAfter := time.Millisecond * 20
	// specify a durationSlice (here we use ExponentialBackoff)
	times, _ := try.NewDurationSlice(
		try.ExponentialBackoff(firstRetryTryAfter, 2), firstRetryTryAfter, 7)

	// initialize a TimesTry struct with the slice
	trier := try.NewTimesTryPtr(times)

	//try 3 times

	//generic part of the trier
	go trier.Do(func(attempt int) (bool, error) {

		//your call goes here
		err := toBeTried("first call")
		// try it 3 times or max tries defined by retrier
		// in case of err == nil (success) it quits the loop
		return attempt < 3, err

	})

	//try max times

	//generic part of the trier
	go trier.Do(func(attempt int) (bool, error) {

		//your call goes here
		err := toBeTried("second call")
		// true means "use max tries defined by retrier"
		// in case of err == nil (success) it quits the loop
		return true, err

	})

	//succeed

	//generic part of the trier
	go trier.Do(func(attempt int) (bool, error) {

		//your call goes here
		err := toBeTriedAndSucceed("third call")
		// true means "use max tries defined by retrier"
		// in case of err == nil (success) it quits the loop
		return true, err

	})

	//wait until keypress
	var input string
	fmt.Scanln(&input)

}

func toBeTried(msg string) error {
	fmt.Println(msg, "i was tried after ", time.Since(initTime))
	return errors.New("something went wrong")
}

func toBeTriedAndSucceed(msg string) error {
	fmt.Println(msg, "i was tried after ", time.Since(initTime))
	return nil
}
