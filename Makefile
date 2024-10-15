# Makefile for Silk test programs

.PHONY: all build run benchmark clean

all: build run

build:
	@echo "Building all test programs..."
	@go build -o bin/basic_arithmetic test_programs/basic_arithmetic/main.go
	@go build -o bin/functions test_programs/functions/main.go
	@go build -o bin/loops test_programs/loops/main.go
	@go build -o bin/parallelism test_programs/parallelism/main.go
	@go build -o bin/conditional_logic test_programs/conditional_logic/main.go

run: build
	@echo "Running basic arithmetic test..."
	@./bin/basic_arithmetic
	@echo "Running conditional logic test..."
	@./bin/conditional_logic
	@echo "Running loops test..."
	@./bin/loops
	@echo "Running parallelism test..."
	@./bin/parallelism

benchmark: build
	@echo "Benchmarking basic arithmetic..."
	@time ./bin/basic_arithmetic
	@echo "Benchmarking conditional logic..."
	@time ./bin/conditional_logic
	@echo "Benchmarking loops..."
	@time ./bin/loops
	@echo "Benchmarking parallelism..."
	@time ./bin/parallelism

clean:
	@echo "Cleaning up binaries..."
	@rm -rf bin
