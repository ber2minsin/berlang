package ast

type NodeType string

const (
	ProgramType        NodeType = "Program"
	NumericLiteralType NodeType = "NumericLiteral"
	IdentifierType     NodeType = "Identifier"
	BinaryExprType     NodeType = "BinaryExpr"
)

type Node interface {
	GetKind() NodeType
}

type Stmt interface {
	Node
	stmtNode()
}

type Expr interface {
	Stmt
	exprNode()
}

type Program struct {
	Kind NodeType `json:"kind"`
	Body []Stmt   `json:"body"`
}

func (p *Program) GetKind() NodeType { return p.Kind }
func (p *Program) stmtNode()         {}

type BinaryExpr struct {
	Kind     NodeType `json:"kind"`
	Left     Expr     `json:"left"`
	Right    Expr     `json:"right"`
	Operator string   `json:"operator"`
}

func (b *BinaryExpr) GetKind() NodeType { return b.Kind }
func (b *BinaryExpr) stmtNode()         {}
func (b *BinaryExpr) exprNode()         {}

type Identifier struct {
	Kind   NodeType `json:"kind"`
	Symbol string   `json:"symbol"`
}

func (i *Identifier) GetKind() NodeType { return i.Kind }
func (i *Identifier) stmtNode()         {}
func (i *Identifier) exprNode()         {}

type NumericLiteral struct {
	Kind  NodeType `json:"kind"`
	Value string   `json:"value"`
}

func (n *NumericLiteral) GetKind() NodeType { return n.Kind }
func (n *NumericLiteral) stmtNode()         {}
func (n *NumericLiteral) exprNode()         {}

func NewProgram() *Program {
	return &Program{
		Kind: ProgramType,
		Body: make([]Stmt, 0),
	}
}

func NewBinaryExpr(left Expr, right Expr, operator string) *BinaryExpr {
	return &BinaryExpr{
		Kind:     BinaryExprType,
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

func NewIdentifier(symbol string) *Identifier {
	return &Identifier{
		Kind:   IdentifierType,
		Symbol: symbol,
	}
}

func NewNumericLiteral(value string) *NumericLiteral {
	return &NumericLiteral{
		Kind:  NumericLiteralType,
		Value: value,
	}
}
