# Test Programs for Silk Executor

This directory contains various test programs designed to validate the functionality of the **Silk Executor**. Each program demonstrates specific features of my language and execution engine, allowing us to confirm that different parts of the Silk system are working correctly. Below is a description of each test program and what it aims to achieve.

## Directory Structure

```
./test_programs/
├── README.md
├── basic_arithmetic
│   └── main.go
├── conditional_logic
│   └── main.go
├── loops
│   └── main.go
└── parallelism
    └── main.go
```

### 1. `basic_arithmetic/main.go`

This program tests **basic arithmetic operations** such as addition, subtraction, multiplication, and division. It ensures that the Silk Executor can correctly perform numeric calculations and handle different types of binary expressions.

- **Purpose**: Verify the implementation of basic arithmetic functionality.
- **Expected Output**: The results of arithmetic operations are printed, allowing us to see if the Executor handles numerical expressions correctly.

### 2. `conditional_logic/main.go`

This program tests **conditional logic**, using `if` statements to demonstrate the Executor's ability to evaluate conditions and execute different branches accordingly. It helps confirm that boolean expressions and conditionals are functioning as intended.

- **Purpose**: Validate the Executor's handling of conditional logic, including true/false evaluations.
- **Expected Output**: Outputs based on condition evaluations, depending on the specific values given to the conditions.

### 3. `loops/main.go`

This program tests **loop functionality**, specifically focusing on `ForLoop` constructs. It iterates over a set of values and prints each value, ensuring that the loop is executed the correct number of times, and that loop control (initialization, condition, and post-iteration) works properly.

- **Purpose**: Confirm that the Executor can correctly handle iterative constructs like `for` loops.
- **Expected Output**: The numbers from 0 to 4 are printed, demonstrating the correct execution of the loop.

### 4. `parallelism/main.go`

This program tests **parallel execution** by executing a set of tasks concurrently. It demonstrates the Silk Executor's ability to handle concurrent execution and synchronize goroutines effectively.

- **Purpose**: Test the Executor's ability to perform multiple operations in parallel, ensuring that goroutines are managed properly.
- **Expected Output**: Outputs from concurrent tasks, which may be printed in a non-deterministic order, depending on the timing of each goroutine.

## How to Run the Programs

To run each program, navigate to the corresponding directory and use the following command:

```sh
go run main.go
```

This will execute the test program, and you can observe the output in the console to verify that the functionality works as expected.

## Contribution

Feel free to add more test programs to this directory to expand the coverage of Silk's features. Each program should aim to test a specific feature of the language or executor, helping us ensure robustness and correctness.
