package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"log"
	"path/filepath"
	"strings"
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
	} else {
		log.Println("package is:", pkgs)
	}
	//time.Sleep(time.Second * 100)

	// 遍历所有的包
	for _, pkg := range pkgs {
		//遍历包中所有语法树
		for i, file := range pkg.Syntax {

			filename := pkg.GoFiles[i]
			log.Println("file name is", filename)
			if filepath.Base(filename) != "sdk.go" {
				continue
			}
			log.Println("file name is", filename)
			//遍历一个文件所有的申明
			for _, decl := range file.Decls {
				// 仅处理函数声明
				if fdecl, ok := decl.(*ast.FuncDecl); ok {

					// 仅处理导出的函数
					if fdecl.Name.IsExported() {
						//获取函数名字
						fmt.Println("Function Name:", fdecl.Name.Name)
						//获取函数注释
						fmt.Printf("Comments:\n%s", getFunComments(fdecl))

						//获取函数原型
						fmt.Printf("Function Declaration:\n%s\n", getFunProtoType(fdecl))

						// 获取函数位置信息
						funcPos := fdecl.Pos()
						funcFile := pkg.Fset.Position(funcPos).Filename
						funcLine := pkg.Fset.Position(funcPos).Line
						fmt.Println("Function Declared at:", fmt.Sprintf("%s:%d", funcFile, funcLine))

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
									//fmt.Printf("obj.Type() is of type: %T\n", obj.Type())
									//parsmProtoType := getParamsProtoType(obj.Type(), pkg)
									//fmt.Println("  Type ProtoType:", parsmProtoType)
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

func getFunProtoType(fdecl *ast.FuncDecl) string {
	var sb strings.Builder

	// 如果是方法，输出接收器
	if fdecl.Recv != nil && len(fdecl.Recv.List) > 0 {
		sb.WriteString("func (")
		for i, field := range fdecl.Recv.List {
			sb.WriteString(getFieldDeclaration(field))
			if i != len(fdecl.Recv.List)-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString(") ")
	} else {
		sb.WriteString("func ")
	}

	// 输出函数名
	sb.WriteString(fdecl.Name.Name)

	// 输出函数参数
	sb.WriteString("(")
	for i, field := range fdecl.Type.Params.List {
		sb.WriteString(getFieldDeclaration(field))
		if i != len(fdecl.Type.Params.List)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")

	// 输出返回值
	if fdecl.Type.Results != nil {
		sb.WriteString(" (")
		for i, field := range fdecl.Type.Results.List {
			sb.WriteString(getFieldDeclaration(field))
			if i != len(fdecl.Type.Results.List)-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString(")")
	}
	return sb.String()
}
func getParamsProtoType(t types.Type, pkg *packages.Package) string {
	// 循环处理，因为可能存在多级的指针，例如 **MyType
	for {
		if ptrType, ok := t.(*types.Pointer); ok {
			t = ptrType.Elem()
			continue
		} else if sliceType, ok := t.(*types.Slice); ok {
			t = sliceType.Elem()
			continue
		}
		break
	}

	// 接下来的代码是处理 *types.Named 的
	named, ok := t.(*types.Named)
	if !ok {
		return ""
	}

	obj := named.Obj()
	typeName := obj.Name()

	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if typeSpec.Name.Name == typeName {
					// 获取类型的源代码，包括注释
					var buf bytes.Buffer
					if genDecl.Doc != nil {
						buf.WriteString(genDecl.Doc.Text())
					}
					err := format.Node(&buf, pkg.Fset, genDecl)
					if err != nil {
						log.Fatalf("Failed to format node: %v", err)
					}
					return buf.String()
				}
			}
		}
	}
	return ""
}

func getFunComments(fdecl *ast.FuncDecl) string {
	var comments string
	if fdecl.Doc != nil {
		for _, c := range fdecl.Doc.List {
			comments += c.Text + "\n"
		}
	}
	return comments
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
func getFieldDeclaration(field *ast.Field) string {
	var sb strings.Builder

	for i, name := range field.Names {
		sb.WriteString(name.Name)
		if i != len(field.Names)-1 {
			sb.WriteString(", ")
		}
	}
	if len(field.Names) > 0 {
		sb.WriteString(" ")
	}
	sb.WriteString(getExprString(field.Type))

	return sb.String()
}

func getExprString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + getExprString(t.X)
	case *ast.SelectorExpr:
		return getExprString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + getExprString(t.Elt)
	case *ast.MapType:
		return "map[" + getExprString(t.Key) + "]" + getExprString(t.Value)
	default:
		return fmt.Sprintf("%T", expr) // 用于调试，显示AST类型
	}
}

type Content struct {
	FuncName string `json:"funcName"`
}
