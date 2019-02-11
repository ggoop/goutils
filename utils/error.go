package utils

import (
	"errors"
	"fmt"
)

// Err : error,string
// Code : xxx xxx xxxx,area-product-app 10bit
type GError struct {
	Err  error
	Code int
}

func (e GError) Error() string {
	return fmt.Sprintf("%s [%d]", e.Err.Error(), e.Code)
}

// err : error,string
// code : xxx xxx xxxx,area-product-app 10bit
func NewError(err interface{}, code int) GError {
	var e error
	if v, ok := err.(GError); ok {
		e = v
	} else if v, ok := err.(error); ok {
		e = v
	} else if v, ok := err.(string); ok {
		e = errors.New(v)
	} else {
		e = errors.New("未知异常")
	}
	return GError{Err: e, Code: code}
}
