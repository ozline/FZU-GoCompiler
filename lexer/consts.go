// consts.go
// 常量定义

package lexer

// Token 类型
type TokenType int

const (
	EOF           TokenType = iota
	IDENTIFIER              // 标识符
	NUMBER                  // 数字
	REAL                    // 实数
	RESERVED_WORD           // 保留字
	TYPE                    // 类型
	OPERATOR                // 运算符
	DELIMITER               // 分隔符
	PACKAGE                 // 包
	IMPORT                  // 导入
	STRING                  // 字符串
)

var TokenTypes = map[TokenType]string{
	EOF:           "文件结束符",
	IDENTIFIER:    "标识符",
	NUMBER:        "数字",
	REAL:          "实数",
	TYPE:          "类型",
	RESERVED_WORD: "保留字",
	OPERATOR:      "运算符",
	DELIMITER:     "分隔符",
	PACKAGE:       "包",
	IMPORT:        "导入",
	STRING:        "字符串",
}

// Token 结构体
type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

var basicTypes = map[string]bool{
	"int":    true,
	"float":  true,
	"string": true,
	"bool":   true,
	"byte":   true,
}

// 保留字列表
var reservedWords = map[string]TokenType{
	// go 的25 个保留关键字
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
	"do":          RESERVED_WORD,
	"while":       RESERVED_WORD,

	// 布尔类型常量
	"true":  RESERVED_WORD,
	"false": RESERVED_WORD,

	"int":    TYPE,
	"bool":   TYPE,
	"string": TYPE,
	"float":  TYPE,
	"byte":   TYPE,
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
	"&":  OPERATOR,
	"|":  OPERATOR,
	"<=": OPERATOR,
	">=": OPERATOR,
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
