package util

import "errors"

// Catch calls the given function and recovers and returns any panic that is thrown with an error. The function
// returns nil if f doesn't panic.
func Catch(f func()) (err error) {
	defer func() {
		switch e := recover().(type) {
		case nil:
		case error:
			err = e
		case string:
			err = errors.New(e)
		default:
			panic(e)
		}
	}()
	f()
	return
}
