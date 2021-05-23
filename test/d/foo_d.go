package d

import (
	"github.com/icattlecoder/terrors"
	"io"
)

func FuncD() error {
	return terrors.Trace(io.EOF)
}