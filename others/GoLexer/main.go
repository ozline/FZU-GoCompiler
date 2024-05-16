package main

import (
	"fmt"
	"os"
)

func main() {
	// 创建一个新的 lexer 实例
	lexer := NewLexer(os.Stdin)

	for {
		token, err := lexer.NextToken()
		if err != nil {
			fmt.Println(err)
			continue
		}
		if token.Type == EOF {
			break
		}
		fmt.Printf("Type: %-10v Value: %v\n", tokenTypes[token.Type], token.Value) // 输出 Token 信息
	}
	fmt.Printf("\n\n\n")
	lexer.PrintSymbols()
}
