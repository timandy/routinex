package astutil

import (
	"go/ast"
	"go/printer"
	"go/token"
	"os"
)

func CreateField(name string, typ string) *ast.Field {
	return &ast.Field{Names: []*ast.Ident{ast.NewIdent(name)}, Type: ast.NewIdent(typ)}
}

func CreateAssignNilStmt(x ast.Expr, name string) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.SelectorExpr{X: x, Sel: &ast.Ident{Name: name}}},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{&ast.Ident{Name: "nil"}},
	}
}

func SaveAs(path string, fset *token.FileSet, af *ast.File) {
	// create dest file
	destFile, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer destFile.Close()
	// write code to dest file
	if err = printer.Fprint(destFile, fset, af); err != nil {
		panic(err)
	}
}
