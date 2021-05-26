package b

import (
	"fmt"
	"github.com/icattlecoder/terrors"
	"github.com/icattlecoder/terrors/test/d"
)

func FuncC() error {
	fmt.Println("funcC called")
	err := fmt.Errorf("%w", d.FuncD())
	fmt.Println("xxxxxxxxxxxxxxxxxxx funcC", terrors.Traced(err))
	return err

}
