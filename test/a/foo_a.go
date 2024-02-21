package a

import (
	"fmt"
	"github.com/icattlecoder/terrors"
	"github.com/icattlecoder/terrors/test/b"
)

type ErrA struct {
	FromFunc string
	err      error
}

func (e *ErrA) Error() string {
	return e.FromFunc
}

func NewErrA(name string) *ErrA {
	return &ErrA{FromFunc: name}
}

func (e *ErrA) Wrap(err error) *ErrA {
	e.err = err
	return e
}

func (e *ErrA) Unwrap() error {
	return e.err
}

func FuncA(cs string) error {
	fmt.Println("func A called")
	err:= b.FuncC()
	fmt.Println("xxxxxxxxxxxxxxxxxxx FuncA",terrors.Traced(err))

	return terrors.Trace(NewErrA("funcA").Wrap(fmt.Errorf("%w", b.FuncC())))
}
