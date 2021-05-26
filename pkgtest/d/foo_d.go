package d

import (
	"github.com/pkg/errors"
	"io"
)

func FuncD() error {
	return errors.New(io.EOF.Error())
}