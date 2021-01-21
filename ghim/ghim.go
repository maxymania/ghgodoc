package ghim

import (
	"go/ast"
	"go/doc"
	"go/build"
	"go/token"
	"go/parser"
	
	"os"
	"strings"
	//"fmt"
)

func filecheck(file os.FileInfo) bool {
	name := file.Name()
	
	if strings.HasPrefix(name,".") { return false }
	if !strings.HasSuffix(name,".go") { return false }
	if strings.HasSuffix(name,"_test.go") { return false }
	return true
}

func Parse(bPkg *build.Package, impPath string) (fset *token.FileSet, pkg *doc.Package, err error) {
	fset = token.NewFileSet()
	pkgSet,err0 := parser.ParseDir(fset, bPkg.Dir, filecheck, parser.ParseComments)
	err = err0
	if err!=nil { return }
	
	var aPkg *ast.Package
	first := ""
	for k,v := range pkgSet {
		if aPkg!=nil && first!="main" && first!="documentation" { break }
		first = k
		aPkg = v
	}
	
	if aPkg!=nil { pkg = doc.New(aPkg,impPath,0) }
	
	return
}

