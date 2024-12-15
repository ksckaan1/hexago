package customerrors

import (
	"errors"
	"fmt"
)

var (
	ErrRunnerNotImplemented = errors.New("runner not implemented")
	ErrDirMustBeFolder      = errors.New("dir must be folder")
	ErrDomainNotFound       = errors.New("domain not found")
	ErrInvalidInstanceName  = errors.New("invalid instance name")
	ErrInvalidPkgName       = errors.New("invalid pkg name")
	ErrInvalidCmdName       = errors.New("invalid cmd name")
	ErrTemplateCanNotParsed = errors.New("template can not parsed")
	ErrAlreadyExist         = errors.New("already exist")
	ErrSuppressed           = errors.New("")
)

// Custom errors

type ErrInitGoModule struct {
	Message string
}

func (e ErrInitGoModule) Error() string {
	return e.Message
}

type ErrTemplateCanNotExecute struct {
	Message string
}

func (e ErrTemplateCanNotExecute) Error() string {
	return e.Message
}

type ErrFormatGoFile struct {
	Message string
}

func (e ErrFormatGoFile) Error() string {
	return e.Message
}

type ErrInvalidPortName struct {
	PortName string
}

func (e ErrInvalidPortName) Error() string {
	return fmt.Sprintf("invalid port name: %s", e.PortName)
}
