package a

import (
	"github.com/icattlecoder/terrors/pkgtest/b"
)

type MyErr struct {
	Msg string
	err error
}

func (e *MyErr) Error() string {
	return e.Msg
}

func (e *MyErr) Unwrap() error {
	return e.err
}

func FuncA(c string) error {
	return b.FuncB()

	//return &MyErr{
	//	Msg: "call from FuncA",
	//	err: b.FuncB(),
	//}
}
