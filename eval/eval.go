package eval

import (
	"fmt"

	"github.com/ibzsy/cardboard/object"
	"github.com/ibzsy/cardboard/parser/ast"
)

var NULL = &object.Null{}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalStatements(node.Statements, env)
	case *ast.UnboxStatement:
		return evalUnboxStatement(node, env)
	case *ast.PutStatement:
		return evalPutStatement(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	// Expressions
	case *ast.BoxExpression:
		return evalBoxExpression(node, env)
	case *ast.CallExpression:
		return evalCallExpression(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)
	case *ast.InfixExpression:
		return evalInfixExpression(node, env)
	case *ast.IntegerLiteral:
		return evalInteger(node)
	}

	// We've encountered an unknown word thats attempting ot be evaluated.
	return throwError("Unknown Word: <%s>", node.TokenLiteral())
}

func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement, env)
		switch result.Type() {
		case object.ERROR_OBJ:
			return result
		case object.UNBOX_OBJ:
			return result.(*object.Unbox).Value
		}
	}
	return result
}

func evalPrefixExpression(expr *ast.PrefixExpression, env *object.Environment) object.Object {
	operand := Eval(expr.Right, env)

	if isError(operand) {
		return operand
	}

	if operand.Type() != object.INTEGER {
		return throwError("Type error. Can't use <%s> Operator with <%s> Type.", expr.Operator, operand.Type())
	}

	value := operand.(*object.Integer).Value

	switch expr.Operator {
	case "+":
		return &object.Integer{Value: value}
	case "-":
		return &object.Integer{Value: -value}
	}

	return throwError("Unknown Operator: <%s>.", expr.Operator)
}

func evalInfixExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	// Eval Arguments
	left := Eval(node.Left, env)
	right := Eval(node.Right, env)

	if isError(left) {
		return left
	} else if isError(right) {
		return right
	}

	// Infix Operations only defined for integers for now.
	if left.Type() != object.INTEGER || right.Type() != object.INTEGER {
		return throwError("Type Mismatch: <%s><%s><%s>", left.Type(), node.Operator, right.Type())
	}

	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch node.Operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	}

	return NULL
}

func evalInteger(node *ast.IntegerLiteral) object.Object {
	return &object.Integer{Value: node.Value}
}

func evalUnboxStatement(unbox *ast.UnboxStatement, env *object.Environment) object.Object {
	val := Eval(unbox.NodeExpression, env)
	if isError(val) {
		return val
	}
	return &object.Unbox{Value: val}
}

func evalPutStatement(stmt *ast.PutStatement, env *object.Environment) object.Object {
	val := Eval(stmt.NodeExpression, env)
	if isError(val) {
		return val
	}
	return env.Set(stmt.NodeIdentifier.Value, val)
}

func evalIdentifier(ident *ast.Identifier, env *object.Environment) object.Object {
	obj, ok := env.Get(ident.Value)
	if !ok {
		return throwError("Unknown identifier: %s.", ident.TokenLiteral())
	}
	return obj
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result.Type() == object.ERROR_OBJ || result.Type() == object.UNBOX_OBJ {
			return result
		}
	}
	return result
}

func evalBoxExpression(box *ast.BoxExpression, env *object.Environment) object.Object {
	return &object.Box{Env: env, ParameterList: box.ParameterList, Body: box.Body}
}

func evalCallExpression(call *ast.CallExpression, env *object.Environment) object.Object {
	var arguments []object.Object

	box := Eval(call.Function, env)
	if isError(box) {
		return box
	}

	for _, arg := range call.Arguments {
		evaluated := Eval(arg, env)
		if isError(evaluated) {
			return evaluated
		}
		arguments = append(arguments, evaluated)
	}

	return applyBoxFunction(box, arguments)
}

func applyBoxFunction(box object.Object, args []object.Object) object.Object {
	fn, ok := box.(*object.Box)

	if !ok {
		return throwError("Type Mismatch Error. Expected Function. Got <%s>", box.Type())
	}

	fn.Env = object.CreateEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.ParameterList {
		fn.Env.Set(param.Value, args[paramIdx])
	}

	evaluated := Eval(fn.Body, fn.Env)

	if isError(evaluated) {
		return evaluated
	}

	if unbox, ok := evaluated.(*object.Unbox); ok {
		return unbox.Value
	}

	return evaluated
}

func throwError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if _, err := obj.(*object.Error); err {
		return true
	}
	return false
}
