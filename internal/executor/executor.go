package executor

import (
	"errors"
	"fmt"
	"sync"

	"silk/internal/models"
)

type Executor struct {
	env       map[string]interface{}
	functions map[string]*models.FunctionDeclaration
	builtins  map[string]func(args []interface{}) (interface{}, error)
}

func NewExecutor() *Executor {
	return &Executor{
		env:       make(map[string]interface{}),
		functions: make(map[string]*models.FunctionDeclaration),
		builtins:  make(map[string]func(args []interface{}) (interface{}, error)),
	}
}

func (e *Executor) Execute(node models.Node) (interface{}, error) {
	switch n := node.(type) {

	case *models.Program:
		var result interface{}
		for _, stmt := range n.Body {
			res, err := e.Execute(stmt)
			if err != nil {
				return nil, err
			}
			result = res
		}
		return result, nil

	case *models.Number:
		return n.Value, nil

	case *models.Variable:
		val, ok := e.env[n.Name]
		if !ok {
			return nil, fmt.Errorf("undefined variable: %s", n.Name)
		}
		return val, nil

	case *models.Assignment:
		val, err := e.Execute(n.Value)
		if err != nil {
			return nil, err
		}
		e.env[n.Variable.Name] = val
		return val, nil

	case *models.BinaryExpression:
		left, err := e.Execute(n.Left)
		if err != nil {
			return nil, err
		}
		right, err := e.Execute(n.Right)
		if err != nil {
			return nil, err
		}

		switch n.Operator {
		case "+":
			return e.add(left, right)
		case "-":
			return e.subtract(left, right)
		case "*":
			return e.multiply(left, right)
		case "/":
			return e.divide(left, right)
		default:
			return nil, fmt.Errorf("unknown operator: %s", n.Operator)
		}

	case *models.IfStatement:
		condition, err := e.Execute(n.Condition)
		if err != nil {
			return nil, err
		}
		condBool, ok := condition.(bool)
		if !ok {
			return nil, errors.New("condition must evaluate to a boolean")
		}
		if condBool {
			return e.Execute(n.Consequent)
		} else if n.Alternate != nil {
			return e.Execute(n.Alternate)
		}
		return nil, nil

	case *models.String:
		return n.Value, nil

	case *models.ComparisonExpression:
		left, err := e.Execute(n.Left)
		if err != nil {
			return nil, err
		}
		right, err := e.Execute(n.Right)
		if err != nil {
			return nil, err
		}
		leftNum, ok1 := left.(float64)
		rightNum, ok2 := right.(float64)
		if !ok1 || !ok2 {
			return nil, errors.New("operands must be numbers")
		}
		switch n.Operator {
		case ">":
			return leftNum > rightNum, nil
		case "<":
			return leftNum < rightNum, nil
		case "==":
			return leftNum == rightNum, nil
		default:
			return nil, fmt.Errorf("unknown comparison operator: %s", n.Operator)
		}

	case *models.ParallelBlock:
		var wg sync.WaitGroup
		errorsCh := make(chan error, len(n.Body))
		for _, childNode := range n.Body {
			wg.Add(1)
			go func(node models.Node) {
				defer wg.Done()
				_, err := e.Execute(node)
				if err != nil {
					errorsCh <- err
				}
			}(childNode)
		}
		wg.Wait()
		close(errorsCh)
		if len(errorsCh) > 0 {
			// Handle errors (for simplicity, return the first one)
			return nil, <-errorsCh
		}
		return nil, nil

	case *models.FunctionDeclaration:
		e.functions[n.Name] = n
		return nil, nil

	case *models.FunctionCall:
		// Check if it's a built-in function
		if builtin, ok := e.builtins[n.Name]; ok {
			args := []interface{}{}
			for _, argNode := range n.Args {
				argVal, err := e.Execute(argNode)
				if err != nil {
					return nil, err
				}
				args = append(args, argVal)
			}
			return builtin(args)
		}

		// Handle user-defined functions
		function, ok := e.functions[n.Name]
		if !ok {
			return nil, fmt.Errorf("undefined function: %s", n.Name)
		}

		// Create a new environment for the function
		newEnv := make(map[string]interface{})
		// Evaluate arguments
		for i, param := range function.Parameters {
			argVal, err := e.Execute(n.Args[i])
			if err != nil {
				return nil, err
			}
			newEnv[param.Name] = argVal
		}

		// Save the current environment and set the new one
		originalEnv := e.env
		e.env = newEnv

		// Execute the function body
		var result interface{}
		for _, stmt := range function.Body {
			res, err := e.Execute(stmt)
			if err != nil {
				e.env = originalEnv
				return nil, err
			}
			result = res
		}

		// Restore the original environment
		e.env = originalEnv
		return result, nil

	case *models.ForLoop:
		// Execute the initialization part
		_, err := e.Execute(n.Initialization)
		if err != nil {
			return nil, err
		}

		// Loop while the condition is true
		for {
			condition, err := e.Execute(n.Condition)
			if err != nil {
				return nil, err
			}
			condBool, ok := condition.(bool)
			if !ok {
				return nil, errors.New("condition must evaluate to a boolean")
			}
			if !condBool {
				break
			}

			// Execute the loop body
			for _, stmt := range n.Body {
				_, err := e.Execute(stmt)
				if err != nil {
					return nil, err
				}
			}

			// Execute the post iteration statement
			_, err = e.Execute(n.Post)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil

	case *models.WhileLoop:
		// Loop while the condition is true
		for {
			condition, err := e.Execute(n.Condition)
			if err != nil {
				return nil, err
			}
			condBool, ok := condition.(bool)
			if !ok {
				return nil, errors.New("condition must evaluate to a boolean")
			}
			if !condBool {
				break
			}

			// Execute the loop body
			for _, stmt := range n.Body {
				_, err := e.Execute(stmt)
				if err != nil {
					return nil, err
				}
			}
		}
		return nil, nil

	default:
		return nil, fmt.Errorf("unknown node type: %T", n)
	}
}

func (e *Executor) Env() map[string]interface{} {
	return e.env
}

func (e *Executor) RegisterFunction(name string, function *models.FunctionDeclaration) {
	if e.functions == nil {
		e.functions = make(map[string]*models.FunctionDeclaration)
	}
	e.functions[name] = function
}

func (e *Executor) RegisterBuiltin(name string, function func(args []interface{}) (interface{}, error)) {
	if e.builtins == nil {
		e.builtins = make(map[string]func(args []interface{}) (interface{}, error))
	}
	e.builtins[name] = function
}

func (e *Executor) add(a, b interface{}) (interface{}, error) {
	switch a := a.(type) {
	case float64:
		b := b.(float64)
		return a + b, nil
	case string:
		b := b.(string)
		return a + b, nil
	default:
		return nil, errors.New("unsupported types for addition")
	}
}

func (e *Executor) subtract(a, b interface{}) (interface{}, error) {
	aNum, ok1 := a.(float64)
	bNum, ok2 := b.(float64)
	if !ok1 || !ok2 {
		return nil, errors.New("operands must be numbers")
	}
	return aNum - bNum, nil
}

func (e *Executor) multiply(a, b interface{}) (interface{}, error) {
	aNum, ok1 := a.(float64)
	bNum, ok2 := b.(float64)
	if !ok1 || !ok2 {
		return nil, errors.New("operands must be numbers")
	}
	return aNum * bNum, nil
}

func (e *Executor) divide(a, b interface{}) (interface{}, error) {
	aNum, ok1 := a.(float64)
	bNum, ok2 := b.(float64)
	if !ok1 || !ok2 {
		return nil, errors.New("operands must be numbers")
	}
	if bNum == 0 {
		return nil, errors.New("division by zero")
	}
	return aNum / bNum, nil
}
