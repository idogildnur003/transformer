# Data Transformer Library

## Objective
A generic Golang library for transforming and processing user data into a structured format.

## Requirements
- Go 1.20.5+

## Setup
```shell
git clone git@github.com:idogildnur003/transformer.git
cd transformer
go mod tidy
mkdir -p data/output  # Ensure the default output directory exists
```

**Note:** The `data/output` directory must exist before running the program, as it serves as the default output path.

## Usage

### Running the CLI
You can execute the `read` command with required flags:

#### **Basic Usage**

```sh
./data-transformer read --input path/to/input.json
```

#### **Flags & Options**

| Flag       | Short | Description                                      | Default Value                         |
| ---------- | ----- | ------------------------------------------------ | ------------------------------------- |
| `--input`  | `-i`  | Path to the input file (Required)                | None (Must be provided)               |
| `--rules`  | `-r`  | Path to the transformation rules file (Optional) | `configs/default_mapping_config.json` |
| `--output` | `-o`  | Path to the output directory (Optional)          | `data/output/`                        |


If no output file is specified, the program will save the transformed data to `data/output` by default.
If no rules file is provided, the program will use a default one.

### Examples:

#### Example 1: Using a custom rules file
```shell
go run cmd/main.go read --input data/input/users.json --output data/output/transformed.json --rules configs/custom_rules.json
```

#### Example 2: Using the default rules file
```shell
go run cmd/main.go read --input data/input/
```

#### Example 3: Using the default rules file and single input file 
```shell
go run cmd/main.go read --input data/input/fake_users_1.json
```

#### Example 4: Running as a built executable
```shell
go build -o transformer cmd/main.go
./transformer read -i data/input/users.json -o data/output/transformed.json -r configs/custom_rules.json
```

## How It Works
1. **Unmarshalling**: The input file is read and converted into structured data.
2. **Transformation**: The data is processed based on predefined rules.
3. **Storage**: The transformed data is saved into structured output files.
4. **Parallel Processing**: Uses goroutines to process files efficiently.

## Reasoning Behind the Storage of Users and Sign-In Activities
I have structured the storage of **Users** and **Sign-In Activities** into two separate files, ensuring a consistent schema for each:
## Data Storage Design

### **Users File Structure**

By default, user data is stored in a structured JSON format:

```json
{
    "id": "id",
    "external_id": "string",
    "mail": "string",
    "type": "string",
    "location": "string",
    "is_enabled": "boolean",
    "first_name": "string",
    "last_name": "string"
}
```

### **Sign-In Activity File Structure**

Sign-in activities are stored separately in an optimized format for scalability:

```json
[
    {
        "requestId": "string",
        "timeStamp": "string",
        "type": "string",
        "userId": "string"
    }
]
```

### **Rationale for This Design**

1. **Optimized Querying & Database Integration**

   - Users and sign-in activities are stored separately to avoid redundant data storage.
   - The `userId` in sign-in activity records serves as a **foreign key**, ensuring efficient database relations.

2. **Scalability & Performance**

   - Storing sign-in activities separately makes it easier to manage large user datasets.
   - The `type` field allows flexibility in categorizing different types of sign-in activities.

3. **Future-Proofing & Extensibility**

   - This structure allows easy expansion without modifying existing user records.
   - Additional metadata for sign-ins can be added without altering the main user structure.

---

## Running Tests

The project includes unit tests for all core components:

```sh
cd transformer/pkg/[module you want to test]
go test
```

Ensure all dependencies are installed before running tests.
