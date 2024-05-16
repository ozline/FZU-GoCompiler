// rules.go
// 文法规则的定义

package parser

import (
	"fmt"

	"github.com/ozline/CoursePractice-GoCompiler/consts"
)

// NewGrammar 初始化一个文法
func NewGrammar(rules []Production, terminals []consts.Terminal) *Grammar {
	return &Grammar{
		Productions: rules,
		Terminals:   terminals,
	}
}

// checkLeftRecursion 检查文法是否存在左递归
func (g *Grammar) CheckLeftRecursion() bool {
	for _, prod := range g.Productions {
		// 检查每个产生式 prod 的体 prod.Body 是否非空，并且头部 prod.Head 是否等于体的第一个符号 prod.Body[0]
		if len(prod.Body) > 0 && prod.Head == prod.Body[0] {
			fmt.Printf("发现左递归: %v -> %v\n", prod.Head, prod.Body)
			return true
		}
	}
	fmt.Println("没有发现左递归")
	return false
}

// isTerminal 判断符号是否为终结符
func (g *Grammar) IsTerminal(symbol consts.Symbol) bool {
	for _, t := range g.Terminals {
		if symbol == consts.Symbol(t) {
			return true
		}
	}
	return false
}

// PrintFirstSet 打印First集合
func (p *Parser) PrintFirstSet() {
	for nonTerminal, first := range p.FirstSet {
		if len(first) == 0 {
			continue
		}
		fmt.Print("FIRST(" + nonTerminal + ") : [ ")
		var isFirst = true
		for terminal := range first {
			if !isFirst {
				fmt.Print(", ")
			} else {
				isFirst = false
			}

			if terminal == EPSILON {
				fmt.Print("ε")
				continue
			}
			fmt.Print(terminal)
		}
		fmt.Println(" ]")
	}
}

// PrintfollowSet 打印Follow集合
func (p *Parser) PrintFollowSet() {
	for nonTerminal, first := range p.FollowSet {
		fmt.Print("FOLLOW(" + nonTerminal + ") : [ ")
		var isFirst = true
		for terminal := range first {
			if !isFirst {
				fmt.Print(", ")
			} else {
				isFirst = false
			}

			if terminal == EPSILON {
				fmt.Print("ε")
				continue
			}
			fmt.Print(terminal)
		}
		fmt.Println(" ]")
	}
}

// InitFirstSet 初始化First集合
func (p *Parser) InitFirstSet() {
	firstSet := make(FirstSet)

	// 初始化所有终结符的First集合
	for _, t := range p.Grammar.Terminals {
		// 每个终结符的First集合只包含自身
		firstSet[consts.Symbol(t)] = map[consts.Terminal]bool{t: true}
	}

	// 初始化所有非终结符的First集合
	for _, p := range p.Grammar.Productions {
		if _, exists := firstSet[p.Head]; !exists {
			// 每个非终结符的First集合为空
			firstSet[p.Head] = make(map[consts.Terminal]bool, 0)
		}
	}

	// 迭代计算First集，直到不再有变化为止
	changed := true
	for changed {
		changed = false
		for _, p := range p.Grammar.Productions {
			headFirstSet := firstSet[p.Head]

			// 如果产生式体为空，就将空（EPSILON）添加到产生式头部 p.Head 的 FIRST 集合中，并标记 FIRST 集合发生了变化。
			if len(p.Body) == 0 {
				if !headFirstSet[EPSILON] {
					headFirstSet[EPSILON] = true
					changed = true
				}
			}

			for _, sym := range p.Body {
				// 特殊情况：如果符号是 EPSILON，直接添加到 head 的 FIRST 集合
				if sym == EPSILON {
					if !headFirstSet[EPSILON] {
						headFirstSet[EPSILON] = true
						changed = true
					}
					break
				}

				// 如果 sym 是一个非终结符，就将 sym 的 FIRST 集合中的所有终结符添加到产生式头部 p.Head 的 FIRST 集合中，并标记 FIRST 集合发生了变化。
				if symFirstSet, isNonTerminal := firstSet[sym]; isNonTerminal {
					for terminal := range symFirstSet {
						// 如果当前符号的First集不包含ε，则将其加入到当前符号的First集中
						if terminal != EPSILON && !headFirstSet[terminal] {
							headFirstSet[terminal] = true
							changed = true
						}
					}
					// 如果 sym 的 FIRST 集合不包含空（EPSILON），就停止遍历产生式体。
					if !symFirstSet[EPSILON] {
						break
					}

					// 如果我们到达了产生式体的末尾，并且所有的符号都包含 EPSILON
					// 那么我们需要将 EPSILON 添加到产生式头部的 FIRST 集合中
					if sym == p.Body[len(p.Body)-1] && symFirstSet[EPSILON] {
						headFirstSet[EPSILON] = true
						changed = true
					}
				} else {
					// 如果 sym 是一个终结符，就将 sym 添加到产生式头部 p.Head 的 FIRST 集合中，并标记 FIRST 集合发生了变化，然后停止遍历产生式体。
					if !headFirstSet[consts.Terminal(sym)] {
						headFirstSet[consts.Terminal(sym)] = true
						changed = true
					}
					break
				}
			}
		}
	}

	p.FirstSet = firstSet
}

// InitfollowSet 初始化Follow集合
// 将 $（表示输入结束）加入到开始符号的 Follow 集中。
// 如果存在一个产生式 A → αBβ，那么所有在 First(β) 中的终结符都在 Follow(B) 中。
// 如果存在一个产生式 A → αB 或一个产生式 A → αBβ 并且 ε 在 First(β) 中，那么所有在 Follow(A) 中的符号都在 Follow(B) 中。
func (p *Parser) InitFollowSet() error {
	followSet := make(FollowSet)
	// 确保First集已被初始化
	if p.FirstSet == nil {
		return fmt.Errorf("FirstSet is not initialized")
	}

	// 初始化Follow集
	for _, prod := range p.Grammar.Productions {
		if _, exists := followSet[prod.Head]; !exists {
			followSet[prod.Head] = make(map[consts.Terminal]bool)
		}
	}

	// 初始化增广产生式头部的Follow集
	if _, exists := followSet[ARGUMENTED_PRODUCTION.Head]; !exists {
		followSet[ARGUMENTED_PRODUCTION.Head] = make(map[consts.Terminal]bool)
	}
	// 将$加入到开始符号的Follow集中
	followSet[ARGUMENTED_PRODUCTION.Head][TERMINATE_SYMBOL] = true

	// 迭代直到没有变化为止
	changed := true
	for changed {
		changed = false
		for _, prod := range p.Grammar.Productions {
			for i := 0; i < len(prod.Body); i++ {
				sym := prod.Body[i]

				// 如果不是非终结符就跳过
				if p.Grammar.IsTerminal(sym) {
					continue
				}
				followSetSym := followSet[sym]
				if followSetSym == nil {
					followSetSym = make(map[consts.Terminal]bool)
				}

				// 如果是产生式的最后一个符号或者下一个符号的First集包含EPSILON
				// 就将产生式头部的Follow集加入到当前符号的Follow集中
				// 将First(β) - {ε} 加入Follow(B)
				if i+1 == len(prod.Body) || p.containsEpsilon(p.computeFirstSetOfSequence(prod.Body[i+1:])) {
					for terminal := range followSet[prod.Head] {
						if !followSetSym[terminal] {
							followSetSym[terminal] = true
							changed = true
						}
					}
				}

				// 将后续符号的First集（除了EPSILON）加入到当前符号的Follow集中
				if i < len(prod.Body)-1 {
					nextSym := prod.Body[i+1]
					nextFirstSet := p.computeFirstSetOfSymbol(nextSym)
					for terminal := range nextFirstSet {
						if terminal != EPSILON && !followSetSym[terminal] {
							followSetSym[terminal] = true
							changed = true
						}
					}
				}

				// 将修改后的followSetSym重新赋值给followSet[sym]
				followSet[sym] = followSetSym
			}
		}
	}

	p.FollowSet = followSet
	return nil
}

// computeFirstSetOfSymbol 计算单个符号的First集合
func (p *Parser) computeFirstSetOfSymbol(sym consts.Symbol) map[consts.Terminal]bool {
	if p.Grammar.IsTerminal(sym) {
		return map[consts.Terminal]bool{consts.Terminal(sym): true}
	}
	return p.FirstSet[sym]
}

// computeFirstSetOfSequence 计算符号串的First集合
func (p *Parser) computeFirstSetOfSequence(sequence []consts.Symbol) map[consts.Terminal]bool {
	firstSetSeq := make(map[consts.Terminal]bool)
	for _, sym := range sequence {
		firstSetSym := p.FirstSet[sym]
		for terminal := range firstSetSym {
			if terminal != EPSILON {
				firstSetSeq[terminal] = true
			}
		}
		if !firstSetSym[EPSILON] {
			break
		}
	}
	return firstSetSeq
}

// containsEpsilon 检查First集合是否包含EPSILON
func (p *Parser) containsEpsilon(firstSet map[consts.Terminal]bool) bool {
	return firstSet[EPSILON]
}
