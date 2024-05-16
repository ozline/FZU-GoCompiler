package parser

// 基于教授提供的文法规则，这里定义了终结符和产生式

const (
	EPSILON = "" // 空字符

	SHIFT  = "shift"  // 移入
	REDUCE = "reduce" // 规约
	ACCEPT = "accept" // 接受
	ERROR  = "error"  // 错误
)

var (

	// TERMINALS 表示所有终结符
	TERMINALS = []Terminal{
		"{", "}", ";", "[", "]", "(", ")",
		"+", "-", "*", "/",
		"||", "&&", "==", "!=", "<", "<=", ">", ">=", "!", "=", "!=",
		"if", "else", "while", "do", "break",
		"true", "false",
		"basic", "id", "num", "real",
		"int", "bool", "string", "float", "byte",
		EPSILON, TERMINATE_SYMBOL,
	}

	// ARGUMENTED_PRODUCTION 表示增广产生式
	ARGUMENTED_PRODUCTION = Production{"program'", []Symbol{"program"}}

	// TERMINATE_SYMBOL 表示终结符
	TERMINATE_SYMBOL = Terminal("$")

	// PRODUCTIONS 表示所有产生式
	// 每一个 PRODUCTION 表示一个产生式，其中 Head 表示产生式的头部，Body 表示产生式的体部
	// 例如，对于产生式 E -> E + T，E 是头部，E + T 是体部
	PRODUCTIONS = []Production{
		{"program", []Symbol{"block"}},
		{"block", []Symbol{"{", "decls", "stmts", "}"}},
		{"decls", []Symbol{"decls", "decl"}},
		{"decls", []Symbol{EPSILON}},
		{"decl", []Symbol{"type", "id", ";"}},
		// {"type", []Symbol{"type[num]"}},
		{"type", []Symbol{"type_array"}},                  // 新增产生式
		{"type_array", []Symbol{"type", "[", "num", "]"}}, // 新增产生式
		{"type", []Symbol{"basic"}},
		{"stmts", []Symbol{"stmts", "stmt"}},
		{"stmts", []Symbol{EPSILON}},
		{"stmt", []Symbol{"loc", "=", "bool", ";"}},
		{"stmt", []Symbol{"loc", "=", "num", ";"}},
		// {"stmt", []Symbol{"loc", "=", "id", ";"}},
		{"stmt", []Symbol{"if", "(", "bool", ")", "stmt"}},
		{"stmt", []Symbol{"if", "(", "bool", ")", "stmt", "else", "stmt"}},
		{"stmt", []Symbol{"while", "(", "bool", ")", "stmt"}},
		{"stmt", []Symbol{"do", "stmt", "while", "(", "bool", ")", ";"}},
		{"stmt", []Symbol{"break", ";"}},
		{"stmt", []Symbol{"block"}},
		// {"loc", []Symbol{"loc[num]"}},
		{"loc", []Symbol{"loc_array"}},                  // 新增产生式
		{"loc_array", []Symbol{"loc", "[", "num", "]"}}, // 新增产生式
		// {"loc", []Symbol{"loc", "[", "num", "]"}},
		{"loc", []Symbol{"id"}},
		{"bool", []Symbol{"bool", "||", "join"}},
		{"bool", []Symbol{"join"}},
		{"join", []Symbol{"join", "&&", "equality"}},
		{"join", []Symbol{"equality"}},
		{"equality", []Symbol{"equality", "==", "rel"}},
		{"equality", []Symbol{"equality", "!=", "rel"}},
		{"equality", []Symbol{"rel"}},
		{"rel", []Symbol{"expr", "<", "expr"}},
		{"rel", []Symbol{"expr", "<=", "expr"}},
		{"rel", []Symbol{"expr", ">=", "expr"}},
		{"rel", []Symbol{"expr", ">", "expr"}},
		{"rel", []Symbol{"expr"}},
		{"expr", []Symbol{"expr", "+", "term"}},
		{"expr", []Symbol{"expr", "-", "term"}},
		{"expr", []Symbol{"term"}},
		{"term", []Symbol{"term", "*", "unary"}},
		{"term", []Symbol{"term", "/", "unary"}},
		{"term", []Symbol{"unary"}},
		{"unary", []Symbol{"!", "unary"}},
		{"unary", []Symbol{"-", "unary"}},
		{"unary", []Symbol{"factor"}},
		{"factor", []Symbol{"(", "bool", ")"}},
		{"factor", []Symbol{"loc"}},
		{"factor", []Symbol{"num"}},
		{"factor", []Symbol{"real"}},
		{"factor", []Symbol{"true"}},
		{"factor", []Symbol{"false"}},
	}
)

/*
这个文法似乎存在 SHIFT/REDUCE 冲突。

运行时提示：设置 SHIFT 发生冲突! 状态: 113 符号: 'else', 期望数字：121, 现有内容为：{ActionType:reduce Number:11}

检查了一圈没发现什么问题，然后问了以下 gpt，说是悬挂 else 问题

stmt 规则定义了 if 语句和 if-else 语句，这就导致了解析器在遇到 else 关键字时无法确定是应该将 else 与前面最近的尚未匹配的 if 配对（即SHIFT动作），还是应该结束当前的 if 语句的解析（即REDUCE动作）

然后发现书上就有这玩意，了解了一下是用扩展巴格斯范式（EBNF）解决的，即在 stmt 规则中将 if-else 语句和 if 语句分开，这样解析器就可以根据 else 关键字前的内容来判断是应该进行 SHIFT 还是 REDUCE 了

解决这个的方案是细化文法，似乎可以写成这样，实际上就是细化一个 matched_stmt 和 unmatched_stmt：
program → block
block → { decls stmts }
decls → decls decl | ε
decl → type id;
type → type[num] | basic
stmts → stmts stmt | ε
stmt → matched_stmt | unmatched_stmt
matched_stmt → loc=bool;
             | while(bool) matched_stmt
             | do matched_stmt while(bool);
             | break;
             | block
             | if (bool) matched_stmt else matched_stmt
unmatched_stmt → if (bool) stmt
               | if (bool) matched_stmt else unmatched_stmt
loc → loc[num] | id
bool → bool || join | join
join → join ＆ equality | equality
equality → equality == rel | equality ！= rel | rel
rel → expr<expr | expr<=expr | expr>=expr | expr>expr | expr
expr → expr+term | expr-term | term
term → term*unary | term/unary | unary
unary → !unary | -unary | factor
factor → (bool) | loc | num | real | true | false
*/
