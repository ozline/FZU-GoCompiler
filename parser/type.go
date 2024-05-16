package parser

import (
	"github.com/ozline/CoursePractice-GoCompiler/consts"
	"github.com/ozline/CoursePractice-GoCompiler/intercoder"
)

// Production 结构体表示一个产生式
type Production struct {
	Head    consts.Symbol       // 产生式的头部
	Body    []consts.Symbol     // 产生式的体部
	Handler func(*Parser) error // 产生式的处理函数

	// 产生式是文法的基本组成部分，它由两部分组成：头部（Head）和体部（Body）。头部是一个非终结符，体部是一个符号序列，每个符号可以是终结符或非终结符。
	// 例如，对于产生式 E -> E + T，E 是头部，E + T 是体部。
}

// Grammar 结构体表示一个文法
type Grammar struct {
	Productions []Production      // 产生式集合
	Terminals   []consts.Terminal // 终结符集合
}

// FirstSet 表示First集
type FirstSet map[consts.Symbol]map[consts.Terminal]bool

// FollowSet 表示Follow集
type FollowSet map[consts.Symbol]map[consts.Terminal]bool

// LR1Item 表示一个LR(1)项
type LR1Item struct {
	Production Production      // 产生式
	Position   int             // 点的位置
	Lookahead  consts.Terminal // 展望符

	// 在 LR(1) 项中，我们用一个点来表示产生式右侧已经识别的部分和尚未识别的部分的分界点。例如，如果 Position 的值是 2，那么表示产生式右侧的前两个符号已经被识别。

	// 在 LR(1) 项中，我们还需要一个展望符（Lookahead Symbol），用来表示在识别完产生式右侧的某个部分后，我们需要展望符的符号。这个符号是一个终结符，用来帮助我们在分析时做出正确的决策。
}

// LR1Items 表示LR(1)项的集合
type LR1Items []LR1Item

// State 表示一个LR(1)状态，包含多个LR(1)项
type State struct {
	Items LR1Items // LR(1) 项集
	Index int      // 状态编号
}

// StateCollection 表示所有状态的集合
type StateCollection []*State

// Parser 结构体
type Parser struct {
	Grammar         *Grammar               // 文法
	FirstSet        FirstSet               // First集
	FollowSet       FollowSet              // Follow 集（后续发现在 LR（1）中并不需要）
	StateCollection StateCollection        // 状态集合
	ActionTable     ActionTable            // Action表，Action 表用来表示状态在某个输入符号下的动作，它是一个二维表，其中每个单元格包含了一个动作类型和一个状态编号。
	GotoTable       GotoTable              // Goto表，Goto 表用来表示状态之间的转移关系，它是一个二维表，其中每个单元格包含了一个状态编号，表示在某个状态下通过某个符号转移到另一个状态。
	SymbolTable     intercoder.SymbolTable // 符号表
	TokenStack      []consts.Symbol        // 符号栈
	StateStack      []int                  // 状态栈
	LastType        consts.Terminal        // 上一个终结符的类型
	LastSize        int                    // 上一个终结符的大小
	LastName        string                 // 上一个终结符的名字
	ThreeAddress    []string               // 三地址码

	// Parser 结构体是整个文法分析器的核心，它包含了文法、First集和状态集合等重要信息。
	// 在文法分析器中，我们需要用到文法的产生式集合、终结符集合、First集合、Follow集合等信息，这些信息都会被封装在 Parser 结构体中。
	// Parser 结构体还包含了状态集合，每个状态包含了多个 LR(1) 项，用来表示文法分析器在某个状态下的所有可能的分析情况。
}

type ActionEntry struct {
	ActionType string // "shift", "reduce", "accept", "error"
	Number     int    // 状态编号或产生式编号
}

type ActionTable map[int]map[consts.Terminal]ActionEntry // Action表，int表示状态编号，Terminal表示终结符
type GotoTable map[int]map[consts.Symbol]int             // Goto表，int表示状态编号，Symbol表示非终结符
