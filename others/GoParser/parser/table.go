// table.go
// LR(1)分析表的构建

package parser

import "fmt"

func (p *Parser) PrintGoToTable() {
	fmt.Println("GOTO 表")
	for i, gotoMap := range p.GotoTable {
		for sym, stateIndex := range gotoMap {
			fmt.Printf("GOTO[%d, %s] = %d\n", i, sym, stateIndex)
		}
	}
}

func (p *Parser) PrintActionTable() {
	fmt.Println("ACTION 表")
	for i, actionMap := range p.ActionTable {
		for sym, actionEntry := range actionMap {
			fmt.Printf("ACTION[%d, %s] = %v\n", i, sym, actionEntry)
		}
	}
}

func (p *Parser) buildGotoTable() {
	// 遍历所有的 LR(1) 状态 state
	for i, state := range p.StateCollection {
		// 对于每个状态 state，遍历其中的所有项 item
		for _, item := range state.Items {
			// 如果项 item 的点位置小于产生式体的长度，说明还有未处理的符号。

			if item.Position < len(item.Production.Body) {
				sym := item.Production.Body[item.Position]

				// 如果 sym 是一个非终结符，计算 sym 的转移后的项集 nextState
				if !p.Grammar.IsTerminal(sym) {
					nextState := p.gotoState(state.Items, sym)

					// 如果 nextState 在状态集合 p.StateCollection 中，将动作表 p.GotoTable[i] 的对应项设置为 nextState 在 p.StateCollection 中的索引。
					if nextStateIndex, exists := p.containsState(p.StateCollection, &State{Items: nextState}); exists {
						if p.GotoTable[i] == nil {
							p.GotoTable[i] = make(map[Symbol]int)
						}

						if existingState, exists := p.GotoTable[i][sym]; exists {
							// 冲突检测：已有状态与当前状态不同
							if existingState != nextStateIndex {
								fmt.Printf("设置 Goto 表发生冲突! 状态: %d 符号: '%s'\n", i, sym)
							}
						} else {
							p.GotoTable[i][sym] = nextStateIndex
						}
					}
				}
			}
		}
	}
}

func (p *Parser) buildActionTable() {
	// 遍历所有的 LR(1) 状态 state
	for i, state := range p.StateCollection {
		// 对于每个状态 state，遍历其中的所有项 item
		for _, item := range state.Items {

			// 如果项 item 的点位置等于产生式体的长度，说明没有未处理的符号，需要执行规约动作或接受动作。
			// 这个 Position 指的是下一个要处理的符号的位置，所以如果等于产生式体的长度，说明没有未处理的符号。

			// 当没有未处理的符号时，执行规约动作或接受动作。
			// 特殊处理，即 item 只有一个 EPSILON
			if item.Position == len(item.Production.Body) || (item.Position == len(item.Production.Body)-1 && item.Production.Body[item.Position] == EPSILON) {
				if item.Production.Head == ARGUMENTED_PRODUCTION.Head && item.Lookahead == TERMINATE_SYMBOL {
					// 接受动作

					// 确保 p.ActionTable[i] 已经初始化，然后再进行赋值
					if p.ActionTable[i] == nil {
						p.ActionTable[i] = make(map[Terminal]ActionEntry)
					}
					p.ActionTable[i][item.Lookahead] = ActionEntry{
						ActionType: ACCEPT,
						Number:     0,
					}
					continue
				}
				// 执行规约动作，动作的参数是产生式在文法的产生式列表中的索引。

				// 由于 LR（1）只有1 个展望符，这边不需要遍历展望符
				// for lookahead := range p.computeLookahead(item.Production.Body[:item.Position], item.Lookahead) {
				if p.ActionTable[i] == nil {
					p.ActionTable[i] = make(map[Terminal]ActionEntry)
				}
				// 找到产生式编号
				prodIndex := -1 // 表示无效索引值，用于检测是否找到了对应的产生式
				for index, prod := range p.Grammar.Productions {
					if equalProductions(prod, item.Production) {
						prodIndex = index
						break
					}
				}

				if existingAction, exists := p.ActionTable[i][item.Lookahead]; exists {
					// 冲突检测：已有动作与当前动作不同
					if existingAction.ActionType != REDUCE || existingAction.Number != prodIndex {
						fmt.Printf("设置 REDUCE 发生冲突! 状态: %d 展望符: '%s'\n", i, item.Lookahead)
					}
				} else {
					p.ActionTable[i][item.Lookahead] = ActionEntry{
						ActionType: REDUCE,
						Number:     prodIndex,
					}
				}
				// }
			} else {
				// 如果项 item 的点位置小于产生式体的长度，说明还有未处理的符号，需要执行移入动作。
				// 这里也有可能是接受动作，需要检查是否是增广产生式。

				// 如果产生式是增广产生式，且点的位置等于产生式体的长度，且展望符是终止符，执行接受动作。
				sym := item.Production.Body[item.Position]

				// 如果 sym 是一个终结符，计算 sym 的转移后的项集 nextState

				if p.Grammar.IsTerminal(sym) {
					nextState := p.gotoState(state.Items, sym)

					if i == 25 {
						fmt.Println("")
					}

					// 如果 nextState 在状态集合 p.StateCollection 中，将动作表 p.ActionTable[i] 的对应项设置为移入动作，动作的参数是 nextState 在 p.StateCollection 中的索引。
					if nextStateIndex, exists := p.containsState(p.StateCollection, &State{Items: nextState}); exists {
						if p.ActionTable[i] == nil {
							p.ActionTable[i] = make(map[Terminal]ActionEntry)
						}

						if existingAction, exists := p.ActionTable[i][Terminal(sym)]; exists {
							// 冲突检测：已有动作与当前动作不同
							if existingAction.ActionType != SHIFT || existingAction.Number != nextStateIndex {
								fmt.Printf("设置 SHIFT 发生冲突! 状态: %d 符号: '%s', 期望数字：%+v, 现有内容为：%+v\n", i, sym, nextStateIndex, existingAction)
							}
						} else {
							p.ActionTable[i][Terminal(sym)] = ActionEntry{
								ActionType: SHIFT,
								Number:     nextStateIndex,
							}
						}
					}
				}
			}
		}
	}
}
