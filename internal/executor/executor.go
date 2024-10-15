package executor

import (
	"errors"
	"fmt"
	"runtime"
	"sync"

	"silk/internal/models"
)

// Environment represents a single scope of variable bindings.
type Environment struct {
	variables  map[string]interface{}
	isReusable bool
}

// Executor is responsible for executing AST nodes and managing environments and functions.
type Executor struct {
	envStack      []Environment                                            // Stack of environments to handle variable scoping.
	functions     map[string]*models.FunctionDeclaration                   // Map of user-defined functions.
	builtins      map[string]func(args []interface{}) (interface{}, error) // Map of built-in functions.
	builtinCache  map[string]func(args []interface{}) (interface{}, error) // Cache for frequently used built-in functions.
	envPool       []Environment                                            // Pool of reusable environments.
	maxGoroutines int                                                      // Maximum number of concurrent goroutines.
	sem           chan struct{}                                            // Semaphore to control goroutine concurrency.
}

// NewExecutor creates a new Executor with an initial environment.
func NewExecutor() *Executor {
	maxGoroutines := runtime.NumCPU() // Set the limit for the number of concurrent goroutines to the number of logical processors.
	return &Executor{
		envStack:      []Environment{{variables: make(map[string]interface{}), isReusable: false}},
		functions:     make(map[string]*models.FunctionDeclaration),
		builtins:      make(map[string]func(args []interface{}) (interface{}, error)),
		builtinCache:  make(map[string]func(args []interface{}) (interface{}, error)),
		envPool:       []Environment{},
		maxGoroutines: maxGoroutines,
		sem:           make(chan struct{}, maxGoroutines),
	}
}

// Execute executes a given AST node and returns the result or an error.
func (e *Executor) Execute(node models.Node) (interface{}, error) {
	switch n := node.(type) {

	case *models.Program:
		// Execute each statement in the program sequentially.
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
		// Return the numeric value.
		return n.Value, nil

	case *models.Variable:
		// Retrieve the value of a variable from the current environment.
		val, ok := e.currentEnv().variables[n.Name]
		if !ok {
			return nil, fmt.Errorf("undefined variable: %s", n.Name)
		}
		return val, nil

	case *models.Assignment:
		// Evaluate the value and assign it to the variable in the current environment.
		val, err := e.Execute(n.Value)
		if err != nil {
			return nil, err
		}
		e.currentEnv().variables[n.Variable.Name] = val
		return val, nil

	case *models.BinaryExpression:
		// Validate operator before evaluating operands to avoid unnecessary computations.
		if !e.isValidOperator(n.Operator) {
			return nil, fmt.Errorf("unknown operator: %s", n.Operator)
		}

		// Evaluate both sides of the binary expression and perform the operation.
		left, err := e.Execute(n.Left)
		if err != nil {
			return nil, err
		}
		right, err := e.Execute(n.Right)
		if err != nil {
			return nil, err
		}

		// Check if both operands are numbers before performing the operation.
		leftNum, ok1 := left.(float64)
		rightNum, ok2 := right.(float64)
		if !ok1 || !ok2 {
			return nil, errors.New("operands must be numbers")
		}

		return e.handleBinaryOperation(n.Operator, leftNum, rightNum)

	case *models.IfStatement:
		// Evaluate the condition and execute the appropriate branch.
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
		// Return the string value.
		return n.Value, nil

	case *models.ComparisonExpression:
		// Evaluate both sides of the comparison and perform the comparison operation.
		left, err := e.Execute(n.Left)
		if err != nil {
			return nil, err
		}
		right, err := e.Execute(n.Right)
		if err != nil {
			return nil, err
		}

		// Check if both operands are numbers before performing the comparison.
		leftNum, ok1 := left.(float64)
		rightNum, ok2 := right.(float64)
		if !ok1 || !ok2 {
			return nil, errors.New("operands must be numbers")
		}

		return e.handleComparison(n.Operator, leftNum, rightNum)

	case *models.ParallelBlock:
		// Execute each statement in parallel using goroutines, with a limit on concurrency.
		var wg sync.WaitGroup
		errors := []error{}
		var mu sync.Mutex
		for _, childNode := range n.Body {
			e.sem <- struct{}{} // Acquire a slot
			wg.Add(1)
			go func(node models.Node) {
				defer wg.Done()
				defer func() { <-e.sem }() // Release the slot
				_, err := e.Execute(node)
				if err != nil {
					mu.Lock()
					errors = append(errors, err)
					mu.Unlock()
				}
			}(childNode)
		}
		wg.Wait()
		if len(errors) > 0 {
			return nil, fmt.Errorf("multiple errors occurred: %v", errors)
		}
		return nil, nil

	case *models.FunctionDeclaration:
		// Register a user-defined function.
		e.functions[n.Name] = n
		return nil, nil

	case *models.FunctionCall:
		// Handle a function call, either built-in or user-defined.
		return e.handleFunctionCall(n)

	case *models.ForLoop:
		// Handle a for loop, including initialization, condition check, and post iteration.
		return e.handleForLoop(n)

	case *models.WhileLoop:
		// Handle a while loop, executing while the condition is true.
		return e.handleWhileLoop(n)

	default:
		return nil, fmt.Errorf("unknown node type: %T", n)
	}
}

// currentEnv returns the current environment from the top of the stack.
func (e *Executor) currentEnv() *Environment {
	return &e.envStack[len(e.envStack)-1]
}

// pushEnv adds a new environment to the stack, reusing one from the pool if available.
func (e *Executor) pushEnv() {
	var newEnv Environment
	if len(e.envPool) > 0 {
		newEnv = e.envPool[len(e.envPool)-1]
		e.envPool = e.envPool[:len(e.envPool)-1]
		newEnv.variables = make(map[string]interface{}) // Reset the environment variables.
	} else {
		newEnv = Environment{variables: make(map[string]interface{}), isReusable: true}
	}
	e.envStack = append(e.envStack, newEnv)
}

// popEnv removes the top environment from the stack and adds it back to the pool if reusable.
func (e *Executor) popEnv() {
	env := e.envStack[len(e.envStack)-1]
	e.envStack = e.envStack[:len(e.envStack)-1]
	if env.isReusable {
		e.envPool = append(e.envPool, env)
	}
}

// Env returns the environment stack.
func (e *Executor) Env() []Environment {
	return e.envStack
}

// CurrentEnv returns the current environment from the top of the stack.
func (e *Executor) CurrentEnv() Environment {
	return *e.currentEnv()
}

// EnvValue retrieves the value of a variable from the current environment.
func (e *Executor) EnvValue(name string) (interface{}, error) {
	val, ok := e.currentEnv().variables[name]
	if !ok {
		return nil, fmt.Errorf("undefined variable: %s", name)
	}
	return val, nil
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

// handleFunctionCall executes a function call, supporting both built-in and user-defined functions.
func (e *Executor) handleFunctionCall(n *models.FunctionCall) (interface{}, error) {
	// Check if it's cached in the built-in function cache.
	if cachedBuiltin, ok := e.builtinCache[n.Name]; ok {
		args := []interface{}{}
		for _, argNode := range n.Args {
			argVal, err := e.Execute(argNode)
			if err != nil {
				return nil, err
			}
			args = append(args, argVal)
		}
		return cachedBuiltin(args)
	}

	// Check if it's a built-in function.
	if builtin, ok := e.builtins[n.Name]; ok {
		// Cache the built-in function for future calls.
		e.builtinCache[n.Name] = builtin
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

	// Handle user-defined function.
	function, ok := e.functions[n.Name]
	if !ok {
		return nil, fmt.Errorf("undefined function: %s", n.Name)
	}

	// Check if the number of arguments matches the number of parameters.
	if len(n.Args) != len(function.Parameters) {
		return nil, fmt.Errorf("function %s expects %d arguments, but got %d", n.Name, len(function.Parameters), len(n.Args))
	}

	// Create a new environment for the function call.
	e.pushEnv()
	defer e.popEnv()
	for i, param := range function.Parameters {
		argVal, err := e.Execute(n.Args[i])
		if err != nil {
			return nil, err
		}
		e.currentEnv().variables[param.Name] = argVal
	}

	// Execute the function body.
	var result interface{}
	// Instead of using retStmt, let's directly check the type and break if necessary
	for _, stmt := range function.Body {
		res, err := e.Execute(stmt)
		if err != nil {
			return nil, err
		}
		if _, ok := stmt.(*models.ReturnStatement); ok {
			result = res
			break
		}
		result = res
	}

	return result, nil
}

// handleBinaryOperation performs arithmetic operations on two operands.
func (e *Executor) handleBinaryOperation(operator string, left, right float64) (interface{}, error) {
	switch operator {
	case "+":
		return left + right, nil
	case "-":
		return left - right, nil
	case "*":
		return left * right, nil
	case "/":
		if right == 0 {
			return nil, errors.New("division by zero")
		}
		return left / right, nil
	default:
		return nil, fmt.Errorf("unknown operator: %s", operator)
	}
}

// handleComparison performs comparison operations on two operands.
func (e *Executor) handleComparison(operator string, left, right float64) (interface{}, error) {
	switch operator {
	case ">":
		return left > right, nil
	case "<":
		return left < right, nil
	case "==":
		return left == right, nil
	default:
		return nil, fmt.Errorf("unknown comparison operator: %s", operator)
	}
}

// handleForLoop executes a for loop, managing initialization, condition, and post-iteration.
func (e *Executor) handleForLoop(n *models.ForLoop) (interface{}, error) {
	// Execute the initialization part of the loop.
	_, err := e.Execute(n.Initialization)
	if err != nil {
		return nil, err
	}

	// Loop while the condition is true.
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

		// Execute the loop body.
		for _, stmt := range n.Body {
			_, err := e.Execute(stmt)
			if err != nil {
				return nil, err
			}
		}

		// Execute the post iteration statement.
		_, err = e.Execute(n.Post)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// handleWhileLoop executes a while loop, continuing as long as the condition is true.
func (e *Executor) handleWhileLoop(n *models.WhileLoop) (interface{}, error) {
	for {
		// Evaluate the condition.
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

		// Execute the loop body.
		for _, stmt := range n.Body {
			_, err := e.Execute(stmt)
			if err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

// isValidOperator checks if the given operator is a valid arithmetic operator.
// It returns true if the operator is valid, and false otherwise.
func (e *Executor) isValidOperator(operator string) bool {
	return operator == "+" || operator == "-" || operator == "*" || operator == "/"
}
