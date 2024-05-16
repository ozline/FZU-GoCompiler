package parser

import "fmt"

// LabelCounter 用于生成新的唯一标签
var LabelCounter int = 0

// BreakLabelStack 用于维护 break 标签的栈
var BreakLabelStack []string = []string{}

// NewLabel 生成一个新的唯一标签
func (p *Parser) NewLabel() string {
	label := fmt.Sprintf("L%d", LabelCounter)
	LabelCounter++
	return label
}

// EmitLabel 将标签输出到三地址码中
func (p *Parser) EmitLabel(label string) {
	p.ThreeAddress = append(p.ThreeAddress, label+":")
}

// GetBreakLabel 返回最内层循环或 switch 语句的结束标签
func (p *Parser) GetBreakLabel() string {
	if len(BreakLabelStack) == 0 {
		panic("发生逻辑错误: BreakLabelStack 为空，但 GetBreakLabel 被调用。")
	}
	// 返回栈顶元素
	return BreakLabelStack[len(BreakLabelStack)-1]
}

// EnterLoop 用于在进入循环或 switch 时调用，将新的 break 标签压入栈
func (p *Parser) EnterLoop() {
	label := p.NewLabel()
	BreakLabelStack = append(BreakLabelStack, label)
}

// ExitLoop 用于在离开循环或 switch 时调用，将 break 标签弹出栈
func (p *Parser) ExitLoop() {
	if len(BreakLabelStack) == 0 {
		panic("发生逻辑错误: BreakLabelStack 为空，但 ExitLoop 被调用。")
	}
	BreakLabelStack = BreakLabelStack[:len(BreakLabelStack)-1]
}

// Emit 输出三地址代码
// result: 结果 opcode: 操作符 operands: 操作数
func (p *Parser) Emit(result string, opcode string, operands ...string) {
	var op string
	// 输出三地址代码
	if opcode == "" {
		op = fmt.Sprintf("%s = %s\n", result, operands[0])
	} else {
		op = fmt.Sprintf("%s = %s %s %s;\n", result, operands[0], opcode, operands[1])
	}

	p.ThreeAddress = append(p.ThreeAddress, op)
}

func (p *Parser) PrintThreeAddress() {
	fmt.Println("\n\n===============三地址码===============")
	for i, code := range p.ThreeAddress {
		fmt.Printf("%d: %s", i, code)
	}
}
