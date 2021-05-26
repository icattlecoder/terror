package parser

import (
	"bytes"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
	"go/token"
	"strings"
)

const terrorsPkgPath = "github.com/icattlecoder/terrors"

type NoTraced struct {
	Positions []token.Position
}

func (e *NoTraced) Error() string {
	buf := bytes.Buffer{}
	for _, l := range e.Positions {
		buf.WriteString(l.String())
		buf.WriteString("\n")
	}
	return buf.String()
}

type PackageParser struct {
	*decorator.Package
	NoTraced
}

func New(p *decorator.Package) *PackageParser {
	return &PackageParser{Package: p}
}

func (p *PackageParser) Pos(node dst.Node) {
	pos := p.Package.Fset.Position(p.Decorator.Ast.Nodes[node].Pos())
	p.NoTraced.Positions = append(p.NoTraced.Positions, pos)
}

func (p *PackageParser) Run(w bool) error {
	for _, syntax := range p.Syntax {
		p.ParseFile(syntax)
	}
	if len(p.NoTraced.Positions) == 0 {
		return nil
	}
	if w {
		return p.Save()
	}
	return &p.NoTraced
}

func (p *PackageParser) ParseFile(file *dst.File) {

	if notrace(file) {
		return
	}

	for _, decl := range file.Decls {
		fdecl, ok := decl.(*dst.FuncDecl)
		if !ok {
			continue
		}
		p.parseFunc(fdecl)
	}
	return
}

func notrace(node dst.Node) bool {
	_, _, ps := dstutil.Decorations(node)
	for _, p := range ps {
		for _, d := range p.Decs {
			if strings.Contains(d, "go:notrace") {
				return true
			}
		}
	}
	return false
}

func (p *PackageParser) parseFunc(node dst.Node) {

	if notrace(node) {
		return
	}

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
		return
	}

	if results == nil || len(results.List) == 0 {
		return
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
		return
	}

	parser := blockListFuncParser{
		funcParser: funcParser{
			errsIndex:     errsIndex,
			result:        results.List,
			PackageParser: p,
		},
	}
	parser.inspect(body)
	return
}

type funcParser struct {
	*PackageParser
	errsIndex []int
	result    []*dst.Field
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
		b.Pos(returnStmt)
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
		b.Pos(returnStmt)
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
		return false
	case *dst.CaseClause:
		nb := blockListFuncParser{
			funcParser: b.funcParser,
			list:       &nodeType.Body,
			root:       node,
		}
		dst.Inspect(nodeType, nb.inspect)
		return false
	case *dst.FuncLit:
		b.parseFunc(nodeType)
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
