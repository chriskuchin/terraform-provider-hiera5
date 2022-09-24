package catch

// Do calls the given function and recovers an Error panic. If an Error panic is recovered, its cause is returned.
//
// A recover of something that wasn't produced with the Error() function will result in a new panic (a rethrow).
func Do(doer func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if ewc, ok := r.(*errorWithCause); ok {
				err = ewc.cause
			} else {
				panic(r)
			}
		}
	}()
	doer()
	return
}
