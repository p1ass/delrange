package delrange

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "delrange is a static analysis tool which detects delete function is called with a value different from range key."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "delrange",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.RangeStmt)(nil),
	}

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		rangeStmt, ok := n.(*ast.RangeStmt)
		if !ok || rangeStmt == nil {
			return
		}

		if _, ok := pass.TypesInfo.TypeOf(rangeStmt.X).(*types.Map); !ok {
			return
		}

		ident, ok := rangeStmt.Key.(*ast.Ident)
		if !ok || ident == nil {
			return
		}
		assignStmt, ok := ident.Obj.Decl.(*ast.AssignStmt)
		if !ok || assignStmt == nil {
			return
		}

		keyIdent := getKeyIdent(assignStmt)
		if keyIdent == nil {
			return
		}
		reportUsingDeleteFuncWithIterationKey(pass, rangeStmt.Body, keyIdent)
	})

	return nil, nil
}

func getKeyIdent(stmt *ast.AssignStmt) *ast.Ident {
	if len(stmt.Lhs) == 0 {
		return nil
	}

	ident, ok := stmt.Lhs[0].(*ast.Ident)
	if !ok || ident == nil {
		return nil
	}

	if ident.Obj.Kind == ast.Var {
		return ident
	}

	return nil
}

func reportUsingDeleteFuncWithIterationKey(pass *analysis.Pass, body *ast.BlockStmt, keyIdent *ast.Ident) {
	sameValueKeyIdents := map[*ast.Object]*ast.Ident{keyIdent.Obj: keyIdent}

	ast.Inspect(body, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		switch n := n.(type) {
		case *ast.AssignStmt:
			updateSameValueKeyIdents(n, keyIdent, sameValueKeyIdents)
		case *ast.CallExpr:
			fun, ok := n.Fun.(*ast.Ident)
			if !ok || fun == nil {
				return true
			}

			if fun.Name != "delete" {
				return true
			}

			if len(n.Args) != 2 {
				return true
			}

			if secondArg, ok := n.Args[1].(*ast.Ident); ok {
				for _, keyIdent := range sameValueKeyIdents {
					if secondArg.Obj == keyIdent.Obj {
						return true
					}
				}
			}

			pass.Reportf(fun.Pos(), "delete function is called with a value different from range key")
		}

		return true
	})
}

func updateSameValueKeyIdents(n *ast.AssignStmt, keyIdent *ast.Ident, sameValueKeyIdents map[*ast.Object]*ast.Ident) {
	for i, r := range n.Rhs {
		ident, ok := r.(*ast.Ident)
		if !ok || ident == nil {
			continue
		}
		if ident.Obj == keyIdent.Obj {
			l := n.Lhs[i].(*ast.Ident)
			sameValueKeyIdents[l.Obj] = l
		}
	}

	for i, l := range n.Lhs {
		li, ok := l.(*ast.Ident)
		if !ok || li == nil {
			continue
		}
		if _, ok := sameValueKeyIdents[li.Obj]; ok {
			r := n.Rhs[i]
			if ri, ok := r.(*ast.Ident); ok {
				if _, exist := sameValueKeyIdents[ri.Obj]; exist {
					continue
				}
			}
			delete(sameValueKeyIdents, li.Obj)

		}
	}
}
