package main

// Token 类型
type TokenType int

const (
	EOF           TokenType = iota
	IDENTIFIER              // 标识符
	NUMBER                  // 数字
	RESERVED_WORD           // 保留字
	OPERATOR                // 运算符
	DELIMITER               // 分隔符
	PACKAGE                 // 包
	IMPORT                  // 导入
	STRING                  // 字符串
)

var tokenTypes = map[TokenType]string{
	EOF:           "文件结束符",
	IDENTIFIER:    "标识符",
	NUMBER:        "数字",
	RESERVED_WORD: "保留字",
	OPERATOR:      "运算符",
	DELIMITER:     "分隔符",
	PACKAGE:       "包",
	IMPORT:        "导入",
	STRING:        "字符串",
}

// Token 结构体
type Token struct {
	Type  TokenType
	Value string
}

// 保留字列表
var reservedWords = map[string]TokenType{
	// 25 个保留关键字
	"if":          RESERVED_WORD,
	"else":        RESERVED_WORD,
	"for":         RESERVED_WORD,
	"func":        RESERVED_WORD,
	"return":      RESERVED_WORD,
	"var":         RESERVED_WORD,
	"break":       RESERVED_WORD,
	"continue":    RESERVED_WORD,
	"package":     RESERVED_WORD,
	"import":      RESERVED_WORD,
	"chan":        RESERVED_WORD,
	"const":       RESERVED_WORD,
	"default":     RESERVED_WORD,
	"defer":       RESERVED_WORD,
	"fallthrough": RESERVED_WORD,
	"case":        RESERVED_WORD,
	"goto":        RESERVED_WORD,
	"interface":   RESERVED_WORD,
	"map":         RESERVED_WORD,
	"range":       RESERVED_WORD,
	"select":      RESERVED_WORD,
	"struct":      RESERVED_WORD,
	"switch":      RESERVED_WORD,
	"type":        RESERVED_WORD,
	"go":          RESERVED_WORD,

	// 类型
	"int":        RESERVED_WORD,
	"bool":       RESERVED_WORD,
	"string":     RESERVED_WORD,
	"float":      RESERVED_WORD,
	"float32":    RESERVED_WORD,
	"float64":    RESERVED_WORD,
	"byte":       RESERVED_WORD,
	"rune":       RESERVED_WORD,
	"uint":       RESERVED_WORD,
	"uint8":      RESERVED_WORD,
	"uint16":     RESERVED_WORD,
	"uint32":     RESERVED_WORD,
	"uint64":     RESERVED_WORD,
	"int8":       RESERVED_WORD,
	"int16":      RESERVED_WORD,
	"int32":      RESERVED_WORD,
	"int64":      RESERVED_WORD,
	"uintptr":    RESERVED_WORD,
	"complex":    RESERVED_WORD,
	"complex64":  RESERVED_WORD,
	"complex128": RESERVED_WORD,
	"error":      RESERVED_WORD,

	// 布尔类型常量
	"true":  RESERVED_WORD,
	"false": RESERVED_WORD,
}

// 运算符列表
var operators = map[string]TokenType{
	"+":  OPERATOR,
	"-":  OPERATOR,
	"*":  OPERATOR,
	"/":  OPERATOR,
	"=":  OPERATOR,
	"==": OPERATOR,
	"!=": OPERATOR,
	"<":  OPERATOR,
	">":  OPERATOR,
	"&&": OPERATOR,
	"||": OPERATOR,
	"!":  OPERATOR,
	// 添加其他运算符...
}

// 分隔符列表
var delimiters = map[string]TokenType{
	";": DELIMITER,
	",": DELIMITER,
	"(": DELIMITER,
	")": DELIMITER,
	"{": DELIMITER,
	"}": DELIMITER,
	"[": DELIMITER,
	"]": DELIMITER,
	":": DELIMITER,
	".": DELIMITER,
	// 添加其他分隔符...
}
