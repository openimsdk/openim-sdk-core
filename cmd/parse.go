package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
)

func main() {

	//filePath, err := filepath.Abs(".\test.go")
	//if err != nil {
	//	panic(err)
	//}
	// 解析Go文件
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "D:\\Goland\\workspace\\Open-IM-SDK-Core\\main\\test.go", nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}
	//myImporter := importer.Default()
	// 创建类型检查器
	//conf := types.Config{Importer: myImporter}
	info := &types.Info{
		Defs: make(map[*ast.Ident]types.Object),
	}
	//// 类型检查
	//_, err = conf.Check("", fset, []*ast.File{node}, info)
	//if err != nil {
	//	panic(err)
	//}
	// 遍历文件中所有函数

	fn := func(pkg *types.Package) string {
		return pkg.Name()
	}
	for _, decl := range node.Decls {
		if f, ok := decl.(*ast.FuncDecl); ok {
			// 打印函数名
			fmt.Println("Function Name: ", f.Name.Name)
			// 打印参数名和类型
			for _, param := range f.Type.Params.List {
				for _, name := range param.Names {
					obj := info.ObjectOf(name)
					typ := obj.Type()
					fmt.Printf("Parameter Name: %s, Type: %s\n", name.Name, types.TypeString(typ, fn))
				}
			}
		}
	}
}
