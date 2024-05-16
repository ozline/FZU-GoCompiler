package main

import (
	"fmt"
	"os"

	"github.com/ozline/CoursePractice-GoCompiler/lexer"
	"github.com/ozline/CoursePractice-GoCompiler/parser"
)

func main() {
	// 创建一个新的 parser 实例
	parser := parser.NewParser()

	// 构造并打印 First 集合
	parser.InitFirstSet()
	// parser.PrintFirstSet()

	// // 构造并打印 Follow 集合
	// parser.InitFollowSet()
	// // parser.PrintFollowSet()

	// 检查文法是否存在左递归
	// parser.Grammar.CheckLeftRecursion()

	// 构建状态集合并输出
	parser.BuildStateCollection()
	parser.PrintStateCollection()

	// 构建分析表
	parser.BuildTables()
	// parser.PrintGoToTable()
	// parser.PrintActionTable()

	// 初始化词法分析器
	var lex *lexer.Lexer
	if len(os.Args) > 1 {
		lex = lexer.NewLexer(os.Stdin)
	} else {
		file, err := os.Open("/Users/ozliinex/projects/code-compiler/tests/case5.in")
		if err != nil {
			fmt.Printf("Failed to open file: %v", err)
			return
		}
		defer file.Close()
		lex = lexer.NewLexer(file)
	}

	if err := parser.Parse(lex); err != nil {
		fmt.Printf("%v", err)
	}

	parser.PrintThreeAddress() // 打印三地址码
	parser.SymbolTable.Print() // 打印符号表
}
