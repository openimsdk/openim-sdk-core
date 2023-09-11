package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

func main() {
	//filePath := "/Users/gordon/GolandProjects/private/fork/openim-sdk-core/internal/conversation_msg" // 替换为你的Go文件路径

	filePath := "D:\\Goland\\fg\\openim-sdk-core\\internal\\conversation_msg"
	// 加载包信息
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedSyntax | packages.NeedFiles | packages.NeedTypesInfo |
			packages.NeedImports | packages.NeedName | packages.NeedDeps,
		Tests: true,
		Dir:   filePath,
	}
	pkgs, err := packages.Load(cfg, "")
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) == 0 {
		log.Fatal("Failed to load packages")
	}

	// 遍历文件中的所有声明
	for _, pkg := range pkgs {
		// for _, file := range pkg.GoFiles {
		// 	fmt.Println("file name", filepath.Base(file))

		// }
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				// 仅处理函数声明
				if fdecl, ok := decl.(*ast.FuncDecl); ok {
					// 仅处理导出的函数
					if fdecl.Name.IsExported() {
						funcName := fdecl.Name.Name
						fmt.Println("Function Name:", funcName)
						// 获取函数位置信息
						funcPos := fdecl.Pos()
						funcFile := pkg.Fset.Position(funcPos).Filename
						funcLine := pkg.Fset.Position(funcPos).Line
						//fmt.Println("Location:", funcFile, "Line:", funcLine)
						fmt.Println("Function Declared at:", fmt.Sprintf("%s:%d", funcFile, funcLine))
						if filepath.Base(funcFile) != "sdk.go" {
							continue
						}
						// 处理函数参数
						fmt.Println("Parameters:")
						for _, param := range fdecl.Type.Params.List {
							for _, name := range param.Names {
								fmt.Println("  Parameter Name:", name.Name)
								obj := pkg.TypesInfo.ObjectOf(name)
								stringParamType := getTypeString(param.Type, pkg.TypesInfo)
								fmt.Println("  Parameter Type:", stringParamType)
								fmt.Println("  Parameter Declared at:", getObjectPosition(obj, pkg.Fset))
								if isCustomType(stringParamType) {
									fmt.Println("  111Type Location:", pkg.PkgPath)
									typePos := getTypePosition(stringParamType, pkg)
									fmt.Println("  Type Location:", typePos, pkg.Name)
								}
							}
						}

						// 处理返回值
						fmt.Println("Return Values:")
						if fdecl.Type.Results != nil {
							for _, result := range fdecl.Type.Results.List {
								if len(result.Names) > 0 {
									for _, name := range result.Names {
										obj := pkg.TypesInfo.ObjectOf(name)
										fmt.Println("  Result Name:", name)
										fmt.Println("  Result Type:", getTypeString(result.Type, pkg.TypesInfo))
										fmt.Println("  Result Declared at:", getObjectPosition(obj, pkg.Fset))
									}
								} else {
									fmt.Println("  Result Type:", getTypeString(result.Type, pkg.TypesInfo))
								}
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
}

// 获取类型的字符串表示形式
func getTypeString(expr ast.Expr, info *types.Info) string {
	return info.TypeOf(expr).String()
}

// 获取类型声明的位置
func getObjectPosition(obj types.Object, fset *token.FileSet) string {
	if obj == nil {
		return "unknown"
	}
	pos := fset.Position(obj.Pos())
	return fmt.Sprintf("%s:%d", pos.Filename, pos.Line)
}

// // 获取类型的位置信息
// // 获取类型的位置信息
//
//	func getTypePosition(typeName string, pkg *packages.Package) string {
//		for _, file := range pkg.Syntax {
//			var typePos token.Pos
//			ast.Inspect(file, func(node ast.Node) bool {
//				if typeSpec, ok := node.(*ast.TypeSpec); ok {
//					if typeSpec.Name.Name == typeName {
//						typePos = typeSpec.Pos()
//						return false // 停止继续遍历
//					}
//				}
//				return true // 继续遍历
//			})
//			if typePos != token.NoPos {
//				typeFile := pkg.Fset.Position(typePos).Filename
//				typeLine := pkg.Fset.Position(typePos).Line
//				return fmt.Sprintf("%s:%d", typeFile, typeLine)
//			}
//		}
//		return ""
//	}
//

func getLastSegment(typeName string, separator string) string {
	index := strings.Index(typeName, separator)
	if index == -1 {
		return typeName
	}
	return typeName[index+len(separator):]
}

// 没有处理当前包的情况
// 获取类型的位置信息
func getTypePosition(typeName string, pkg *packages.Package) string {
	for _, pkgInfo := range pkg.Imports {
		// fmt.Println("compare::::", pkgInfo.Name)
		for _, file := range pkgInfo.Syntax {
			var typePos token.Pos
			// fmt.Println("compare file::::", pkgInfo.Name)

			ast.Inspect(file, func(node ast.Node) bool {
				if typeSpec, ok := node.(*ast.TypeSpec); ok {
					if pkgInfo.Name == "sdk_struct" && typeSpec.Name.Name == "MsgStruct" && typeName == "*github.com/openimsdk/openim-sdk-core/v3/sdk_struct.MsgStruct" {
						fmt.Println("msg struct", typeSpec.Name.Name, fmt.Sprintf("%s.%s", pkgInfo.PkgPath, typeSpec.Name.Name))
					}
					if fmt.Sprintf("%s.%s", pkgInfo.PkgPath, typeSpec.Name.Name) == getLastSegment(typeName, "*") {
						fmt.Println("compare::::", getTypeString(typeSpec.Type, pkgInfo.TypesInfo), typeName)
						typePos = typeSpec.Pos()
						return false // 停止继续遍历
					}
				}
				// fmt.Println("not compare::::")
				return true // 继续遍历
			})
			if typePos != token.NoPos {
				typeFile := pkg.Fset.Position(typePos).Filename
				typeLine := pkg.Fset.Position(typePos).Line
				return fmt.Sprintf("%s:%d", typeFile, typeLine)
			}
		}
	}
	return ""
}

func isCustomType(typeName string) bool {
	basicTypes := map[string]bool{
		"bool":            true,
		"byte":            true,
		"complex64":       true,
		"complex128":      true,
		"float32":         true,
		"float64":         true,
		"int":             true,
		"int8":            true,
		"int16":           true,
		"int32":           true,
		"int64":           true,
		"rune":            true,
		"string":          true,
		"uint":            true,
		"uint8":           true,
		"uint16":          true,
		"uint32":          true,
		"uint64":          true,
		"uintptr":         true,
		"context.Context": true,
		"error":           true,
		"[]string":        true,
	}

	return !basicTypes[typeName]
}
