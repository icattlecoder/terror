package test

import (
	"errors"
	"io"
)

//func FuncA() {
//
//}
//
func FuncErr() error {
	return errors.New("func err")
}

func IFuncErr() (interface{}, error) {
	return nil, errors.New("func err")
}

func INamedFuncErr() (i interface{}, err error) {
	i, err = nil, io.EOF
	return
}

func FuncWithError() (interface{}, error) {
	var i int
	if i > 1 {
		return nil, errors.New("i>1")
	}

	if i > 2 {
		if err := FuncErr(); err != nil {
			return nil, err
		}
	}

	_ = map[string]interface{}{
		"f": func() {

		},
		"fe": func() error {
			return errors.New("")
		},
	}

	switch i {
	case 0:
		return nil, FuncErr()
	case 1:
		if _, err := IFuncErr(); err != nil {
			return nil, err
		}
		return IFuncErr()
	}
	f := func() error {
		return io.EOF
	}
	if err := f(); err != nil {
		return nil, err
	}


	_, err := IFuncErr()
	if err != nil {
		return IFuncErr()
	}
	return IFuncErr()
}
