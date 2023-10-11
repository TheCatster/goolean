# Goolean

Goolean is a command-line application written in Go, designed to simplify Boolean algebra expressions and generate truth tables. Utilizing the Cobra framework, it functions as a REPL (Read-Evaluate-Print Loop) to provide an interactive user experience.

## Features

- [ ] Simplification of Boolean algebra expressions.
- [ ] Generation of truth tables.
- [x] Support for the following operators: AND (`&`), OR (`|`), NOT (`!`), NAND, NOR, and XOR (`⊕`).
- [x] Parentheses support for expression grouping.
- [x] Interactive REPL interface for real-time expression evaluation.

## Installation

To install Goolean, ensure you have Go installed on your machine, then run the following command:

```bash
go get -u github.com/thecatster/goolean
```

## Usage

Once installed, run the application using the following command:

```bash
goolean
```

In the REPL, enter your Boolean algebra expressions, and Goolean will simplify them and print out the corresponding truth tables. Type `exit` to quit the REPL.

```goolean
goolean> a & b
```

## Development

To contribute to Goolean, clone the repository and create a new branch for your feature or bug fix. After making your changes, open a pull request to merge them into the main branch.

```bash
git clone https://github.com/thecatster/goolean.git
cd goolean
git checkout -b feature/your-feature
```

## Testing

`WIP`: Goolean comes with a suite of unit tests to ensure its core functionality works as expected. Run the tests using the following command:

```bash
go test ./...
```

## License

Goolean is licensed under the GPLv3 License. See the [LICENSE](LICENSE) file for details.
