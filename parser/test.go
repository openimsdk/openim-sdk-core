package main

import (
	"fmt"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/packages"
	"log"
)

func main() {
	filePath := "D:\\Goland\\fg\\openim-sdk-core\\parser\\" // 替换为你的Go文件路径

	// 加载包信息
	cfg := &packages.Config{
		Mode:  packages.NeedTypes | packages.NeedSyntax | packages.NeedFiles,
		Tests: true,
		Dir:   filePath,
	}
	fmt.Println("111111111111111111")
	pkgs, err := packages.Load(cfg, "")
	if err != nil {
		fmt.Println("333333333333", err)
		log.Fatal(err)
	}
	if len(pkgs) == 0 {
		log.Fatal("Failed to load packages")
	}
	fmt.Println("222222222222222222222222")
	// 遍历文件中的所有声明
	for _, file := range pkgs[0].Syntax {
		for _, decl := range file.Decls {
			// 仅处理函数声明
			if fdecl, ok := decl.(*ast.FuncDecl); ok {
				// 仅处理导出的函数
				if fdecl.Name.IsExported() {
					funcName := fdecl.Name.Name
					fmt.Println("Function Name:", funcName)

					// 处理函数参数
					fmt.Println("Parameters:")
					for _, param := range fdecl.Type.Params.List {
						for _, name := range param.Names {
							fmt.Println("  Parameter Name:", name.Name)
						}
						fmt.Println("  Parameter Type:", param.Type)

						//paramType := getTypeString(param.Type, pkgs[0].TypesInfo)
						//fmt.Println("  Parameter Type:", paramType)
					}

					// 处理返回值
					fmt.Println("Return Values:")
					if fdecl.Type.Results != nil {
						for _, result := range fdecl.Type.Results.List {
							//resultType := getTypeString(result.Type, pkgs[0].TypesInfo)
							//fmt.Println("  Result Type:", resultType)
							fmt.Println("  Result Type:", result.Type)

						}
					} else {
						fmt.Println("  No return values")
					}

					fmt.Println("----------------------")
				}
			}
		}
	}
}

// 获取类型的字符串表示形式
func getTypeString(expr ast.Expr, info *types.Info) string {
	return info.TypeOf(expr).String()
}
