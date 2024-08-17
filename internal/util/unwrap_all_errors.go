package util

import "errors"

func UnwrapAllErrors(err error) error {
	if err == nil {
		return nil
	}
	lastErr := err
	for {
		e := errors.Unwrap(lastErr)
		if e == nil {
			return lastErr
		}
		lastErr = e
	}
}
