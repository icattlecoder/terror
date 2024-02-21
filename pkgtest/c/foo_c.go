package c

import (
	"fmt"
	"github.com/icattlecoder/terrors/pkgtest/d"
	"github.com/pkg/errors"
)

func FuncC() error {
	return fmt.Errorf("abc%w",errors.WithStack(d.FuncD()))
}
