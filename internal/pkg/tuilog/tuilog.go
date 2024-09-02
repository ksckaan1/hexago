package tuilog

import (
	"github.com/samber/do"
)

type TUILog struct{}

func New(i *do.Injector) (*TUILog, error) {
	return &TUILog{}, nil
}

func (t *TUILog) Info(msg string, title ...string) {
	t.log(infoType, msg, title...)
}

func (t *TUILog) Success(msg string, title ...string) {
	t.log(successType, msg, title...)
}

func (t *TUILog) Warning(msg string, title ...string) {
	t.log(warningType, msg, title...)
}

func (t *TUILog) Error(msg string, title ...string) {
	t.log(errorType, msg, title...)
}