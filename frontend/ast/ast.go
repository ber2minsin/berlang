package ast

type NodeType string

const (
	ProgramType        NodeType = "Program"
	NumericLiteralType NodeType = "NumericLiteral"
	IdentifierType     NodeType = "Identifier"
	BinaryExprType     NodeType = "BinaryExpr"
	VarDeclType        NodeType = "VarDecl"
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
	Kind NodeType
	Body []Stmt
}

func (p *Program) GetKind() NodeType { return p.Kind }
func (p *Program) stmtNode()         {}

type BinaryExpr struct {
	Kind     NodeType
	Left     Expr
	Right    Expr
	Operator string
}

func (b *BinaryExpr) GetKind() NodeType { return b.Kind }
func (b *BinaryExpr) stmtNode()         {}
func (b *BinaryExpr) exprNode()         {}

type Identifier struct {
	Kind NodeType
	Name string
}

func (i *Identifier) GetKind() NodeType { return i.Kind }
func (i *Identifier) stmtNode()         {}
func (i *Identifier) exprNode()         {}

type NumericLiteral struct {
	Kind  NodeType
	Value string
}

func (n *NumericLiteral) GetKind() NodeType { return n.Kind }
func (n *NumericLiteral) stmtNode()         {}
func (n *NumericLiteral) exprNode()         {}

type VarDecl struct {
	Kind  NodeType
	Name  string
	Type  string // TODO actually define these types so we can check
	Value *Expr
}

func (n *VarDecl) GetKind() NodeType { return n.Kind }
func (n *VarDecl) stmtNode()         {}
func (n *VarDecl) exprNode()         {}

func NewVarDecl(name string, vartype string, value *Expr) *VarDecl {
	return &VarDecl{Kind: VarDeclType, Name: name, Type: vartype, Value: value}
}

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

func NewIdentifier(name string) *Identifier {
	return &Identifier{
		Kind: IdentifierType,
		Name: name,
	}
}

func NewNumericLiteral(value string) *NumericLiteral {
	return &NumericLiteral{
		Kind:  NumericLiteralType,
		Value: value,
	}
}
