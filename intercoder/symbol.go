// symbol.go
// 符号表的定义

package intercoder

import (
	"fmt"
	"math/rand"
	"strconv"
)

// SymbolType 定义了符号的类型
type SymbolType string

const (
	// SymbolTypeVar 表示一个变量
	SymbolTypeVar SymbolType = "VAR"
	// SymbolTypeFunc 表示一个函数
	SymbolTypeFunc SymbolType = "FUNC"
	// SymbolTypeArray 表示一个数组
	SymbolTypeArray SymbolType = "ARRAY"
	// SymbolTypeBasic 表示一个基本类型
	SymbolTypeBasic SymbolType = "BASIC"
)

// Position 用于表示符号在源代码中的位置
type Position struct {
	Line   int
	Column int
}

// SymbolInfo 用于表示符号的信息
type SymbolInfo struct {
	Type      SymbolType
	Positions []Position
	Name      string // 符号的名字
	Scope     int    // 符号的作用域级别
	Addr      string // 符号在中间代码中的地址或名称
	ArraySize int    // (如果是数组) 数组的大小
	DataType  string // (如果是变量) 数据类型
}

// SymbolTable 表示符号表
type SymbolTable struct {
	table map[string]SymbolInfo // 使用符号的名字作为键
	scope int                   // 当前的作用域级别
}

// NewSymbolTable 创建一个新的符号表
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		table: make(map[string]SymbolInfo),
		scope: 0,
	}
}

// EnterScope 进入一个新的作用域
func (st *SymbolTable) EnterScope() {
	st.scope++
}

// ExitScope 离开当前作用域
func (st *SymbolTable) ExitScope() {
	if st.scope > 0 {
		st.scope--
	}
	// 清除当前作用域的所有符号
	for name, symbol := range st.table {
		if symbol.Scope == st.scope {
			delete(st.table, name)
		}
	}
}

func (st *SymbolTable) NewTempAddr() string {
	return "t" + strconv.Itoa(rand.Intn(1000)) // 生成一个随机地址，没必要上纲上线

}

// Define 定义一个新的符号
func (st *SymbolTable) Define(name string, symbolType SymbolType) error {
	if _, exists := st.table[name]; exists {
		return fmt.Errorf("symbol %s already defined", name)
	}

	st.table[name] = SymbolInfo{
		Name:  name,
		Type:  symbolType,
		Scope: st.scope,
		Addr:  st.NewTempAddr(),
	}
	return nil
}

// DefineData 定义一个变量
// 如果不是数组，size 输入 0
// 如果不是基本数据类型 basic，dataType 输入 ""
func (st *SymbolTable) DefineData(name string, dataType string, size int) error {
	if _, exists := st.table[name]; exists {
		return fmt.Errorf("symbol %s already defined", name)
	}

	var symbolType SymbolType

	if size != 0 {
		symbolType = SymbolTypeArray // 如果是数组，类型改为数组
	} else {
		symbolType = SymbolTypeVar // 如果不是数组，类型改为变量
	}

	st.table[name] = SymbolInfo{
		Name:      name,
		Type:      symbolType,
		Scope:     st.scope,
		Addr:      st.NewTempAddr(),
		DataType:  dataType,
		ArraySize: size,
	}
	return nil
}

// Lookup 查找一个符号
func (st *SymbolTable) Lookup(name string) (SymbolInfo, bool) {
	symbol, exists := st.table[name]
	return symbol, exists
}

// Print 打印符号表的内容，用于调试
func (st *SymbolTable) Print() {
	fmt.Println("\n\n===============符号表===============")
	for _, symbol := range st.table {
		fmt.Printf("名称: %s, 类型: %s, 作用域: %d 地址: %s\n", symbol.Name, symbol.Type, symbol.Scope, symbol.Addr)
	}
}
