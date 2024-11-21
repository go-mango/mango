package goweb

type abortError struct {
	status int
	err    error
}

func Abort(status int, err error) {
	panic(&abortError{
		status: status,
		err:    err,
	})
}
