package test

import (
	"fmt"
	"github.com/icattlecoder/terrors"
	"github.com/icattlecoder/terrors/test/a"
	"testing"
)

func TestPrint(t *testing.T) {

	//err := terrors.Trace(a.FuncA("b"))
	//fmt.Println(err.Error())
	//fmt.Println(terrors.Print(err))
	err := terrors.Trace(a.FuncA("c"))
	//fmt.Println(err.Error())
	fmt.Println(terrors.Traced(err),"___________________")
	fmt.Printf("%+v", terrors.Unwrap(err))
}