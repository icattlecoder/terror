//go:xnotrace
package test

import (
	"errors"
	"io"

	"github.com/icattlecoder/terrors"
)

//func FuncA() {
//
//}
//
func FuncErr() error {
	return terrors.Trace(errors.New("func err"))
}

func IFuncErr() (interface{}, error) {
	return nil, errors.New("func err")
}

func INamedFuncErr() (i interface{}, err error) {
	i, err = nil, io.EOF
	err = terrors.Trace(err)
	return
}

func FuncWithError() (interface{}, error) {
	var i int
	if i > 1 {
		return nil, terrors.Trace(errors.New("i>1"))
	}

	if i > 2 {
		if err := FuncErr(); err != nil {
			return nil, terrors.Trace(err)
		}
	}

	_ = map[string]interface{}{
		"f": func() {

		},
		"fe": func() error {
			return terrors.Trace(errors.New(""))
		},
	}

	switch i {
	case 0:
		return nil, terrors.Trace(FuncErr())
	case 1:
		if _, err := IFuncErr(); err != nil {
			return nil, terrors.Trace(err)
		}
		result0, result1 := IFuncErr()
		return result0, terrors.Trace(result1)
	}
	f := func() error {
		return terrors.Trace(io.EOF)
	}
	if err := f(); err != nil {
		return nil, terrors.Trace(err)
	}
	_, err := IFuncErr()
	if err != nil {
		result0, result1 := IFuncErr()
		return result0, terrors.Trace(result1)
	}
	result0, result1 := IFuncErr()
	return result0, terrors.Trace(result1)
}
