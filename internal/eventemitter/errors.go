package eventemitter

import "errors"

var ErrInvalidNumberOfArguments = errors.New("events: invalid number of arguments")
var ErrInvalidArgumentType = errors.New("events: invalid argument type")
