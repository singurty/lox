package resolver

import (
	"errors"

	"github.com/singurty/lox/ast"
)

// function types enum
type functionType int
type classType int
const (
	NONE = iota
	FUNCTION
	METHOD
	INITIALIZER
	CLASS
	SUBCLASS
)

type Resolver struct {
	stack []map[string]bool
	Locals map[ast.Expr]int
	currentFunction functionType
	currentClass classType
	insideLoop bool
}

func NewResolver() *Resolver {
	resolver := &Resolver{
		stack: make([]map[string]bool, 0),
		Locals: make(map[ast.Expr]int),
		insideLoop: false,
		currentFunction: NONE,
		currentClass: NONE,
	}
	return resolver
}

// add new scope
func (r *Resolver) push(entry map[string]bool) {
	r.stack = append(r.stack, entry)
}

// remove current scope
func (r *Resolver) pop() map[string]bool {
	n := len(r.stack) -1
	entry := r.stack[n]
	r.stack = r.stack[:n]
	return entry
}

// get the map at the top of the Stack without removing it
func (r *Resolver) peek() map[string]bool {
	return r.stack[len(r.stack) - 1]
}

func (r *Resolver) beginScope() {
	r.push(make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.pop()
}

func (r *Resolver) Resolve(statements []ast.Stmt) error {
	for _, statement := range statements {
		err := r.resolveStmt(statement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveStmt(statement ast.Stmt) error {
	switch s := statement.(type) {
	case *ast.Block:
		err := r.blockStmt(s)
		if err != nil {
			return err
		}
	case *ast.Var:
		err := r.varStmt(s)
		if err != nil {
			return err
		}
	case *ast.Function:
		err := r.functionStmt(s)
		if err != nil {
			return err
		}
	case *ast.Class:
		err := r.classStmt(s)
		if err != nil {
			return err
		}
	case *ast.ExprStmt:
		err := r.expressionStmt(s)
		if err != nil {
			return err
		}
	case *ast.If:
		err := r.ifStmt(s)
		if err != nil {
			return err
		}
	case *ast.PrintStmt:
		err := r.printStmt(s)
		if err != nil {
			return err
		}
	case *ast.Return:
		err := r.returnStmt(s)
		if err != nil {
			return err
		}
	case *ast.While:
		err := r.whileStmt(s)
		if err != nil {
			return err
		}
	case *ast.For:
		err := r.forStmt(s)
		if err != nil {
			return err
		}
	case *ast.Break:
		err := r.breakStmt()
		if err != nil {
			return err
		}
	case *ast.Continue:
		err := r.continueStmt()
		if err != nil {
			return err
		}
	default:
		return errors.New("unknown statement")
	}
	return nil
}

func (r *Resolver) resolveExpr(expression ast.Expr) error {
	switch e := expression.(type) {
	case *ast.Variable:
		err := r.variableExpr(e)
		if err != nil {
			return err
		}
	case *ast.Assign:
		err := r.assignExpr(e)
		if err != nil {
			return err
		}
	case *ast.Set:
		err := r.setExpr(e)
		if err != nil {
			return err
		}
	case *ast.Binary:
		err := r.binaryExpr(e)
		if err != nil {
			return err
		}
	case *ast.Call:
		err := r.callExpr(e)
		if err != nil {
			return err
		}
	case *ast.Get:
		err := r.resolveExpr(e.Object)
		if err != nil {
			return err
		}
	case *ast.Lambda:
		err := r.lambdaExpr(e)
		if err != nil {
			return err
		}
	case *ast.Grouping:
		err := r.groupgingExpr(e)
		if err != nil {
			return err
		}
	case *ast.Literal:
		err := r.literalExpr(e)
		if err != nil {
			return err
		}
	case *ast.Logical:
		err := r.logicalExpr(e)
		if err != nil {
			return err
		}
	case *ast.Unary:
		err := r.unaryExpr(e)
		if err != nil {
			return err
		}
	case *ast.This:
		err := r.thisExpr(e)
		if err != nil {
			return err
		}
	case *ast.Super:
		err := r.superExpr(e)
		if err != nil {
			return err
		}
	default:
		return errors.New("unknown expression")
	}
	return nil
}

func (r *Resolver) blockStmt(block *ast.Block) error {
	r.beginScope()
	err := r.Resolve(block.Statements)
	if err != nil {
		return err
	}
	r.endScope()
	return nil
}

func (r *Resolver) varStmt(statement *ast.Var) error {
	err := r.declare(statement.Name.Lexeme)
	if err != nil {
		return err
	}
	if statement.Initializer != nil {
		err := r.resolveExpr(statement.Initializer)
		if err != nil {
			return err
		}
	}
	r.define(statement.Name.Lexeme)
	return nil
}

func (r *Resolver) classStmt(class *ast.Class) error {
	if class.SuperClass != nil && class.SuperClass.Name.Lexeme == class.Name.Lexeme {
		return errors.New("A class cannot inherit from itself.")
	}
	enclosing := r.currentClass
	if class.SuperClass != nil {
		r.currentClass = SUBCLASS
		r.beginScope()
		r.peek()["super"] = true
	} else {
		r.currentClass = CLASS
	}
	r.beginScope()
	r.peek()["this"] = true
	r.declare(class.Name.Lexeme)
	r.define(class.Name.Lexeme)
	if class.SuperClass != nil {
		r.variableExpr(class.SuperClass)
	}
	for _, method := range class.Methods {
		var declaration functionType
		if method.Name.Lexeme == "init" {
			declaration = INITIALIZER
		} else {
			declaration = METHOD
		}
		err := r.resolveFunction(method, declaration)
		if err != nil {
			return err
		}
	}
	r.endScope()
	if class.SuperClass != nil {
		r.endScope()
	}
	r.currentClass = enclosing
	return nil
}

func (r *Resolver) declare(name string) error {
	if len(r.stack) == 0 {
		return nil
	}
	if _, ok := r.peek()[name]; ok {
		return errors.New("A variable with the same name already exists in this scope")
	}
	r.peek()[name] = false
	return nil
}

func (r *Resolver) define(name string) {
	if len(r.stack) == 0 {
		return
	}
	r.peek()[name] = true
}

func (r *Resolver) variableExpr(expr *ast.Variable) error {
	if len(r.stack) > 0 {
		if value, ok := r.peek()[expr.Name.Lexeme]; ok && !value {
			return errors.New("Can't read local variable in its own initializer.")
		}
	}
	return r.resolveLocal(expr, expr.Name.Lexeme)
}

func (r *Resolver) resolveLocal(expr ast.Expr, name string) error {
	for i := len(r.stack) - 1; i >= 0; i-- {
		if _, ok := r.stack[i][name]; ok {
			r.Locals[expr] = len(r.stack) - 1 - i
			return nil
		}
	}
	return nil
}

func (r *Resolver) assignExpr(expr *ast.Assign) error {
	err := r.resolveExpr(expr.Value)
	if err != nil {
		return err
	}
	return r.resolveLocal(expr, expr.Name.Lexeme)
}

func (r *Resolver) setExpr(expr *ast.Set) error {
	err := r.resolveExpr(expr.Object)
	if err != nil {
		return err
	}
	err = r.resolveExpr(expr.Value)
	if err != nil {
		return err
	}
	return nil
}

func (r *Resolver) functionStmt(stmt *ast.Function) error {
	err := r.declare(stmt.Name.Lexeme)
	if err != nil {
		return err
	}
	r.define(stmt.Name.Lexeme)
	return r.resolveFunction(stmt, FUNCTION)
}

func (r *Resolver) resolveFunction(function *ast.Function, typeFunction functionType) error {
	enclosingFunction := r.currentFunction
	r.currentFunction = typeFunction
	r.beginScope()
	for _, param := range function.Parameters {
		err := r.declare(param.Lexeme)
		if err != nil {
			return err
		}
		r.define(param.Lexeme)
	}
	err := r.Resolve(function.Body)
	if err != nil {
		return err
	}
	r.endScope()
	r.currentFunction = enclosingFunction
	return nil
}

func (r *Resolver) expressionStmt(stmt *ast.ExprStmt) error {
	return r.resolveExpr(stmt.Expression)
}

func (r *Resolver) ifStmt(stmt *ast.If) error {
	err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return err
	}
	err = r.resolveStmt(stmt.ThenBranch)
	if err != nil {
		return err
	}
	if stmt.ElseBranch != nil {
		return r.resolveStmt(stmt.ElseBranch)
	}
	return nil
}

func (r *Resolver) printStmt(stmt *ast.PrintStmt) error {
	return r.resolveExpr(stmt.Expression)
}

func (r *Resolver) returnStmt(stmt *ast.Return) error {
	if r.currentFunction == NONE {
		return errors.New("Cannot return from top-level code.")
	}
	if stmt.Value != nil {
		if r.currentFunction == INITIALIZER {
			return errors.New("Cannot return from an initializer.")
		}
		return r.resolveExpr(stmt.Value)
	}
	return nil
}

func (r *Resolver) whileStmt(stmt *ast.While) error {
	prevStatus := r.insideLoop
	r.insideLoop = true
	err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return err
	}
	err = r.resolveStmt(stmt.Body)
	if err != nil {
		return err
	}
	r.insideLoop = prevStatus
	return nil
}

func (r *Resolver) forStmt(stmt *ast.For) error {
	prevStatus := r.insideLoop
	r.insideLoop = true
	err := r.resolveStmt(stmt.Initializer)
	if err != nil {
		return err
	}
	err = r.resolveExpr(stmt.Condition)
	if err != nil {
		return err
	}
	err = r.resolveExpr(stmt.Increment)
	if err != nil {
		return err
	}
	err = r.resolveStmt(stmt.Body)
	if err != nil {
		return err
	}
	r.insideLoop = prevStatus
	return nil
}

func (r *Resolver) breakStmt() error {
	if r.insideLoop {
		return nil
	}
	return errors.New("Cannot break outside of a loop")
}

func (r *Resolver) continueStmt() error {
	if r.insideLoop {
		return nil
	}
	return errors.New("Cannot continue outside of a loop")
}

func (r *Resolver) binaryExpr(expr *ast.Binary) error {
	err := r.resolveExpr(expr.Left)
	if err != nil {
		return err
	}
	return r.resolveExpr(expr.Right)
}

func (r *Resolver) callExpr(expr *ast.Call) error {
	err := r.resolveExpr(expr.Callee)
	if err != nil {
		return err
	}
	for _, argument := range expr.Arguments {
		err := r.resolveExpr(argument)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) lambdaExpr(expr *ast.Lambda) error {
	enclosingFunction := r.currentFunction
	r.currentFunction = FUNCTION
	r.beginScope()
	for _, param := range expr.Parameters {
		err := r.declare(param.Lexeme)
		if err != nil {
			return err
		}
		r.define(param.Lexeme)
	}
	err := r.Resolve(expr.Body)
	if err != nil {
		return err
	}
	r.endScope()
	r.currentFunction = enclosingFunction
	return nil
}

func (r *Resolver) groupgingExpr(expr *ast.Grouping) error {
	return r.resolveExpr(expr.Expression)
}

func (r *Resolver) literalExpr(expr *ast.Literal) error {
	return nil
}

func (r *Resolver) logicalExpr(expr *ast.Logical) error {
	err := r.resolveExpr(expr.Left)
	if err != nil {
		return err
	}
	return r.resolveExpr(expr.Right)
}

func (r *Resolver) unaryExpr(expr *ast.Unary) error {
	return r.resolveExpr(expr.Right)
}

func (r *Resolver) thisExpr(expr *ast.This) error {
	if r.currentClass == NONE {
		return errors.New("Cannot use \"this\" outside of a class.")
	}
	return r.resolveLocal(expr, expr.Keyword.Lexeme)
}

func (r *Resolver) superExpr(expr *ast.Super) error {
	if r.currentClass != SUBCLASS {
		return errors.New("Cannot use \"super\" outside of a subclass.")
	}
	r.resolveLocal(expr, expr.Keyword.Lexeme)
	return nil
}
