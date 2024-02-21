package test

import (
	"fmt"
	"github.com/icattlecoder/terrors/pkgtest/a"
	"testing"
)

func TestPrint(t *testing.T) {

	err:=a.FuncA("b")
	fmt.Printf("%+v",err)
}