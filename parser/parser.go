package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func test() {
	exprs := []string{
		// bad expressions
		"1 *",
		"-1 * 2", // negative numbers are not allowed
		"1 (1)",
		"1 ?",
		"1 ? 2",
		"1 :",
		"1 : 2",
		"2 * (3 * (4 + 5)",
		"2 * (3 * (4 + 5)))",
		// good expressions
		"1 + 2",
		"1 + 2 * 3",
		"(1 + 2) * 3",
		"2 * (3 * (4 + 5))",
		"1 + (1 + 2) * 3",
		"3 * (1 + 2)",
		"3 * (1 + 2) + 1",
		"1 ",
		"1 > 2 ? 3 : 4",
		"1 > 2 ? (5 + 3) / 2 : 3 == 3 ? 6 * (5 - 2) : 0", // weee
	}

	for _, expr := range exprs {
		p, err := parseExpression(expr)
		if err != nil {
			fmt.Println(err)
		} else {
			n := p.expression(-1)
			if n.kind() == nodeError {
				fmt.Printf("%q: invalid expression\n", expr)
			} else {
				fmt.Printf("%q: %v\n", expr, n)
			}
		}
	}
}

var (
	numbers = "0123456789"
	symbols = "*/%+-=!<>|&?:"
)

func parseExpression(s string) (*parser, error) {
	var nodes []node
	var pin, pos int
	for pos < len(s) {
		// Skip spaces.
		for pos < len(s) && s[pos] == ' ' {
			pos++
			pin = pos
		}
		// Get parentheses.
		if pos < len(s) && (s[pos] == '(' || s[pos] == ')') {
			if n, err := newSymbolNode(string(s[pos])); err != nil {
				return nil, err
			} else {
				nodes = append(nodes, n)
				pos++
				pin = pos
				continue
			}
		}
		// Get int.
		for pos < len(s) && strings.IndexRune(numbers, rune(s[pos])) >= 0 {
			pos++
		}
		if pos > pin {
			if n, err := newIntNode(s[pin:pos]); err != nil {
				return nil, err
			} else {
				nodes = append(nodes, n)
				pin = pos
				continue
			}
		}
		// Get symbol.
		for pos < len(s) && strings.IndexRune(symbols, rune(s[pos])) >= 0 {
			pos++
		}
		if pos > pin {
			if n, err := newSymbolNode(s[pin:pos]); err != nil {
				return nil, err
			} else {
				nodes = append(nodes, n)
				pin = pos
				continue
			}
		}
		// We should not reach this point.
		if pos < len(s) {
			return nil, fmt.Errorf("Invalid char %q", string(s[pos]))
		}
	}
	return &parser{nodes: nodes}, nil
}

// ----------------------------------------------------------------------------

type nodeKind int

// node types
const (
	nodeError nodeKind = iota
	nodeBool
	nodeInt
	nodeMul        // *
	nodeDiv        // /
	nodeMod        // %
	nodeAdd        // +
	nodeSub        // - (binary)
	nodeEq         // ==
	nodeNotEq      // !=
	nodeGt         // >
	nodeGte        // >=
	nodeLt         // <
	nodeLte        // <=
	nodeOr         // ||
	nodeAnd        // &&
	nodeIf         // ?
	nodeElse       // :
	nodeNot        // !
	nodeLeftParen  // (
	nodeRightParen // )
)

// binding power
var nodePower = map[nodeKind]int{
	nodeMul:        9,
	nodeDiv:        9,
	nodeMod:        9,
	nodeAdd:        8,
	nodeSub:        8,
	nodeEq:         7,
	nodeNotEq:      7,
	nodeGt:         7,
	nodeGte:        7,
	nodeLt:         7,
	nodeLte:        7,
	nodeOr:         6,
	nodeAnd:        5,
	nodeIf:         4,
	nodeLeftParen:  3,
	nodeRightParen: 0,
	nodeElse:       0,
}

var defaultNodes = map[string]node{
	"*":  binaryOpNode{k: nodeMul},
	"/":  binaryOpNode{k: nodeDiv},
	"%":  binaryOpNode{k: nodeMod},
	"+":  binaryOpNode{k: nodeAdd},
	"-":  binaryOpNode{k: nodeSub},
	"==": binaryOpNode{k: nodeEq},
	"!=": binaryOpNode{k: nodeNotEq},
	">":  binaryOpNode{k: nodeGt},
	">=": binaryOpNode{k: nodeGte},
	"<":  binaryOpNode{k: nodeLt},
	"<=": binaryOpNode{k: nodeLte},
	"||": binaryOpNode{k: nodeOr},
	"&&": binaryOpNode{k: nodeAnd},
	"?":  ifNode{},
	":":  noopNode{k: nodeElse},
	"(":  leftParenNode{},
	")":  noopNode{k: nodeRightParen},
}

type parser struct {
	nodes []node
	index int
}

// Pratt's algorithm.
func (p *parser) expression(rbp int) node {
	if p.index >= len(p.nodes) {
		return invalidExpression
	}
	right := p.nodes[p.index]
	p.index++
	left := right.nud(p)
	for p.index < len(p.nodes) && rbp < nodePower[p.nodes[p.index].kind()] {
		right = p.nodes[p.index]
		p.index++
		left = right.led(p, left)
	}
	return left
}

// ----------------------------------------------------------------------------

type node interface {
	kind() nodeKind
	nud(*parser) node
	led(*parser, node) node
}

// ----------------------------------------------------------------------------

var invalidExpression = noopNode{nodeError, errors.New("Invalid expression")}

type noopNode struct {
	k   nodeKind
	err error
}

func (n noopNode) kind() nodeKind {
	return n.k
}

func (n noopNode) nud(p *parser) node {
	return invalidExpression
}

func (n noopNode) led(p *parser, left node) node {
	return invalidExpression
}

func (n noopNode) error() error {
	return n.err
}

// ----------------------------------------------------------------------------

func newIntNode(s string) (node, error) {
	if value, err := strconv.ParseInt(s, 10, 64); err == nil {
		return intNode(value), nil
	}
	return nil, fmt.Errorf("Invalid int %q", s)
}

type intNode int64

func (n intNode) kind() nodeKind {
	return nodeInt
}

func (n intNode) nud(p *parser) node {
	return n
}

func (n intNode) led(p *parser, left node) node {
	return invalidExpression
}

// ----------------------------------------------------------------------------

type boolNode bool

func (n boolNode) kind() nodeKind {
	return nodeBool
}

func (n boolNode) nud(p *parser) node {
	return n
}

func (n boolNode) led(p *parser, left node) node {
	return invalidExpression
}

// ----------------------------------------------------------------------------

func newSymbolNode(s string) (node, error) {
	if n, ok := defaultNodes[s]; ok {
		return n, nil
	}
	return nil, fmt.Errorf("Invalid symbol %q", s)
}

type binaryOpNode struct {
	k nodeKind
}

func (n binaryOpNode) kind() nodeKind {
	return n.k
}

func (n binaryOpNode) nud(p *parser) node {
	return invalidExpression
}

func (n binaryOpNode) led(p *parser, left node) node {
	right := p.expression(nodePower[n.kind()])
	switch x := left.(type) {
	case intNode:
		switch y := right.(type) {
		case intNode:
			switch n.k {
			case nodeMul:
				return x * y
			case nodeDiv:
				return x / y
			case nodeMod:
				return x % y
			case nodeAdd:
				return x + y
			case nodeSub:
				return x - y
			case nodeEq:
				return boolNode(x == y)
			case nodeNotEq:
				return boolNode(x != y)
			case nodeGt:
				return boolNode(x > y)
			case nodeGte:
				return boolNode(x >= y)
			case nodeLt:
				return boolNode(x < y)
			case nodeLte:
				return boolNode(x <= y)
			}
		}
	case boolNode:
		switch y := right.(type) {
		case boolNode:
			switch n.k {
			case nodeAnd:
				return boolNode(x && y)
			case nodeOr:
				return boolNode(x || y)
			case nodeEq:
				return boolNode(x == y)
			case nodeNotEq:
				return boolNode(x != y)
			}
		}
	}
	return invalidExpression
}

// ----------------------------------------------------------------------------

type ifNode struct{}

func (n ifNode) kind() nodeKind {
	return nodeIf
}

func (n ifNode) nud(p *parser) node {
	return invalidExpression
}

func (n ifNode) led(p *parser, left node) node {
	if cond, ok := left.(boolNode); ok {
		rv1 := p.expression(0)
		if p.index >= len(p.nodes) || p.nodes[p.index].kind() != nodeElse {
			return invalidExpression
		}
		p.index++
		rv2 := p.expression(0)
		if cond {
			return rv1
		}
		return rv2
	}
	return invalidExpression
}

// ----------------------------------------------------------------------------

type leftParenNode struct{}

func (n leftParenNode) kind() nodeKind {
	return nodeLeftParen
}

func (n leftParenNode) nud(p *parser) node {
	right := p.expression(0)
	if p.index >= len(p.nodes) || p.nodes[p.index].kind() != nodeRightParen {
		return invalidExpression
	}
	p.index++
	return right
}

func (n leftParenNode) led(p *parser, left node) node {
	return invalidExpression
}
