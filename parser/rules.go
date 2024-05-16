package parser

import (
	"fmt"
	"strconv"

	"github.com/ozline/CoursePractice-GoCompiler/consts"
)

// program' -> program 增广文法，不需要编写函数
func genARGUMENTED_PRODUCTION(p *Parser) error { return nil }

// program → block
func genProgram(p *Parser) error {
	// 这个规则不需要生成任何中间代码，因为它只是一个开始符号
	return nil
}

// block → { decls stmts }
func genBlock(p *Parser) error {
	// 这个规则不需要生成中间代码，因为它是一个结构性的规则
	return nil
}

// decls → decls decl
func genDecls(p *Parser) error {
	// 由于声明不涉及运行时计算，所以这里不生成中间代码
	return nil
}

// decls → ε
func genDeclsEpsilon(p *Parser) error {
	// 空规则，不生成任何代码
	return nil
}

// decl → type id;
func genDecl(p *Parser) error {
	// 获取变量类型
	varName := p.TokenStack[len(p.TokenStack)-2]
	// 使用暂存的类型信息来重新定义变量
	if err := p.SymbolTable.DefineData(string(varName), string(p.LastType), p.LastSize); err != nil {
		fmt.Println("定义符号时出错:", err)
	}
	fmt.Printf("[符号表] 定义变量 %s 类型为 %s size:%d\n", varName, p.LastType, p.LastSize)

	// 重置 lastType 和 lastSize
	p.LastType = consts.Terminal("") // 重置类型为空
	p.LastSize = 4                   // 默认大小为 4
	return nil
}

// type → type_array
func genTypeArray(p *Parser) error {
	// 这里不需要生成中间代码，因为genLocArrayFinal会处理具体的数组赋值
	return nil
}

// type_array -> type[num]
func genTypeArrayFinal(p *Parser) error {
	arraySizeToken := p.TokenStack[len(p.TokenStack)-2]
	arraySize, err := strconv.Atoi(string(arraySizeToken)) // 转换为整数
	if err != nil {
		return fmt.Errorf("[符号表] 数组大小转换错误")
	}
	p.LastSize = arraySize
	fmt.Println("[符号表] 触发数组类型定义， 数组大小为", arraySize)
	return nil
}

// type → basic
func genBasicType(p *Parser) error {
	varType := p.TokenStack[len(p.TokenStack)-1]
	p.LastType = consts.Terminal(varType)
	fmt.Println("[符号表] 触发类型定义，该类型为", varType)
	return nil
}

// stmts -> stmts stmt
func genStmts(p *Parser) error {
	// 这部分没弄明白，问 LLM 的 OvO
	// 由于 stmts -> stmts stmt 规则表示一个语句序列，我们这里可能不需要生成特定的代码。
	// 通常，每个 stmt 在被解析时会自己生成必要的代码。
	// 因此，这里返回 nil 表示没有错误，也不需要额外的操作。
	return nil
}

// stmts -> ε
func genStmtsEpsilon(p *Parser) error {
	// 对于空语句（epsilon），同样不需要生成代码。
	return nil
}

// stmt -> loc=bool;
func genStmt(p *Parser) error {
	// 获取变量名
	varName := p.LastName
	// 获取变量值
	varValue := p.TokenStack[len(p.TokenStack)-2]
	fmt.Printf("[符号表] 将变量 %s 赋值为 %s\n", varName, varValue)

	// 生成赋值代码
	p.Emit(varName, "", string(varValue)) // 赋值操作
	return nil
}

// stmt → if(bool) stmt
func genStmtIf(p *Parser) error {
	// 生成 if 语句的代码
	condition := p.TokenStack[len(p.TokenStack)-3]
	afterStmtLabel := p.NewLabel() // 创建新的标签，用于 if 之后的代码位置

	// 生成条件为假时跳转的代码
	p.Emit("ifFalse", string(condition), "goto", afterStmtLabel)

	// 处理 if 语句的 stmt 部分
	if err := genStmt(p); err != nil {
		return err
	}

	// 标记 if 语句之后的代码位置
	p.EmitLabel(afterStmtLabel)

	return nil
}

// stmt -> if(bool) stmt else stmt
func genStmtIfElse(p *Parser) error {
	// 生成 if-else 语句的代码
	condition := p.TokenStack[len(p.TokenStack)-5]
	elseLabel := p.NewLabel()      // 创建新的标签，用于 else 语句的代码位置
	afterStmtLabel := p.NewLabel() // 创建新的标签，用于 if-else 之后的代码位置

	// 生成条件为假时跳转到 else 的代码
	p.Emit("ifFalse", string(condition), "goto", elseLabel)

	// 处理 if 语句的 stmt 部分
	if err := genStmt(p); err != nil {
		return err
	}

	// 从 if 直接跳转到 if-else 语句之后的代码位置
	p.Emit("goto", afterStmtLabel)

	// 标记 else 语句的开始位置
	p.EmitLabel(elseLabel)

	// 处理 else 语句的 stmt 部分
	if err := genStmt(p); err != nil {
		return err
	}

	// 标记 if-else 语句之后的代码位置
	p.EmitLabel(afterStmtLabel)

	return nil
}

// stmt -> while(bool) stmt
func genStmtWhile(p *Parser) error {
	// 生成 while 语句的代码
	startLabel := p.NewLabel() // 创建新的标签，用于 while 循环开始的位置
	condition := p.TokenStack[len(p.TokenStack)-3]
	afterStmtLabel := p.NewLabel() // 创建新的标签，用于 while 循环之后的代码位置

	// 标记循环开始的位置
	p.EmitLabel(startLabel)

	// 生成条件为假时跳转的代码
	p.Emit("ifFalse", string(condition), "goto", afterStmtLabel)

	// 处理 while 语句的 stmt 部分
	if err := genStmt(p); err != nil {
		return err
	}

	// 循环结束后跳回循环开始
	p.Emit("goto", startLabel)

	// 标记循环之后的代码位置
	p.EmitLabel(afterStmtLabel)
	return nil
}

// stmt -> do stmt while(bool);
func genStmtDoWhile(p *Parser) error {
	// 生成 do-while 语句的代码
	startLabel := p.NewLabel() // 创建新的标签，用于 do-while 循环开始的位置
	condition := p.TokenStack[len(p.TokenStack)-2]

	// 标记循环开始的位置
	p.EmitLabel(startLabel)

	// 处理 do-while 语句的 stmt 部分
	if err := genStmt(p); err != nil {
		return err
	}

	// 生成条件为真时重复循环的代码
	p.Emit("if", string(condition), "goto", startLabel)
	return nil
}

// stmt -> break;
func genStmtBreak(p *Parser) error {
	// 生成 break 语句的代码
	breakLabel := p.GetBreakLabel() // 获取跳出循环或 switch 的标签
	p.Emit("goto", breakLabel)
	return nil
}

// stmt -> block
func genStmtBlock(p *Parser) error {
	// p.SymbolTable.EnterScope()
	fmt.Println("[符号表] 进入新的作用域")
	return nil
}

// loc -> id
func genLoc(p *Parser) error {
	// 获取变量名
	varName := p.TokenStack[len(p.TokenStack)-1]
	fmt.Printf("[符号表] 触发变量 %s 赋值\n", varName)
	p.LastName = string(varName)
	return nil
}

// loc -> loc_array
func genLocArray(p *Parser) error { return nil }

// loc_array -> loc[num]
func genLocArrayFinal(p *Parser) error {
	// 获取数组名
	arrayName := p.TokenStack[len(p.TokenStack)-3]
	// 获取数组索引
	arrayIndex := p.TokenStack[len(p.TokenStack)-2]
	fmt.Printf("[符号表] 触发数组变量 %s 赋值，索引为 %s\n", arrayName, arrayIndex)
	// 生成数组变量赋值代码
	p.Emit(string(arrayName), string(arrayIndex), "") // 数组变量赋值操作
	return nil
}

// bool -> join || bool
func genBoolOr(p *Parser) error {
	// 生成逻辑或操作的中间代码
	op1 := p.TokenStack[len(p.TokenStack)-3]
	op2 := p.TokenStack[len(p.TokenStack)-1]
	result := p.SymbolTable.NewTempAddr()          // 生成一个新的临时变量
	p.Emit(result, string(op1), string(op2), "||") // 逻辑或操作
	return nil
}

// bool -> join
func genBool(p *Parser) error {
	// 不需要生成中间代码，直接返回
	return nil
}

// join -> equality
func genJoin(p *Parser) error {
	// 这里不需要生成中间代码，因为相等性比较 equality 会单独处理
	return nil
}

// join -> join && equality
func genJoinAnd(p *Parser) error {
	// 生成逻辑与操作的中间代码
	op1 := p.TokenStack[len(p.TokenStack)-3]
	op2 := p.TokenStack[len(p.TokenStack)-1]
	result := p.SymbolTable.NewTempAddr()          // 生成一个新的临时变量
	p.Emit(result, string(op1), string(op2), "&&") // 逻辑与操作
	return nil
}

// equality -> equality == rel
func genEqualityEqual(p *Parser) error {
	// 生成等于操作的中间代码
	op1 := p.TokenStack[len(p.TokenStack)-3]
	op2 := p.TokenStack[len(p.TokenStack)-1]
	result := p.SymbolTable.NewTempAddr()          // 生成一个新的临时变量
	p.Emit(result, string(op1), string(op2), "==") // 等于操作
	return nil
}

// equality -> equality != rel
func genEqualityNotEqual(p *Parser) error {
	// 生成不等比较操作的中间代码
	op1 := p.TokenStack[len(p.TokenStack)-3]
	op2 := p.TokenStack[len(p.TokenStack)-1]
	result := p.SymbolTable.NewTempAddr()          // 生成一个新的临时变量
	p.Emit(result, "!=", string(op1), string(op2)) // 不等比较操作
	return nil
}

// equality -> rel
func genEquality(p *Parser) error {
	// 这里不需要生成中间代码，因为相对表达式 rel 会单独处理
	return nil
}

// rel -> expr < expr
func genRelLess(p *Parser) error {
	// 生成小于比较操作的中间代码
	op1 := p.TokenStack[len(p.TokenStack)-3]
	op2 := p.TokenStack[len(p.TokenStack)-1]
	result := p.SymbolTable.NewTempAddr()         // 生成一个新的临时变量
	p.Emit(result, "<", string(op1), string(op2)) // 小于比较操作
	return nil
}

// rel -> expr <= expr
func genRelLessEqual(p *Parser) error {
	// 生成小于等于比较操作的中间代码
	op1 := p.TokenStack[len(p.TokenStack)-3]
	op2 := p.TokenStack[len(p.TokenStack)-1]
	result := p.SymbolTable.NewTempAddr()          // 生成一个新的临时变量
	p.Emit(result, "<=", string(op1), string(op2)) // 小于等于比较操作
	return nil
}

// rel -> expr >= expr
func genRelGreaterEqual(p *Parser) error {
	// 生成大于等于比较操作的中间代码
	op1 := p.TokenStack[len(p.TokenStack)-3]
	op2 := p.TokenStack[len(p.TokenStack)-1]
	result := p.SymbolTable.NewTempAddr()
	p.Emit(result, ">=", string(op1), string(op2))
	return nil
}

// rel -> expr > expr
func genRelGreater(p *Parser) error {
	// 生成大于比较操作的中间代码
	op1 := p.TokenStack[len(p.TokenStack)-3]
	op2 := p.TokenStack[len(p.TokenStack)-1]
	result := p.SymbolTable.NewTempAddr()
	p.Emit(result, ">", string(op1), string(op2))
	return nil
}

// rel -> expr
func genRel(p *Parser) error {
	// 这里不需要生成中间代码，因为 expr 自身已经生成了相应的中间代码。
	return nil
}

// expr -> expr + term
func genExprAdd(p *Parser) error {
	// 生成加法操作的中间代码
	op1 := p.TokenStack[len(p.TokenStack)-3]
	op2 := p.TokenStack[len(p.TokenStack)-1]
	result := p.SymbolTable.NewTempAddr()
	p.Emit(result, "+", string(op1), string(op2))
	return nil
}

// expr -> expr - term
func genExprSub(p *Parser) error {
	op1 := p.TokenStack[len(p.TokenStack)-3]
	op2 := p.TokenStack[len(p.TokenStack)-1]
	result := p.SymbolTable.NewTempAddr()
	p.Emit(result, "-", string(op1), string(op2))
	return nil
}

// expr -> term
func genExpr(p *Parser) error {
	// 这里不需要生成中间代码，因为 term 将单独处理
	return nil
}

// term -> term * unary
func genTermMul(p *Parser) error {
	// 生成除法操作的中间代码
	op1 := p.TokenStack[len(p.TokenStack)-3]
	op2 := p.TokenStack[len(p.TokenStack)-1]
	result := p.SymbolTable.NewTempAddr()
	p.Emit(result, "*", string(op1), string(op2))
	return nil
}

// term -> term / unary
func genTermDiv(p *Parser) error {
	// 生成除法操作的中间代码
	op1 := p.TokenStack[len(p.TokenStack)-3]
	op2 := p.TokenStack[len(p.TokenStack)-1]
	result := p.SymbolTable.NewTempAddr()
	p.Emit(result, "/", string(op1), string(op2))
	return nil
}

// term -> unary
func genTerm(p *Parser) error {
	// 这里不需要生成中间代码，因为一元操作 unary 将单独处理
	return nil
}

// unary -> !unary
func genUnaryNot(p *Parser) error {
	// 生成逻辑非操作的中间代码
	op := p.TokenStack[len(p.TokenStack)-1]
	result := p.SymbolTable.NewTempAddr()
	p.Emit(result, "!", string(op))
	return nil
}

// unary -> -unary
func genUnaryNeg(p *Parser) error {
	// 生成负号操作的中间代码
	op := p.TokenStack[len(p.TokenStack)-1]
	result := p.SymbolTable.NewTempAddr()
	p.Emit(result, "-", string(op))
	return nil
}

// unary -> factor
func genUnary(p *Parser) error {
	// 不需要生成中间代码，因为因子 factor 将单独处理
	return nil
}

// factor -> (bool)
func genFactorBool(p *Parser) error {
	// 不需要生成中间代码，因为布尔表达式 bool 将单独处理
	return nil
}

// factor -> loc
func genFactorLoc(p *Parser) error {
	// 不需要生成中间代码，因为位置 loc 将单独处理
	return nil
}

// factor -> num
func genFactorNum(p *Parser) error {
	// 不需要生成中间代码，因为数字 num 将单独处理
	return nil
}

// factor -> real
func genFactorReal(p *Parser) error {
	// 不需要生成中间代码，因为实数 real 将单独处理
	return nil
}

// factor -> true
func genFactorTrue(p *Parser) error {
	// 不需要生成中间代码，因为常量 true 将单独处理
	return nil
}

// factor -> false
func genFactorFalse(p *Parser) error {
	// 不需要生成中间代码，因为常量 false 将单独处理
	return nil
}
