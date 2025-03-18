# Data Transformer Library

## Objective
A generic Golang library for transforming and processing user data into a structured format.

## Requirements
- Go 1.20.5+

## Setup
```shell
git clone [your repo]
cd data-transformer
go mod tidy
```

## Usage
This CLI tool processes user data based on transformation rules.

### Running the CLI
You can execute the `read` command with required flags:

```shell
go run cmd/main.go read --input <input_file> --output <output_file> --rules <rules_file>
```

### Flags:
- `--input (-i)`: Path to the input JSON file (required).
- `--output (-o)`: Path to the output JSON file (required).
- `--rules (-r)`: Path to the transformation rules file (optional, defaults to `configs/default_rules.json`).

If no rules file is provided, the program will use a default one.

## How It Works
1. **Unmarshalling**: The input file is read and converted into structured data.
2. **Transformation**: The data is processed based on predefined rules.
3. **Storage**: The transformed data is saved into structured output files.
4. **Parallel Processing**: Uses goroutines to process files efficiently.

## Running Tests
Run the following command to execute unit tests:
```shell
go test ./...
```
