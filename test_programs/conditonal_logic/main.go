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
			// temperature = 50
			&models.Assignment{
				Variable: &models.Variable{Name: "temperature"},
				Value:    &models.Number{Value: 50},
			},
			// if temperature > 100 { state = "Gas" } else if temperature > 0 { state = "Liquid" } else { state = "Solid" }
			&models.IfStatement{
				Condition: &models.ComparisonExpression{
					Operator: ">",
					Left:     &models.Variable{Name: "temperature"},
					Right:    &models.Number{Value: 100},
				},
				Consequent: &models.Assignment{
					Variable: &models.Variable{Name: "state"},
					Value:    &models.String{Value: "Gas"},
				},
				Alternate: &models.IfStatement{
					Condition: &models.ComparisonExpression{
						Operator: ">",
						Left:     &models.Variable{Name: "temperature"},
						Right:    &models.Number{Value: 0},
					},
					Consequent: &models.Assignment{
						Variable: &models.Variable{Name: "state"},
						Value:    &models.String{Value: "Liquid"},
					},
					Alternate: &models.Assignment{
						Variable: &models.Variable{Name: "state"},
						Value:    &models.String{Value: "Solid"},
					},
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
	stateValue := exec.Env()["state"]
	fmt.Printf("State = %v\n", stateValue)
}
