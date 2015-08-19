package exception

func Throw(err error) {
	panic(err)
}

func ThrowIf(err error) {
	if err != nil {
		Throw(err)
	}
}

func Try(fn func()) (err error) {
	defer Catch(func(e error) { err = e })
	fn()
	return nil
}

type ErrorHandler func(err error)

func Catch(handler ErrorHandler) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			handler(err)
		} else {
			panic(r)
		}
	}
}
