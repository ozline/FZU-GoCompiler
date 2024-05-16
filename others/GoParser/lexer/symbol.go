// symbol.go
// 符号表的定义

package lexer

// Position 用于表示符号在源代码中的位置
type Position struct {
	Line   int
	Column int
}

// SymbolInfo 用于表示符号的信息
type SymbolInfo struct {
	Type      TokenType
	Positions []Position
}

// SymbolTable 符号表
type SymbolTable map[string]SymbolInfo
