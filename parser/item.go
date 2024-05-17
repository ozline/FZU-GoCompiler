// item.go
// LR(1)项的定义和状态集合的构建

package parser

import (
	"fmt"
	"slices"

	"github.com/ozline/CoursePractice-GoCompiler/consts"
)

// PrintStateCollection 打印状态集合
func (p *Parser) PrintStateCollection() {
	fmt.Println("状态集合 - 共有", len(p.StateCollection), "个状态")
	for i, state := range p.StateCollection {
		fmt.Println("状态", i)
		for _, item := range state.Items {
			fmt.Printf("%v -> %v 位置：%v 展望符：%v 长度: %v\n", item.Production.Head, item.Production.Body, item.Position, item.Lookahead, len(item.Production.Body))
		}
	}
}

// buildItems 构建LR(1)项集合
// 该函数是构建状态集合的基础，它会计算文法的闭包和转移，以构建状态集合。
func (p *Parser) BuildStateCollection() {
	// 初始化一个产生式，产生式的头部是文法的开始符号，体部是文法的第一个产生式体，点的位置为0，展望符为终止符
	startProd := ARGUMENTED_PRODUCTION
	// 产生式的点位置为0，展望符为终止符
	startItem := LR1Item{Production: startProd, Position: 0, Lookahead: TERMINATE_SYMBOL}

	initialState := &State{Items: p.closure(LR1Items{startItem}), Index: 0}

	states := StateCollection{initialState}
	toProcess := []*State{initialState}

	// 循环直到不再有新状态
	for len(toProcess) > 0 {
		state := toProcess[0]     // 获取当前状态
		toProcess = toProcess[1:] // 出队当前状态

		// 对每个符号，计算转移并尝试添加新状态
		for _, sym := range p.getAllSymbols() {
			gotoItems := p.gotoState(state.Items, sym)
			if len(gotoItems) > 0 {
				newState := &State{Items: p.closure(gotoItems), Index: len(states)}
				// 如果新状态不存在，就添加到状态集合中
				if _, exists := p.containsState(states, newState); !exists {
					states = append(states, newState)
					toProcess = append(toProcess, newState) // 入队新状态
				}
			}
		}
	}

	p.StateCollection = states
}

// getAllSymbols 返回文法中所有的符号（终结符和非终结符）
// 该函数用于计算状态转移时需要遍历的所有符号。
func (p *Parser) getAllSymbols() []consts.Symbol {
	var symbols []consts.Symbol
	for _, prod := range p.Grammar.Productions {
		symbols = append(symbols, prod.Head)
		for _, sym := range prod.Body {
			symbols = append(symbols, sym)
		}
	}
	// 去重并返回
	return uniqueSymbols(symbols)
}

// uniqueSymbols 返回唯一符号集合
// 该函数用于去重符号集合，因为我们需要在状态转移时遍历所有符号。
func uniqueSymbols(symbols []consts.Symbol) []consts.Symbol {
	seen := make(map[consts.Symbol]bool)
	var unique []consts.Symbol
	for _, sym := range symbols {
		// 如果符号不在 seen 中，就添加到 unique 中
		if _, ok := seen[sym]; !ok && sym != EPSILON {
			seen[sym] = true
			unique = append(unique, sym)
		}
	}
	return unique
}

// gotoState 计算一个状态通过一个符号的转移
// 它检查项集中每个项的下一个符号，如果与给定符号相匹配，就创建一个新的项，其中点位置向前移动一位。
/*
	例如，我们有一个如下项集{A -> B.C, a, D -> E.F, b}
	如果我们调用 gotoState(items, 'C')，那么函数会遍历这个项集，对于第一个项 A -> B.C, a，它的点位置是 1，小于产生式体 B.C 的长度 2，所以下一个未识别的符号是 C，与给定的符号 C 相匹配，所以这个项可以通过 C 转移到一个新的项 A -> BC., a。对于第二个则不会转移
	所以新的项集就是 {A -> BC., a}
*/
func (p *Parser) gotoState(items LR1Items, sym consts.Symbol) LR1Items {
	var gotoItems LR1Items = make(LR1Items, 0)
	for _, item := range items {
		// 对于每个项 item，检查它的点位置 item.Position 是否小于产生式体 item.Production.Body 的长度。如果是，说明点还没有移动到产生式体的末尾，也就是说还有未识别的符号。
		// 如果还有未识别的符号，检查下一个未识别的符号 item.Production.Body[item.Position] 是否等于给定的符号 sym。如果是，说明这个项可以通过 sym 转移到一个新的项。
		if item.Position < len(item.Production.Body) && item.Production.Body[item.Position] == sym {
			newItem := LR1Item{
				Production: item.Production,
				Position:   item.Position + 1,
				Lookahead:  item.Lookahead,
			}
			gotoItems = append(gotoItems, newItem)
		}
	}
	return p.closure(gotoItems)
}

// containsState 检查一个状态是否已经存在于状态集合中，并返回状态的索引和是否存在的布尔值。
// 如果状态存在，则返回状态的索引和true；如果不存在，则返回-1和false。
func (p *Parser) containsState(states StateCollection, state *State) (int, bool) {
	for _, s := range states {
		if len(s.Items) != len(state.Items) {
			continue
		}
		if p.equalStates(s, state) {
			return s.Index, true
		}
	}
	return -1, false
}

// equalStates 检查两个状态是否相等（包含相同的LR(1)项）
// 由于状态集合是一个集合，所以我们需要检查状态是否相等，以避免重复添加相同的状态。
func (p *Parser) equalStates(s1, s2 *State) bool {
	// 如果两个状态的项数不相等，就肯定不相等
	if len(s1.Items) != len(s2.Items) {
		return false
	}

	for _, item1 := range s1.Items {
		found := false
		// 对于 s1 的每一个 LR(1) 项 item1，在 s2 的 LR(1) 项中查找一个相同的项 item2。
		// 相同的定义是，item1 和 item2 的产生式、点位置和展望符都相等。
		for _, item2 := range s2.Items {
			if equalProductions(item1.Production, item2.Production) && item1.Position == item2.Position && item1.Lookahead == item2.Lookahead {
				found = true
				break
			}
		}

		if !found {
			// fmt.Printf("状态不相等, 项[%v, %v, %v] 没有找到\n", item1.Production.Head, item1.Production.Body, item1.Lookahead)
			return false
		}
	}
	return true
}

func itemKey(item LR1Item) string {
	return fmt.Sprintf("%v|%v|%v|%v", item.Production.Head, item.Production.Body, item.Lookahead, item.Position)
}

// closure 计算LR(1)项集的闭包
// 闭包是一个重要的概念，用来计算一个状态的所有可能项。在构建状态集合时，我们需要计算每个状态的闭包，以便在状态转移时能够正确地处理展望符。
// items 是一个 LR(1) 项集，expanded 是一个映射，用来记录哪些项已经扩展过了。
func (p *Parser) closure(items LR1Items) LR1Items {
	closure := make(LR1Items, len(items))
	copy(closure, items)

	expanded := make(map[string]bool, 0)

	changed := true // 标记是否闭包发生了变化，如果没有变化，就不需要继续扩展了
	for changed {
		changed = false

		for _, item := range closure {
			if expanded[itemKey(item)] {
				// 如果我们已经扩展过这个项，就跳过它
				// fmt.Printf("项已经扩展过: [%v, %v, %v]\n", item.Production.Head, item.Production.Body, item.Lookahead)
				continue
			}

			// fmt.Printf("正在检查以下项的闭包：[%v, %v, %v]\n", item.Production.Head, item.Production.Body, item.Lookahead)

			expanded[itemKey(item)] = true // 标记当前项为已扩展

			// 如果点的位置在产生式的末尾，就不需要扩展了
			if item.Position >= len(item.Production.Body) {
				continue
			}

			// 取出点位置的下一个符号nextSym。如果nextSym是一个非终结符，则需要扩展这个项。
			nextSym := item.Production.Body[item.Position]
			if !p.Grammar.IsTerminal(nextSym) {
				// 计算item的展望符集lookaheads
				lookaheads := p.computeLookahead(item.Production.Body[item.Position+1:], item.Lookahead)

				// 遍历语法的所有产生式prod。如果prod的头部是nextSym，则需要创建一个新的项newItem，并将其添加到closure中。
				// newItem的产生式是prod，点位置是0，展望符是lookaheads中的每一个符号。
				for _, prod := range p.Grammar.Productions {
					// （拓展）如果产生式的头部是下一个符号，就添加新的项到闭包
					if prod.Head == nextSym {
						if len(prod.Body) == 0 {
							newItem := LR1Item{
								Production: prod,
								Position:   0,
								Lookahead:  item.Lookahead,
							}
							if !contains(closure, newItem) {
								closure = append(closure, newItem)
								changed = true
							}
						} else {
							for lookahead := range lookaheads {
								newItem := LR1Item{
									Production: prod,
									Position:   0,
									Lookahead:  lookahead,
								}

								if !contains(closure, newItem) {
									closure = append(closure, newItem)
									changed = true
								}
							}
						}
					}
				}
			}
		}
	}
	return closure
}

// contains 检查LR(1)项是否已经存在于集合中
// 由于LR(1)项集合是一个集合，所以我们需要检查项是否已经存在，以避免重复添加相同的项。
func contains(items LR1Items, item LR1Item) bool {
	for _, i := range items {
		if equalLR1Items(i, item) {
			return true
		}
	}
	return false
}

// equalLR1Items 检查两个LR(1)项是否相等
func equalLR1Items(i1, i2 LR1Item) bool {
	return equalProductions(i1.Production, i2.Production) && i1.Position == i2.Position && i1.Lookahead == i2.Lookahead
}

// equalProductions 检查两个产生式是否相等
func equalProductions(p1, p2 Production) bool {
	if p1.Head != p2.Head {
		return false
	}
	if !slices.Equal(p1.Body, p2.Body) {
		return false
	}
	return true
}

// computeLookahead 计算展望符的First集合
// 展望符是在识别产生式右侧的某个部分后，我们需要展望符的符号。这个符号是一个终结符，用来帮助我们在分析时做出正确的决策。
func (p *Parser) computeLookahead(symbols []consts.Symbol, lookahead consts.Terminal) map[consts.Terminal]bool {
	// 如果后面没有符号了，展望符就是当前项的展望符
	if len(symbols) == 0 {
		return map[consts.Terminal]bool{lookahead: true}
	}

	allNullable := true
	firstSet := make(map[consts.Terminal]bool)
	for _, sym := range symbols {
		if p.Grammar.IsTerminal(sym) {
			// 如果是终结符，就将这个终结符添加到firstSet中
			firstSet[consts.Terminal(sym)] = true
		}

		// 将 sym 的 FIRST 集合中的每一个终结符添加到 firstSet 中，但是忽略空（EPSILON）。
		// 我们只关心可以立即开始的符号，而不是可以是空的符号
		for terminal := range p.FirstSet[sym] {
			if terminal != EPSILON { // 忽略EPSILON
				firstSet[terminal] = true
			}
		}

		// 如果当前符号的First集合不包含EPSILON，那么它不是可空的，停止
		// 因为我们已经找到了一个可以产生终结符的符号
		/*
			在计算 FIRST 集合时，我们只关心可以立即开始的符号。在产生式右侧的符号序列中，如果一个符号的 FIRST 集合包含 EPSILON，那么就意味着这个符号可以为空，我们需要继续查看下一个符号。如果一个符号的 FIRST 集合不包含 EPSILON，那么这个符号就不能为空，我们就可以停止查看后续的符号。
		*/
		if !firstSet[EPSILON] {
			allNullable = false
			break
		}
	}

	// 如果所有符号都是可空的，或者没有符号，添加原始的展望符
	// 例如，假设有 A -> B C D，其中 B、C、D 都可以推导出空串，那么在识别 A 时，我们可能会立即看到 B、C 和 D 后面的符号，也就是原始的展望符
	if allNullable {
		firstSet[lookahead] = true
	}

	return firstSet
}
