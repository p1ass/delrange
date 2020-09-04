package delrange

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "delrange is a static analysis tool which detects references to loop iterator variable."

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
		(*ast.ForStmt)(nil),
		(*ast.RangeStmt)(nil),
	}

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.RangeStmt:
			checkRangeStmt(pass, n)
		}
	})

	return nil, nil
}

func checkRangeStmt(pass *analysis.Pass, rangeStmt *ast.RangeStmt) {
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
	reportUsingIterVarRef(pass, rangeStmt.Body, keyIdent)
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

func reportUsingIterVarRef(pass *analysis.Pass, body *ast.BlockStmt, keyIdent *ast.Ident) {
	ast.Inspect(body, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		callExpr, ok := n.(*ast.CallExpr)
		if !ok || callExpr == nil {
			return true
		}

		fun, ok := callExpr.Fun.(*ast.Ident)
		if !ok || fun == nil {
			return true
		}

		if fun.Name != "delete" {
			return true
		}

		if len(callExpr.Args) != 2 {
			return true
		}

		if secondArg, ok := callExpr.Args[1].(*ast.Ident); ok && secondArg.Obj == keyIdent.Obj {
			return true
		}

		pass.Reportf(fun.Pos(), "delete function is called with a value different from range key")

		return true
	})
}
