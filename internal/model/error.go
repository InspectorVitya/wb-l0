package model

import "errors"

var (
	ErrNotExistOrder = errors.New("not exist order")
	ErrEmptyOrderID  = errors.New("empty order id")
)
