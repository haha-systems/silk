package main

import (
	"fmt"

	"silk/internal/executor"
	"silk/internal/models"
)

func main() {
	// Construct AST
	program := &models.Program{
		Body: []models.Node{
			// x = 5
			&models.Assignment{
				Variable: &models.Variable{Name: "x"},
				Value:    &models.Number{Value: 5},
			},
			// y = 3
			&models.Assignment{
				Variable: &models.Variable{Name: "y"},
				Value:    &models.Number{Value: 3},
			},
			// z = (x + y) * 2
			&models.Assignment{
				Variable: &models.Variable{Name: "z"},
				Value: &models.BinaryExpression{
					Operator: "*",
					Left: &models.BinaryExpression{
						Operator: "+",
						Left:     &models.Variable{Name: "x"},
						Right:    &models.Variable{Name: "y"},
					},
					Right: &models.Number{Value: 2},
				},
			},
		},
	}

	// Execute
	exec := executor.NewExecutor()
	_, err := exec.Execute(program)
	if err != nil {
		fmt.Printf("Execution error: %v\n", err)
		return
	}

	// Output the result
	zValue := exec.Env()["z"]
	fmt.Printf("z = %v\n", zValue)
}
