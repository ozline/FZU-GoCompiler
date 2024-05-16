# Parser - 文法分析器

## 文法规则

该编译原理实验要求文法规则如下

```text
program → block
block → { decls stmts }
decls → decls decl | ε
decl → type id;
type → type[num] | basic
stmts → stmts stmt | ε
stmt → loc=bool;
        | if(bool) stmt
        | if(bool) stmt else stmt
        | while(bool) stmt
        | do stmt while(bool);
        | break;
        | block
loc → loc[num] | id
bool → bool || join | join
join → join ＆& equality | equality
equality → equality == rel | equality != rel | rel
rel → expr<expr | expr<=expr | expr>=expr | expr>expr | expr
expr → expr+term | expr-term | term
term → term*unary | term/unary | unary
unary → !unary | -unary | factor
factor → (bool) | loc | num | real | true | false
```

分析可知
```text
program → block：程序由一个块构成。
block → { decls stmts }：一个块由声明（decls）和语句（stmts）组成，被大括号包围。
decls → decls decl | ε：声明可以是多个声明后跟一个声明，或者为空（ε表示空）。
decl → type id;：一个声明由一个类型后跟一个标识符组成，以分号结束。
type → type[num] | basic：类型可以是数组类型或基本类型。
stmts → stmts stmt | ε：语句可以是多个语句后跟一个语句，或者为空。
stmt → ...：定义了多种语句的形式，包括赋值、if语句、while循环、do-while循环、break语句和块。
loc → loc[num] | id：位置可以是数组元素或标识符。
bool → bool || join | join：布尔表达式可以是由逻辑或连接的布尔表达式，或者更简单的形式。
join → join && equality | equality：更加具体的布尔表达式，通过逻辑与连接。
equality → ...：等式表达式，可以是相等或不等的比较。
rel → ...：关系表达式，包括小于、小于等于、大于等于、大于的比较。
expr → ...：表达式可以是加法、减法或更简单的项。
term → ...：项可以是乘法、除法或更简单的单元。
unary → ...：单元可以是逻辑非、负号或更简单的因子。
factor → ...：因子可以是括号内的布尔表达式、位置、数字、真实数值、true或false。
```

## 实现的具体步骤

要实现一个文法分析器，步骤如下

### 1. initGrammar 初始化文法

定义文法的产生式，区分终结符和非终结符，并初始化产生式集合。

### 2. initTerminal 初始化终结符

定义文法中的终结符，即输入流中可以出现的符号

### 3. initFirstSet 初始化 First 集

计算每个非终结符的First集，即可以从该非终结符推导出的起始符号集合。

### 4. buildItems 构建状态集

构建项集族，即状态集，每个状态都包含一组LR(1)项，表示在分析过程中的不同点。

### 5. constructTable 构建 LR(1) 分析表

根据状态集构建LR(1)分析表，包括ACTION表和GOTO表，这些表用于在语法分析过程中指导分析器的行为。

### 6. parse 进行 LR(1) 分析

使用构建好的LR(1)分析表对输入的程序进行分析，通过维护状态栈和符号栈来进行移进、归约和接受操作。

# 实验时遇到的困难

基本都来自于 go 语言本身

## go1.22中的for loop问题

当我们使用
```go
for _, v := range list {
    go func() {
        fmt.Println(v)
    }
}
```
这样的式子时，当 list 一直被 append，v 的值会一直是同一个，这个在 go1.22（上个月）的新版本中才被修复

## 指针比较问题

我在`item.go`中实现一个函数如下所示
```go
// contains 检查LR(1)项是否已经存在于集合中
// 由于LR(1)项集合是一个集合，所以我们需要检查项是否已经存在，以避免重复添加相同的项。
func contains(items LR1Items, item LR1Item) bool {
	for _, i := range items {
		if i.Production == item.Production && i.Position == item.Position && i.Lookahead == item.Lookahead {
			return true
		}
	}
	return false
}
```
这本来是正确的逻辑，但是问题出在 Production 这是一个指针`*Production`，比较时实际上是比较两个指针本身，会导致比较一直失败（因为指针地址不可能一样），即使内容本身是一样的

这个困扰了我很久的时间

在我修复这个 bug 后，我并没有及时的去看其余代码是否可以修正，因此我又出现了状态无限递增的问题

结果发现在 equalState() 函数中也出现了这个校验，导致状态无限递增下去了