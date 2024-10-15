// test_programs/functions/main.go

package main

import (
	"fmt"
	"silk/internal/executor"
	"silk/internal/models"
)

func main() {
	exec := executor.NewExecutor()

	// Define a simple function that adds two numbers.
	function := &models.FunctionDeclaration{
		Name: "add",
		Parameters: []*models.Variable{
			{Name: "a"},
			{Name: "b"},
		},
		Body: []models.Node{
			&models.ReturnStatement{
				Value: &models.BinaryExpression{
					Left:     &models.Variable{Name: "a"},
					Right:    &models.Variable{Name: "b"},
					Operator: "+",
				},
			},
		},
	}

	// Register the function.
	exec.RegisterFunction(function.Name, function)

	// Call the function.
	functionCall := &models.FunctionCall{
		Name: "add",
		Args: []models.Node{
			&models.Number{Value: 3},
			&models.Number{Value: 5},
		},
	}

	// Execute the function call.
	result, err := exec.Execute(functionCall)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Result: %v\n", result)
}
