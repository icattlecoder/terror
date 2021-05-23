package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/gopackages"
	"go/token"
	"golang.org/x/tools/go/packages"
	"io/ioutil"
	"os"
)

var (
	dir = flag.String("dir", "", "project dir")
	pkg = flag.String("pkg", "", "pkg path")
)

const terrorsPkgPath = "github.com/icattlecoder/terrors"

func main() {
	flag.Parse()

	pkgs, err := decorator.Load(&packages.Config{
		Dir:  *dir,
		Mode: packages.LoadAllSyntax,
	}, *pkg)

	if err != nil {
		panic(err)
	}

	if *dir == "" || *pkg == "" {
		flag.PrintDefaults()
		return
	}

	r := decorator.NewRestorerWithImports(*pkg, gopackages.New(*dir))
	for i, p := range pkgs {
		for _, syntax := range p.Syntax {
			if !parseFile(syntax) {
				continue
			}
			buf := bytes.Buffer{}
			if err := r.Fprint(&buf, syntax); err != nil {
				fmt.Println(p.GoFiles[i], "failed")
				continue
			}

			if err := safeWriteFile(buf.Bytes(), p.GoFiles[i]); err != nil {
				fmt.Println(p.GoFiles[i], "failed")
				continue
			}
			fmt.Println(p.GoFiles[i], "ok")
		}
	}
}

func safeWriteFile(data []byte, filename string) error {

	f, err := ioutil.TempFile("", "terrors")
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		return err
	}
	return os.Rename(f.Name(), filename)
}

//func refactor() {
//	f, err := decorator.Parse(code)
//	if err != nil {
//		panic(err)
//	}
//	if !parseFile(f) {
//		return
//	}
//	r := decorator.NewRestorerWithImports("github.com/icattlecoder/terrors/refacterr/test", gopackages.New("/Users/wangming/Projects/terrors/refacterr"))
//	if err := r.Print(f); err != nil {
//		panic(err)
//	}
//}

func parseFile(file *dst.File) (hasErrResult bool) {

	for _, decl := range file.Decls {
		fdecl, ok := decl.(*dst.FuncDecl)
		if !ok {
			continue
		}
		if parseFunc(fdecl) {
			hasErrResult = true
		}
	}
	return
}

func parseFunc(node dst.Node) (hasErrResult bool) {

	var errsIndex []int
	var body dst.Node
	var results *dst.FieldList
	switch f := node.(type) {
	case *dst.FuncDecl:
		results = f.Type.Results
		body = f.Body
	case *dst.FuncLit:
		results = f.Type.Results
		body = f.Body
	default:
		return false
	}

	if results == nil || len(results.List) == 0 {
		return false
	}

	for i, r := range results.List {
		rt, ok := r.Type.(*dst.Ident)
		if !ok {
			continue
		}
		if rt.Name == "error" {
			errsIndex = append(errsIndex, i)
		}
	}
	if len(errsIndex) == 0 {
		return false
	}

	parser := blockListFuncParser{
		funcParser: funcParser{
			errsIndex: errsIndex,
			result:    results.List,
		},
	}
	parser.inspect(body)
	return parser.hasErrHandled
}

type funcParser struct {
	errsIndex     []int
	result        []*dst.Field
	hasErrHandled bool
}

type blockListFuncParser struct {
	funcParser
	list *[]dst.Stmt
	root dst.Node
}

func (b *blockListFuncParser) handleReturnStmt(returnStmt *dst.ReturnStmt) bool {

	if len(returnStmt.Results) == len(b.result) {
		for _, i := range b.errsIndex {
			returnStmt.Results[i] = &dst.CallExpr{
				Fun: &dst.Ident{
					Name: "Trace",
					Path: terrorsPkgPath,
				},
				Args: []dst.Expr{returnStmt.Results[i]},
			}
		}
		b.hasErrHandled = true
		return true
	}

	if len(returnStmt.Results) == 0 {
		list := *b.list
		returnStmt := list[len(list)-1]
		list = list[0 : len(list)-1]
		for _, i := range b.errsIndex {
			assignStmt := &dst.AssignStmt{
				Lhs: cloneIdent(b.result[i].Names...),
				Tok: token.ASSIGN,
				Rhs: toTraceCall(b.result[i].Names...),
			}
			list = append(list, assignStmt)
		}
		list = append(list, returnStmt)
		*b.list = list
		b.hasErrHandled = true
		return true
	}

	if len(returnStmt.Results) == 1 {
		assignStmt := &dst.AssignStmt{
			Lhs: createIdent("result", len(b.funcParser.result)),
			Tok: token.DEFINE,
			Rhs: returnStmt.Results,
		}
		returnStmt := &dst.ReturnStmt{
			Results: createIdent("result", len(b.funcParser.result)),
		}
		for _, i := range b.errsIndex {
			returnStmt.Results[i] = &dst.CallExpr{
				Fun: &dst.Ident{
					Name: "Trace",
					Path: terrorsPkgPath,
				},
				Args: []dst.Expr{returnStmt.Results[i]},
			}
		}
		l := (*b.list)[:len(*b.list)-1]
		l = append(l, assignStmt, returnStmt)
		*b.list = l
		return true
	}

	return true
}

func (b *blockListFuncParser) inspect(node dst.Node) bool {

	if b.root == node {
		return true
	}

	switch nodeType := node.(type) {
	case *dst.ReturnStmt:
		b.handleReturnStmt(nodeType)
		return false
	case *dst.BlockStmt:
		nb := blockListFuncParser{
			funcParser: b.funcParser,
			list:       &nodeType.List,
			root:       node,
		}
		dst.Inspect(nodeType, nb.inspect)
		if nb.hasErrHandled {
			b.hasErrHandled = true
		}
		return false
	case *dst.CaseClause:
		nb := blockListFuncParser{
			funcParser: b.funcParser,
			list:       &nodeType.Body,
			root:       node,
		}
		dst.Inspect(nodeType, nb.inspect)
		if nb.hasErrHandled {
			b.hasErrHandled = true
		}
		return false
	case *dst.FuncLit:
		ok := parseFunc(nodeType)
		if ok {
			b.hasErrHandled = true
		}
		return false
	default:
		return true
	}
}

func toTraceCall(ident ...*dst.Ident) (ret []dst.Expr) {
	for _, i := range ident {
		ret = append(ret, &dst.CallExpr{
			Fun: &dst.Ident{
				Name: "Trace",
				Path: terrorsPkgPath,
			},
			Args: cloneIdent(i),
		})
	}
	return
}

func cloneIdent(ident ...*dst.Ident) (ret []dst.Expr) {
	for _, i := range ident {
		ret = append(ret, &dst.Ident{Name: i.Name})
	}
	return
}

func createIdent(namePrefix string, n int) (exprs []dst.Expr) {
	for i := 0; i < n; i++ {
		exprs = append(exprs, &dst.Ident{
			Name: fmt.Sprintf("%s%d", namePrefix, i),
		})
	}
	return
}
