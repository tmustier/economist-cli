package errors

import (
	"errors"
	"fmt"
)

type UserError struct {
	Msg string
}

func (e UserError) Error() string {
	return e.Msg
}

func NewUserError(format string, args ...any) error {
	return UserError{Msg: fmt.Sprintf(format, args...)}
}

func IsUserError(err error) bool {
	var ue UserError
	return errors.As(err, &ue)
}

type PaywallError struct{}

func (PaywallError) Error() string {
	return "paywall detected"
}

func IsPaywallError(err error) bool {
	var pe PaywallError
	return errors.As(err, &pe)
}
