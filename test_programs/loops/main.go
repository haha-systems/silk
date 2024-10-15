package main

import (
	"fmt"
	"silk/internal/executor"
	"silk/internal/models"
)

func main() {
	// Create a ForLoop node to test the executor
	forLoop := &models.ForLoop{
		Initialization: &models.Assignment{
			Variable: &models.Variable{Name: "i"},
			Value:    &models.Number{Value: 0},
		},
		Condition: &models.ComparisonExpression{
			Operator: "<",
			Left:     &models.Variable{Name: "i"},
			Right:    &models.Number{Value: 5},
		},
		Post: &models.Assignment{
			Variable: &models.Variable{Name: "i"},
			Value: &models.BinaryExpression{
				Operator: "+",
				Left:     &models.Variable{Name: "i"},
				Right:    &models.Number{Value: 1},
			},
		},
		Body: []models.Node{
			&models.FunctionCall{
				Name: "print",
				Args: []models.Node{&models.Variable{Name: "i"}},
			},
		},
	}

	// Create the main program AST
	program := &models.Program{
		Body: []models.Node{
			forLoop,
		},
	}

	// Create the executor
	exec := executor.NewExecutor()

	// Register built-in print function
	exec.RegisterBuiltin("print", func(args []interface{}) (interface{}, error) {
		fmt.Println(args...)
		return nil, nil
	})

	// Execute the program
	_, err := exec.Execute(program)
	if err != nil {
		fmt.Printf("Execution error: %v\n", err)
	}
}
