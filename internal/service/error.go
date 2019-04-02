package service

import "github.com/hashicorp/go-multierror"

//~~~~~~~~~~~~~~~~~~~~~~| IError

// Error - wrap @err and append @errs to the wrapper
func Error(err error, errs ...error) error {
	return multierror.Append(err, errs...)
}

// GetErrors - unpack the error onto atomic errors
func GetErrors(err error) []error {
	if mr, ok := err.(*multierror.Error); ok {
		return mr.Errors
	}
	return []error{err}
}

// Panic - invore panic with result of an @Error() function
func Panic(err error, errs ...error) {
	panic(Error(err, errs...))
}
