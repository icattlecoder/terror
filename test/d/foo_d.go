package d

import (
	"fmt"
	"github.com/icattlecoder/terrors"
	"io"
)

func FuncD() error {
	fmt.Println("func D called")
	return terrors.Trace(io.EOF)
}