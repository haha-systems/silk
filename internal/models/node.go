package models

type NodeType string

const (
	NodeTypeProgram         NodeType = "Program"
	NodeTypeNumber          NodeType = "Number"
	NodeTypeVariable        NodeType = "Variable"
	NodeTypeBinaryExpr      NodeType = "BinaryExpression"
	NodeTypeAssignment      NodeType = "Assignment"
	NodeTypeIf              NodeType = "IfStatement"
	NodeTypeFunctionCall    NodeType = "FunctionCall"
	NodeTypeReturnStatement NodeType = "ReturnStatement"
)

type Node interface {
	GetType() NodeType
}

type Program struct {
	Body []Node
}

func (p *Program) GetType() NodeType {
	return NodeTypeProgram
}

type Number struct {
	Value float64
}

func (n *Number) GetType() NodeType {
	return NodeTypeNumber
}

type Variable struct {
	Name string
}

func (v *Variable) GetType() NodeType {
	return NodeTypeVariable
}

type BinaryExpression struct {
	Operator string
	Left     Node
	Right    Node
}

func (be *BinaryExpression) GetType() NodeType {
	return NodeTypeBinaryExpr
}

type Assignment struct {
	Variable *Variable
	Value    Node
}

func (a *Assignment) GetType() NodeType {
	return NodeTypeAssignment
}

type IfStatement struct {
	Condition  Node
	Consequent Node
	Alternate  Node
}

func (ifs *IfStatement) GetType() NodeType {
	return NodeTypeIf
}

type String struct {
	Value string
}

func (s *String) GetType() NodeType {
	return "String"
}

type ComparisonExpression struct {
	Operator string
	Left     Node
	Right    Node
}

func (ce *ComparisonExpression) GetType() NodeType {
	return "ComparisonExpression"
}

type ParallelBlock struct {
	Body []Node
}

func (pb *ParallelBlock) GetType() NodeType {
	return "ParallelBlock"
}

type FunctionCall struct {
	Name string
	Args []Node
}

func (fc *FunctionCall) GetType() NodeType {
	return "FunctionCall"
}

type FunctionDeclaration struct {
	Name       string
	Parameters []*Variable
	Body       []Node
}

func (fd *FunctionDeclaration) GetType() NodeType {
	return "FunctionDeclaration"
}

type ForLoop struct {
	Initialization Node
	Condition      Node
	Post           Node
	Body           []Node
}

func (fl *ForLoop) GetType() NodeType {
	return "ForLoop"
}

type WhileLoop struct {
	Condition Node
	Body      []Node
}

func (wl *WhileLoop) GetType() NodeType {
	return "WhileLoop"
}

type ReturnStatement struct {
	Value Node
}

func (rs *ReturnStatement) GetType() NodeType {
	return "ReturnStatement"
}
