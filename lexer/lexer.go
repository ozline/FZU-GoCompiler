// lexer.go
// 词法分析器的实现

package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// Lexer 词法分析器结构体
type Lexer struct {
	reader           *bufio.Reader
	line             int
	column           int
	lastReadRuneSize int   // 上次读取的字符大小
	lineLengths      []int // 用于存储每行的长度
}

// NewLexer 创建一个新的词法分析器实例
func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		reader: bufio.NewReader(reader),
		line:   1, // 第一行
		column: 0, // 第零列
	}
}

// isLetter 检查字符是否是字母
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}

// isDigit 检查字符是否是数字
func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

// skipWhitespace 跳过空白字符
func (l *Lexer) skipWhitespace() {
	for {
		ch, err := l.readRune()
		if err != nil {
			return
		}
		if !unicode.IsSpace(ch) {
			l.unreadRune()
			return
		}
	}
}

// skipComment 跳过注释
func (l *Lexer) skipComment() {
	for {
		ch, err := l.readRune()
		if err != nil || ch == '\n' {
			return
		}
	}
}

// readRune 读取一个字符
func (l *Lexer) readRune() (rune, error) {
	ch, size, err := l.reader.ReadRune()
	l.lastReadRuneSize = size
	if ch == '\n' {
		l.lineLengths = append(l.lineLengths, l.column)
		l.line++
		l.column = 0
	} else {
		l.column++
	}

	return ch, err
}

// unreadRune 撤销读取的字符
func (l *Lexer) unreadRune() {
	l.reader.UnreadRune()
	if l.column > 0 {
		l.column--
	} else if l.line > 1 {
		// 回退到上一行的末尾
		l.line--
		l.column = l.lineLengths[l.line-1]     // 获取上一行的长度
		l.lineLengths = l.lineLengths[:l.line] // 移除当前行的长度记录
	}
}

// NextToken 读取下一个Token
func (l *Lexer) NextToken() (Token, error) {
	l.skipWhitespace()

	ch, err := l.readRune()
	if err == io.EOF {
		return Token{Type: EOF}, nil
	}

	// 检查是否为注释并跳过
	if ch == '/' {
		nextCh, _ := l.readRune()
		if nextCh == '/' {
			l.skipComment()
			return l.NextToken()
		} else {
			l.unreadRune()
		}
	}

	// 检查是否为字符串字面量
	if ch == '"' {
		var sb strings.Builder
		for {
			ch, err := l.readRune()
			if err != nil || ch == '"' {
				break
			}
			sb.WriteRune(ch)
		}
		return Token{Type: STRING, Value: sb.String(), Line: l.line, Column: l.column}, nil
	}

	// 检查是否为字母（可能是保留字或标识符）
	if isLetter(ch) {
		// startLine := l.line     // 记录标识符开始的行
		// startColumn := l.column // 记录标识符开始的列
		var sb strings.Builder
		sb.WriteRune(ch)
		for {
			ch, err := l.readRune()
			if err != nil || !isLetter(ch) && !isDigit(ch) && ch != '_' {
				l.unreadRune()
				break
			}
			sb.WriteRune(ch)
		}
		word := sb.String()
		tokenType := IDENTIFIER
		// position := Position{Line: startLine, Column: startColumn}

		if tType, ok := reservedWords[word]; ok {
			tokenType = tType
		}
		return Token{Type: tokenType, Value: word, Line: l.line, Column: l.column}, nil
	}

	// 检查是否为数字
	if isDigit(ch) {
		var sb strings.Builder
		sb.WriteRune(ch)
		for {
			ch, err := l.readRune()
			if err != nil || (!isDigit(ch) && ch != '.') {
				l.unreadRune()
				break
			}
			sb.WriteRune(ch)
		}
		tokenValue := sb.String()
		// 检查是否为实数
		if strings.Contains(tokenValue, ".") {
			return Token{Type: REAL, Value: tokenValue, Line: l.line, Column: l.column}, nil
		} else {
			return Token{Type: NUMBER, Value: tokenValue, Line: l.line, Column: l.column}, nil
		}
	}

	// 检查是否为运算符
	if tokenType, ok := operators[string(ch)]; ok {
		return Token{Type: tokenType, Value: string(ch), Line: l.line, Column: l.column}, nil
	}

	// 检查是否为分隔符
	if tokenType, ok := delimiters[string(ch)]; ok {
		return Token{Type: tokenType, Value: string(ch), Line: l.line, Column: l.column}, nil
	}

	// 未知字符
	return Token{}, fmt.Errorf(">>> 读取字符错误：未知字符 '%c', 位于 第 %d 行, 第 %d 列", ch, l.line, l.column)
}
