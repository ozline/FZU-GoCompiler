// parser.go
// 文法分析器的实现

package parser

import (
	"fmt"

	"github.com/ozline/CoursePractice-GoCompiler/consts"
	"github.com/ozline/CoursePractice-GoCompiler/intercoder"
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
		SymbolTable:     *intercoder.NewSymbolTable(),
	}
	return parser
}

func (p *Parser) BuildTables() {
	p.buildGotoTable()
	p.buildActionTable()
}

func TokenToTerminal(token lexer.Token) consts.Terminal {
	// fmt.Println("Token:(", token.Value, ") - (", token.Type, ")", lexer.TokenTypes[token.Type])
	if token.Type == lexer.NUMBER {
		return consts.Terminal("num")
	}
	if token.Type == lexer.REAL {
		return consts.Terminal("real")
	}
	if token.Type == lexer.IDENTIFIER {
		return consts.Terminal("id")
	}
	if token.Type == lexer.TYPE { // 类型
		return consts.Terminal("basic")
	}
	if token.Type == lexer.EOF {
		return consts.Terminal(TERMINATE_SYMBOL)
	}
	return consts.Terminal(token.Value)
}

func TokenToSymbol(token lexer.Token) consts.Symbol {
	return consts.Symbol(TokenToTerminal(token))
}

func (p *Parser) Parse(l *lexer.Lexer) error {
	var token lexer.Token
	var err error
	// 初始化分析栈，初始状态为 0
	p.StateStack = []int{0}                                         // 状态栈
	p.TokenStack = []consts.Symbol{consts.Symbol(TERMINATE_SYMBOL)} // 预留一个空位，用于处理状态 0 的转移
	cnt := int(0)
	p.SymbolTable.EnterScope() // 进入一个新的作用域

	// 主循环，直到接受或遇到错误
	readNextToken := true
	fmt.Printf("\n\n===============开始解析===============")
	for {
		fmt.Printf("\n\n=====================================\n")
		// 打印状态和符号栈
		cnt++
		fmt.Printf("第 %d 步\n", cnt)
		fmt.Printf("状态栈: %v\n", p.StateStack)
		fmt.Printf("符号栈: %v\n", p.TokenStack)
		// 查看栈顶状态
		state := p.StateStack[len(p.StateStack)-1]

		if readNextToken {
			token, err = l.NextToken()
			if err != nil {
				return err
			}
			readNextToken = false
		}

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
			p.StateStack = append(p.StateStack, action.Number)
			p.TokenStack = append(p.TokenStack, consts.Symbol(token.Value))
			fmt.Printf("执行移入操作\n")
			readNextToken = true
			break
		case REDUCE:
			// 规约操作：使用产生式规约，并将相应的符号数从栈中弹出
			production := p.Grammar.Productions[action.Number]
			fmt.Printf("使用产生式 %v -> %v 规约\n", production.Head, production.Body)

			// 处理规约操作
			if err := production.Handler(p); err != nil {
				return err
			}

			// 当我们执行归约操作时，需要将产生式右侧的符号从栈中弹出，并将产生式左侧的符号推入栈中
			if !(len(production.Body) == 1 && production.Body[0] == EPSILON) {
				// fmt.Printf("开始弹出%d 个: ", len(production.Body))
				for range production.Body {
					// fmt.Printf("%v ", p.TokenStack[len(p.TokenStack)-1])
					p.TokenStack = p.TokenStack[:len(p.TokenStack)-1]
					p.StateStack = p.StateStack[:len(p.StateStack)-1]
				}
				// fmt.Println()
			}

			var gotoState int
			// 将产生式头部推入栈中
			p.TokenStack = append(p.TokenStack, production.Head)
			topState := p.StateStack[len(p.StateStack)-1]
			gotoState, ok = p.GotoTable[topState][production.Head]
			if !ok {
				return fmt.Errorf("解析错误：无法在状态 %v 中找到产生式 %v 的转移状态\n", topState, production)
			}

			// 将产生式头部和转移后的状态推入栈中
			fmt.Printf("转移状态到 %d\n", gotoState)
			p.StateStack = append(p.StateStack, gotoState)
			break
		case ACCEPT:
			// 接受操作：成功完成分析
			fmt.Println("\n\n>>> 成功完成解析.")
			p.SymbolTable.ExitScope() // 确保退出全局作用域
			return nil

		case ERROR:
			// 错误操作：打印错误并退出
			return fmt.Errorf("解析错误: %s\n", action.ActionType)
		}
	}
}
