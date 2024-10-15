package main

import (
	"fmt"
	"math/rand"
	"time"

	"silk/internal/executor"
	"silk/internal/models"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Define the 'compute' function
	computeFunc := &models.FunctionDeclaration{
		Name: "compute",
		Parameters: []*models.Variable{
			{Name: "n"},
		},
		Body: []models.Node{
			// Simulate computation with a sleep
			&models.FunctionCall{
				Name: "sleepRandom",
				Args: []models.Node{},
			},
			// Print the result
			&models.FunctionCall{
				Name: "printResult",
				Args: []models.Node{
					&models.Variable{Name: "n"},
				},
			},
		},
	}

	// List of numbers to process
	numbers := []float64{1, 2, 3, 4, 5}
	var functionCalls []models.Node

	for _, num := range numbers {
		num := num // Capture the loop variable
		functionCalls = append(functionCalls, &models.FunctionCall{
			Name: "compute",
			Args: []models.Node{&models.Number{Value: num}},
		})
	}

	// Create the main program AST
	program := &models.Program{
		Body: []models.Node{
			computeFunc, // Include the function definition
			&models.ParallelBlock{
				Body: functionCalls, // Execute function calls in parallel
			},
		},
	}

	// Create the executor
	exec := executor.NewExecutor()

	// Register built-in functions
	exec.RegisterBuiltin("sleepRandom", func(args []interface{}) (interface{}, error) {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return nil, nil
	})

	exec.RegisterBuiltin("printResult", func(args []interface{}) (interface{}, error) {
		n := args[0]
		fmt.Printf("Processed number: %v\n", n)
		return nil, nil
	})

	// Execute the program
	_, err := exec.Execute(program)
	if err != nil {
		fmt.Printf("Execution error: %v\n", err)
		return
	}
}
