package golang

import (
	"go/ast"
	"go/parser"
	"go/token"

	"fm.tul.cz/dupl/syntax"
)

const (
	BadNode = iota
	File
	ArrayType
	AssignStmt
	BasicLit
	BinaryExpr
	BlockStmt
	BranchStmt
	CallExpr
	CaseClause
	ChanType
	CommClause
	CompositeLit
	DeclStmt
	DeferStmt
	Ellipsis
	EmptyStmt
	ExprStmt
	Field
	FieldList
	ForStmt
	FuncDecl
	FuncLit
	FuncType
	GenDecl
	GoStmt
	Ident
	IfStmt
	IncDecStmt
	IndexExpr
	InterfaceType
	KeyValueExpr
	LabeledStmt
	MapType
	ParenExpr
	RangeStmt
	ReturnStmt
	SelectStmt
	SelectorExpr
	SendStmt
	SliceExpr
	StarExpr
	StructType
	SwitchStmt
	TypeAssertExpr
	TypeSpec
	TypeSwitchStmt
	UnaryExpr
	ValueSpec
)

// Parse the given file and return uniform syntax tree.
func Parse(filename string) (*syntax.Node, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return nil, err
	}
	return trans(file), nil
}

// trans transforms given golang AST to uniform tree structure.
func trans(node ast.Node) (o *syntax.Node) {
	o = syntax.NewNode()
	o.Pos, o.End = node.Pos(), node.End()

	switch n := node.(type) {
	case *ast.ArrayType:
		o.Type = ArrayType
		if n.Len != nil {
			o.AddChildren(trans(n.Len))
		}
		o.AddChildren(trans(n.Elt))

	case *ast.AssignStmt:
		o.Type = AssignStmt
		for _, e := range n.Rhs {
			o.AddChildren(trans(e))
		}

		for _, e := range n.Lhs {
			o.AddChildren(trans(e))
		}

	case *ast.BasicLit:
		o.Type = BasicLit

	case *ast.BinaryExpr:
		o.Type = BinaryExpr
		o.AddChildren(trans(n.X), trans(n.Y))

	case *ast.BlockStmt:
		o.Type = BlockStmt
		for _, stmt := range n.List {
			o.AddChildren(trans(stmt))
		}

	case *ast.BranchStmt:
		o.Type = BranchStmt
		if n.Label != nil {
			o.AddChildren(trans(n.Label))
		}

	case *ast.CallExpr:
		o.Type = CallExpr
		o.AddChildren(trans(n.Fun))
		for _, arg := range n.Args {
			o.AddChildren(trans(arg))
		}

	case *ast.CaseClause:
		o.Type = CaseClause
		for _, e := range n.List {
			o.AddChildren(trans(e))
		}
		for _, stmt := range n.Body {
			o.AddChildren(trans(stmt))
		}

	case *ast.ChanType:
		o.Type = ChanType
		o.AddChildren(trans(n.Value))

	case *ast.CommClause:
		o.Type = CommClause
		if n.Comm != nil {
			o.AddChildren(trans(n.Comm))
		}
		for _, stmt := range n.Body {
			o.AddChildren(trans(stmt))
		}

	case *ast.CompositeLit:
		o.Type = CompositeLit
		if n.Type != nil {
			o.AddChildren(trans(n.Type))
		}
		for _, e := range n.Elts {
			o.AddChildren(trans(e))
		}

	case *ast.DeclStmt:
		o.Type = DeclStmt
		o.AddChildren(trans(n.Decl))

	case *ast.DeferStmt:
		o.Type = DeferStmt
		o.AddChildren(trans(n.Call))

	case *ast.Ellipsis:
		o.Type = Ellipsis
		if n.Elt != nil {
			o.AddChildren(trans(n.Elt))
		}

	case *ast.EmptyStmt:
		o.Type = EmptyStmt

	case *ast.ExprStmt:
		o.Type = ExprStmt
		o.AddChildren(trans(n.X))

	case *ast.Field:
		o.Type = Field
		for _, name := range n.Names {
			o.AddChildren(trans(name))
		}
		o.AddChildren(trans(n.Type))

	case *ast.FieldList:
		o.Type = FieldList
		for _, field := range n.List {
			o.AddChildren(trans(field))
		}

	case *ast.File:
		o.Type = File
		for _, decl := range n.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
				// skip import declarations
				continue
			}
			o.AddChildren(trans(decl))
		}

	case *ast.ForStmt:
		o.Type = ForStmt
		if n.Init != nil {
			o.AddChildren(trans(n.Init))
		}
		if n.Cond != nil {
			o.AddChildren(trans(n.Cond))
		}
		if n.Post != nil {
			o.AddChildren(trans(n.Post))
		}
		o.AddChildren(trans(n.Body))

	case *ast.FuncDecl:
		o.Type = FuncDecl
		if n.Recv != nil {
			o.AddChildren(trans(n.Recv))
		}
		o.AddChildren(trans(n.Name), trans(n.Type))
		if n.Body != nil {
			o.AddChildren(trans(n.Body))
		}

	case *ast.FuncLit:
		o.Type = FuncLit
		o.AddChildren(trans(n.Type), trans(n.Body))

	case *ast.FuncType:
		o.Type = FuncType
		o.AddChildren(trans(n.Params))
		if n.Results != nil {
			o.AddChildren(trans(n.Results))
		}

	case *ast.GenDecl:
		o.Type = GenDecl
		for _, spec := range n.Specs {
			o.AddChildren(trans(spec))
		}

	case *ast.GoStmt:
		o.Type = GoStmt
		o.AddChildren(trans(n.Call))

	case *ast.Ident:
		o.Type = Ident

	case *ast.IfStmt:
		o.Type = IfStmt
		if n.Init != nil {
			o.AddChildren(trans(n.Init))
		}
		o.AddChildren(trans(n.Cond), trans(n.Body))
		if n.Else != nil {
			o.AddChildren(trans(n.Else))
		}

	case *ast.IncDecStmt:
		o.Type = IncDecStmt
		o.AddChildren(trans(n.X))

	case *ast.IndexExpr:
		o.Type = IndexExpr
		o.AddChildren(trans(n.X), trans(n.Index))

	case *ast.InterfaceType:
		o.Type = InterfaceType
		o.AddChildren(trans(n.Methods))

	case *ast.KeyValueExpr:
		o.Type = KeyValueExpr
		o.AddChildren(trans(n.Key), trans(n.Value))

	case *ast.LabeledStmt:
		o.Type = LabeledStmt
		o.AddChildren(trans(n.Label), trans(n.Stmt))

	case *ast.MapType:
		o.Type = MapType
		o.AddChildren(trans(n.Key), trans(n.Value))

	case *ast.ParenExpr:
		o.Type = ParenExpr
		o.AddChildren(trans(n.X))

	case *ast.RangeStmt:
		o.Type = RangeStmt
		if n.Key != nil {
			o.AddChildren(trans(n.Key))
		}
		if n.Value != nil {
			o.AddChildren(trans(n.Value))
		}
		o.AddChildren(trans(n.X), trans(n.Body))

	case *ast.ReturnStmt:
		o.Type = ReturnStmt
		for _, e := range n.Results {
			o.AddChildren(trans(e))
		}

	case *ast.SelectStmt:
		o.Type = SelectStmt
		o.AddChildren(trans(n.Body))

	case *ast.SelectorExpr:
		o.Type = SelectorExpr
		o.AddChildren(trans(n.X), trans(n.Sel))

	case *ast.SendStmt:
		o.Type = SendStmt
		o.AddChildren(trans(n.Chan), trans(n.Value))

	case *ast.SliceExpr:
		o.Type = SliceExpr
		o.AddChildren(trans(n.X))
		if n.Low != nil {
			o.AddChildren(trans(n.Low))
		}
		if n.High != nil {
			o.AddChildren(trans(n.High))
		}
		if n.Max != nil {
			o.AddChildren(trans(n.Max))
		}

	case *ast.StarExpr:
		o.Type = StarExpr
		o.AddChildren(trans(n.X))

	case *ast.StructType:
		o.Type = StructType
		o.AddChildren(trans(n.Fields))

	case *ast.SwitchStmt:
		o.Type = SwitchStmt
		if n.Init != nil {
			o.AddChildren(trans(n.Init))
		}
		if n.Tag != nil {
			o.AddChildren(trans(n.Tag))
		}
		o.AddChildren(trans(n.Body))

	case *ast.TypeAssertExpr:
		o.Type = TypeAssertExpr
		o.AddChildren(trans(n.X))
		if n.Type != nil {
			o.AddChildren(trans(n.Type))
		}

	case *ast.TypeSpec:
		o.Type = TypeSpec
		o.AddChildren(trans(n.Name), trans(n.Type))

	case *ast.TypeSwitchStmt:
		o.Type = TypeSwitchStmt
		if n.Init != nil {
			o.AddChildren(trans(n.Init))
		}
		o.AddChildren(trans(n.Assign), trans(n.Body))

	case *ast.UnaryExpr:
		o.Type = UnaryExpr
		o.AddChildren(trans(n.X))

	case *ast.ValueSpec:
		o.Type = ValueSpec
		for _, name := range n.Names {
			o.AddChildren(trans(name))
		}
		if n.Type != nil {
			o.AddChildren(trans(n.Type))
		}
		for _, val := range n.Values {
			o.AddChildren(trans(val))
		}

	default:
		o.Type = BadNode

	}

	return o
}
