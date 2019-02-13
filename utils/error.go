package utils

import (
	"errors"
	"fmt"
)

// err : error ,
// code :x xx xxx xxxx,level-product-app-number 10bit
type GError struct {
	Err  error
	Code int
}

func (e GError) Error() string {
	return fmt.Sprintf("%s [%d]", e.Err.Error(), e.Code)
}

// err : error ,
// code :x xx xxx xxxx,level-product-app-number 10bit
func NewError(err interface{}, code int) GError {
	var e error
	if v, ok := err.(GError); ok {
		e = v
	} else if v, ok := err.(error); ok {
		e = v
	} else if v, ok := err.(string); ok {
		e = errors.New(v)
	} else {
		e = errors.New(fmt.Sprint(err))
	}
	return GError{Err: e, Code: code}
}
