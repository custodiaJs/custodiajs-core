# VNH1 Programming Language

Welcome to `vnh1`, a concise and powerful programming language designed for efficient script execution and simplicity in mind. `vnh1` is geared towards developers looking for a straightforward syntax with powerful control flow and data manipulation capabilities.

## Key Features

- **Simplified Syntax**: Designed to be easy to read and write, eliminating unnecessary complexity.
- **Data Types**: Supports integers, floating-point numbers, strings, and boolean values.
- **Control Flow**: Includes `if` statements and `switch` cases for dynamic script flows.
- **Remote Block Calls (`rblockcall`)**: A unique feature for performing remote procedure calls within scripts.
- **Compilation to Bytecode**: `vnh1` scripts are compiled into an efficient bytecode, optimized for fast execution.

## Getting Started

Before diving into `vnh1`, ensure you have the necessary runtime environment and compiler installed on your system. Detailed installation instructions will be provided on the official `vnh1` website.

### Writing Your First Script

Create a simple `vnh1` script that demonstrates variable declaration, control flow, and a remote block call:

\```vnh1
// Variable declarations
greeting := "Hello, World!"
age := 30
pi := 3.14159
isHappy := true

// Basic if statement
if age > 18 {
  println(greeting)
}

// Remote block call
rblockcall("http://example.com/data", {}, <>) {
  // Handle response
}
catch(error) {
  // Error handling
}

// Switch statement
switch (age) {
  case 30:
    println("You are 30 years old.")
  default:
    println("Your age is unknown.")
}
\```

### Compilation and Execution

After writing your script, use the `vnh1` compiler to compile your script into bytecode:

\```sh
vnh1 compile myscript.vnh1
\```

This command generates a bytecode file that can be executed by the `vnh1` runtime:

\```sh
vnh1 run myscript.vnh1c
\```

## Documentation

For a comprehensive guide to `vnh1`'s syntax, data types, control structures, and more, refer to the [official documentation](#).

## Contributing

Contributions to the development and enhancement of `vnh1` are welcome. Please check out the [contribution guidelines](#) for more information.

## License

`vnh1` is open-source software licensed under the [MIT License](LICENSE).