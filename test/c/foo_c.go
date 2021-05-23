package c

import (
	"github.com/icattlecoder/terrors"
	"github.com/icattlecoder/terrors/test/d"
)

func FuncC() error {
	return terrors.Trace(d.FuncD())
}
