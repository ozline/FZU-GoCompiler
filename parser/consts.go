package parser

import (
	"github.com/ozline/CoursePractice-GoCompiler/consts"
)

// 基于教授提供的文法规则，这里定义了终结符和产生式

const (
	EPSILON = "" // 空字符

	SHIFT  = "shift"  // 移入
	REDUCE = "reduce" // 规约
	ACCEPT = "accept" // 接受
	ERROR  = "error"  // 错误
)

var (

	// TerminalS 表示所有终结符
	TERMINALS = []consts.Terminal{
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
	ARGUMENTED_PRODUCTION = Production{"program'", []consts.Symbol{"program"}, genARGUMENTED_PRODUCTION}

	// TERMINATE_SYMBOL 表示终结符
	TERMINATE_SYMBOL = consts.Terminal("$")

	// PRODUCTIONS 表示所有产生式
	// 每一个 PRODUCTION 表示一个产生式，其中 Head 表示产生式的头部，Body 表示产生式的体部
	// 例如，对于产生式 E -> E + T，E 是头部，E + T 是体部
	PRODUCTIONS = []Production{
		0: {"program", []consts.Symbol{"block"}, genProgram},
		1: {"block", []consts.Symbol{"{", "decls", "stmts", "}"}, genBlock},
		2: {"decls", []consts.Symbol{"decls", "decl"}, genDecls},
		3: {"decls", []consts.Symbol{EPSILON}, genDeclsEpsilon},
		4: {"decl", []consts.Symbol{"type", "id", ";"}, genDecl},
		// {"type", []consts.Symbol{"type[num]"}},
		5:  {"type", []consts.Symbol{"type_array"}, genTypeArray},                       // 新增产生式
		6:  {"type_array", []consts.Symbol{"type", "[", "num", "]"}, genTypeArrayFinal}, // 新增产生式
		7:  {"type", []consts.Symbol{"basic"}, genBasicType},
		8:  {"stmts", []consts.Symbol{"stmts", "stmt"}, genStmts},
		9:  {"stmts", []consts.Symbol{EPSILON}, genStmtsEpsilon},
		10: {"stmt", []consts.Symbol{"loc", "=", "bool", ";"}, genStmt},
		11: {"stmt", []consts.Symbol{"loc", "=", "num", ";"}, genStmt},
		// {"stmt", []consts.Symbol{"loc", "=", "id", ";"}},
		12: {"stmt", []consts.Symbol{"if", "(", "bool", ")", "stmt"}, genStmtIf},
		13: {"stmt", []consts.Symbol{"if", "(", "bool", ")", "stmt", "else", "stmt"}, genStmtIfElse},
		14: {"stmt", []consts.Symbol{"while", "(", "bool", ")", "stmt"}, genStmtWhile},
		15: {"stmt", []consts.Symbol{"do", "stmt", "while", "(", "bool", ")", ";"}, genStmtDoWhile},
		16: {"stmt", []consts.Symbol{"break", ";"}, genStmtBreak},
		17: {"stmt", []consts.Symbol{"block"}, genStmtBlock},
		// {"loc", []consts.Symbol{"loc[num]"}},
		18: {"loc", []consts.Symbol{"loc_array"}, genLocArray},                       // 新增产生式
		19: {"loc_array", []consts.Symbol{"loc", "[", "num", "]"}, genLocArrayFinal}, // 新增产生式
		// {"loc", []consts.Symbol{"loc", "[", "num", "]"}},
		20: {"loc", []consts.Symbol{"id"}, genLoc},
		21: {"bool", []consts.Symbol{"bool", "||", "join"}, genBoolOr},
		22: {"bool", []consts.Symbol{"join"}, genBool},
		23: {"join", []consts.Symbol{"join", "&&", "equality"}, genJoinAnd},
		24: {"join", []consts.Symbol{"equality"}, genJoin},
		25: {"equality", []consts.Symbol{"equality", "==", "rel"}, genEqualityEqual},
		26: {"equality", []consts.Symbol{"equality", "!=", "rel"}, genEqualityNotEqual},
		27: {"equality", []consts.Symbol{"rel"}, genEquality},
		28: {"rel", []consts.Symbol{"expr", "<", "expr"}, genRelLess},
		29: {"rel", []consts.Symbol{"expr", "<=", "expr"}, genRelLessEqual},
		30: {"rel", []consts.Symbol{"expr", ">=", "expr"}, genRelGreaterEqual},
		31: {"rel", []consts.Symbol{"expr", ">", "expr"}, genRelGreater},
		32: {"rel", []consts.Symbol{"expr"}, genRel},
		33: {"expr", []consts.Symbol{"expr", "+", "term"}, genExprAdd},
		34: {"expr", []consts.Symbol{"expr", "-", "term"}, genExprSub},
		35: {"expr", []consts.Symbol{"term"}, genExpr},
		36: {"term", []consts.Symbol{"term", "*", "unary"}, genTermMul},
		37: {"term", []consts.Symbol{"term", "/", "unary"}, genTermDiv},
		38: {"term", []consts.Symbol{"unary"}, genTerm},
		39: {"unary", []consts.Symbol{"!", "unary"}, genUnaryNeg},
		40: {"unary", []consts.Symbol{"-", "unary"}, genUnaryNot},
		41: {"unary", []consts.Symbol{"factor"}, genUnary},
		42: {"factor", []consts.Symbol{"(", "bool", ")"}, genFactorBool},
		43: {"factor", []consts.Symbol{"loc"}, genFactorLoc},
		44: {"factor", []consts.Symbol{"num"}, genFactorNum},
		45: {"factor", []consts.Symbol{"real"}, genFactorReal},
		46: {"factor", []consts.Symbol{"true"}, genFactorTrue},
		47: {"factor", []consts.Symbol{"false"}, genFactorFalse},
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
