package a

import (
	"github.com/icattlecoder/terrors"
	"github.com/icattlecoder/terrors/test/b"
	"github.com/icattlecoder/terrors/test/d"
)

func FuncA(c string) error {
	if c == "b" {
		return terrors.Trace(b.FuncC())
	}
	return terrors.Trace(d.FuncD())
}
