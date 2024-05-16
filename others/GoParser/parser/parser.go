// parser.go
// 文法分析器的实现

package parser

import (
	"fmt"

	"github.com/ozline/CoursePractice-GoCompiler/lexer"
)

// NewParser 创建一个新的 Parser 实例
func NewParser() *Parser {
	grammar := NewGrammar(PRODUCTIONS, TERMINALS)

	parser := &Parser{
		Grammar:         grammar,
		StateCollection: []*State{},
		ActionTable:     make(ActionTable),
		GotoTable:       make(GotoTable),
	}
	return parser
}

func (p *Parser) BuildTables() {
	p.buildGotoTable()
	p.buildActionTable()
}

func TokenToTerminal(token lexer.Token) Terminal {
	// fmt.Println("Token:(", token.Value, ") - (", token.Type, ")", lexer.TokenTypes[token.Type])
	if token.Type == lexer.NUMBER {
		return Terminal("num")
	}
	if token.Type == lexer.REAL {
		return Terminal("real")
	}
	if token.Type == lexer.IDENTIFIER {
		return Terminal("id")
	}
	if token.Type == lexer.TYPE { // 类型
		return Terminal("basic")
	}
	if token.Type == lexer.EOF {
		return Terminal(TERMINATE_SYMBOL)
	}
	return Terminal(token.Value)
}

func TokenToSymbol(token lexer.Token) Symbol {
	return Symbol(TokenToTerminal(token))
}

func (p *Parser) Parse(l *lexer.Lexer) error {
	// 初始化分析栈，初始状态为 0
	stateStack := []int{0}      // 状态栈
	tokenStack := []Symbol{"$"} // 符号栈`$`表示输入结束
	cnt := int(0)

	token, err := l.NextToken()
	if err != nil {
		return err
	}

	// 主循环，直到接受或遇到错误
	for {
		fmt.Printf("\n\n=====================================\n")
		// 打印状态和符号栈
		cnt++
		fmt.Printf("第 %d 步\n", cnt)
		fmt.Printf("状态栈: %v\n", stateStack)
		fmt.Printf("符号栈: %v\n", tokenStack)
		// 查看栈顶状态
		state := stateStack[len(stateStack)-1]

		// 将 Token 转换为终结符
		terminal := TokenToTerminal(token)

		// 根据当前状态和读取的 Token（终结符）查找 Action 表中的动作
		fmt.Printf("当前状态: %d, 当前符号: %s 转换后: %s\n", state, token.Value, terminal)

		action, ok := p.ActionTable[state][terminal]
		if !ok {
			// 如果没有找到动作，打印错误消息并退出
			return fmt.Errorf("解析错误：无法找到状态 %d 和符号 %s 的动作\n", state, terminal)
		}

		fmt.Printf("动作类别: %s 期望下一步状态: %d\n", action.ActionType, action.Number)

		switch action.ActionType {
		case SHIFT:
			// 移入操作：将 Token 和新状态推入栈中
			stateStack = append(stateStack, action.Number)
			tokenStack = append(tokenStack, TokenToSymbol(token))
			fmt.Printf("执行移入操作\n")
			token, err = l.NextToken()
			if err != nil {
				return err
			}

		case REDUCE:
			// 规约操作：使用产生式规约，并将相应的符号数从栈中弹出
			production := p.Grammar.Productions[action.Number]
			fmt.Printf("使用产生式 %v -> %v 规约\n", production.Head, production.Body)

			if !(len(production.Body) == 1 && production.Body[0] == EPSILON) {
				for range production.Body {
					tokenStack = tokenStack[:len(tokenStack)-1]
					stateStack = stateStack[:len(stateStack)-1]
				}
			}

			var gotoState int
			tokenStack = append(tokenStack, production.Head)
			topState := stateStack[len(stateStack)-1]
			gotoState, ok = p.GotoTable[topState][production.Head]
			if !ok {
				return fmt.Errorf("解析错误：无法在状态 %v 中找到产生式 %v 的转移状态\n", topState, production)
			}

			// 将产生式头部和转移后的状态推入栈中
			fmt.Printf("转移状态到 %d\n", gotoState)
			stateStack = append(stateStack, gotoState)

		case ACCEPT:
			// 接受操作：成功完成分析
			fmt.Println("成功完成解析.")
			return nil

		case ERROR:
			// 错误操作：打印错误并退出
			return fmt.Errorf("解析错误: %s\n", action.ActionType)
		}

		// 打印栈的内容
		// fmt.Printf("分析栈: \n\n%v\n", stack)
	}
}
