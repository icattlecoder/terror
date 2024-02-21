package terrors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type caller struct {
	file  string
	line  int
	frame runtime.Frame
}

// funcname removes the path prefix component of a function's name reported by func.Name().
func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}

func (c *caller) String() string {
	//i := strings.LastIndex(c.frame.Function, ".")
	//_, f := filepath.Split(c.file)
	return fmt.Sprintf("%s\n\t%s:%d", c.frame.Function, c.file, c.line)
}

type chainError struct {
	caller *stack
	err    error
}

func (c *chainError) Format(f fmt.State, verb rune) {
	fmt.Println("-------------------+++++++++++++++-")

	c.caller.Format(f,verb)

	//switch verb {
	//case 'v':
	//	if f.Flag('+') {
	//		io.WriteString(f, printIdent(c, "\n", ""))
	//		return
	//	}
	//	fallthrough
	//case 's':
	//	io.WriteString(f, c.Error())
	//case 'q':
	//	fmt.Fprintf(f, "%q", c.Error())
	//}
}

//func (c *chainError) getTraceError() string {
//
//	var cerr *chainError
//	if errors.As(c.err, &cerr) {
//		return c.caller.String() + "\n" + cerr.getTraceError()
//	}
//	if c.err != nil {
//		return c.caller.String() + "\n" + c.err.Error()
//	}
//	return c.caller.String()
//}

func (c *chainError) Unwrap() error {
	return c.err
}

func (c *chainError) Error() string {
	return c.err.Error()
}

func Trace(err error) error {
	if err == nil {
		return nil
	}

	var cerr *chainError
	if errors.As(err, &cerr){
		fmt.Println(err,"--------")
		return err
	}
	fmt.Println(err,"++++++++")
	cs:=callers()
	pChain := &chainError{caller: cs, err: err}
	return pChain
}

func Traced(err error) bool  {
	var cerr *chainError
	ok:= errors.As(err, &cerr)
	return ok
}

func Unwrap(err error)  error {
	var cerr *chainError
	if errors.As(err, &cerr){
		return cerr
	}
	return err
}