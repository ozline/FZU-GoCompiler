package main

type Position struct {
	Line   int
	Column int
}

type SymbolInfo struct {
	Type      TokenType
	Positions []Position
}

type SymbolTable map[string]SymbolInfo
