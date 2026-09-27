# CDSL — Custom Description Scripting Language

CDSL (Custom Description Scripting Language) is a lightweight domain-specific language (DSL) for describing robotic movement and control instructions in a simple, human-readable format.

The CDSL interpreter is written in **Go** and provides an interactive terminal interface for defining values, evaluating expressions, controlling configured components, and introducing delays between operations.

## Features

* Terminal-based
* Variables and variable references
* Arithmetic expressions
* Configurable hardware components
* Component value-range validation
* Robotic movement instructions
* Command history
* Input validation and error reporting
* Serial Communication

## Build

### Requirements

Make sure the following tools are installed:

* Go 1.27
* Git
* Make

| Command          | Description                                                           |
| ---------------- | --------------------------------------------------------------------- |
| `make`           | Build CDSL.                                                           |
| `make check`     | Check that the required dependencies are installed.                   |
| `make clean`     | Remove the built `cdsl` binary.                                       |
| `make rebuild`   | Clean and build the project from scratch.                             |
| `make install`   | Install the `cdsl` binary to `~/.local/bin`.                          |
| `make uninstall` | Remove the installed `cdsl` binary.                                   |
| `make config`    | Create the default configuration file if it does not already exist.   |


## Usage

Run:
```sh
cdsl
```

The REPL starts:

```text
CDSL Interactive Mode
Type `exit` to quit.
Type `clear` to clear screen.

cdsl>  
```

Instructions can then be entered directly:

```text
set speed 50
print $speed
move BASE $speed
wait 1000
```

Command history is automatically stored in:

```text
~/.cdsl_history
```

## Language

CDSL programs are composed of instructions, with each instruction written on a separate line.

The current instruction set consists of:

* `set` - Create or update a variable
* `print` - Print a value or variable
* `move` - Move a component to a specified value
* `wait` - Pause execution for a specified duration
* `printc` - Print the loaded configuration

Instruction names are **case-sensitive** and should be written in lowercase.

---

## `set`

Creates a variable or changes the value of an existing variable.

### Syntax

```text
set <variable> <value>
```

### Examples

```text
set A 121
set speed 100
set angle 45.5
```

Values can be numeric literals, variables, or expressions.

Variables are referenced using `$`:

```text
set A 121
set B $A
```

Variables store floating-point values.

### Variable Names

Variable names must:

* Start with a letter or `_`
* Contain only letters, digits, and `_`

Valid:

```text
set speed 100
set _offset 10
set motor_angle 90
```

Invalid:

```text
set 123value 10
set motor-angle 90
set motor.angle 90
```

---

## `print`

Prints a value or the value of a variable.

### Syntax

```text
print <value>
```

### Examples

Print a variable:

```text
set A 121
print $A
```

Expressions can also be evaluated:

```text
print #[10+20*2]
```

An undefined variable results in an error:

```text
print $M
```

`set` and `print` are language conveniences for storing and inspecting values. They do not directly control hardware.

---

## `move`

Moves a configured component to a specified value.

### Syntax

```text
move <component> <value>
```

The component name identifies the hardware component being controlled.

### Examples

Move a component to a literal value:

```text
move JOINT_2 120
```

Use a variable:

```text
set A 121
move JOINT_1 $A
```

Use an expression:

```text
set offset 10
move JOINT_1 #[30+$offset]
```

The value must resolve to a valid numeric value.

Invalid values are rejected:

```text
move JOINT abc
move JOINT 12abc
move JOINT 12.0.1
```

An undefined variable also results in an error:

```text
move BASE $B
```

Once the variable exists, it can be used:

```text
set B 133
move BASE $B
```

The resulting value must also fall within the configured limits of the specified component.

---

## `wait`

Pauses execution for a specified number of milliseconds.

### Syntax

```text
wait <duration_ms>
```

### Examples

```text
wait 500
wait 1000
```

The duration can be stored in a variable:

```text
set delay 500
wait $delay
```

Expressions can also be used:

```text
wait #[250+250]
```

The duration must be a non-negative integer.

Invalid examples:

```text
wait -500
wait abc
```

---

## Variables

Variables are created using `set` and referenced using `$`.

```text
set speed 100
set angle 45.5

print $speed
move JOINT_1 $angle
wait $speed
```

Variables store floating-point values.

Using an undefined variable results in an error:

```text
move BASE $UNKNOWN
```

Variables can also be used inside expressions:

```text
set A 100
set B 50

print #[$A+$B]
```

---

## Expressions

Expressions can be used wherever a numeric value is accepted.

Expressions are enclosed in `#[...]`.

**Expressions must not contain whitespace.**

```text
#[10+20]     // Valid
#[10 + 20]   // Invalid
```

### Arithmetic Operators

| Operator | Operation      |
| -------- | -------------- |
| `+`      | Addition       |
| `-`      | Subtraction    |
| `*`      | Multiplication |
| `/`      | Division       |

### Examples

```text
set A 100
set B #[50+25]

print #[$A+$B]

move JOINT_1 #[45*2]

wait #[250+250]
```

Parentheses can be used to control evaluation order:

```text
print #[(10+20)*2]
```

Variables can be referenced inside expressions using `$`:

```text
set A 100
set B 50

print #[$A+$B]
```

Unary `+` and `-` are also supported:

```text
print #[-$A]
print #[+$B]
```

Division by zero and references to undefined variables result in an interpreter error.

---

## Components and Configuration

CDSL components are configured through:

```text
~/.config/cdsl/config.toml
```

Each component has an allowed value range.

For example:

```toml
[component.BASE]
min = 0
max = 180
```

A `move` instruction validates the requested value against the configured component limits before executing it.

The configuration directory and file are created automatically if they do not already exist.

Components are identified by their configured names and IDs.

---

## Communication

CDSL supports communication through a serial connection.

The hardware configuration specifies the serial port, baud rate, and configured components.

The communication protocol uses a fixed-size packet:

```text
+------------+-----------+-----------+-------------+---------+
| Start      | Packet ID | Component | Float Value | CRC     |
| 1 byte     | 2 bytes   | 1 byte    | 4 bytes     | 2 bytes |
+------------+-----------+-----------+-------------+---------+
```

The packet contains:

* **Start** — Packet start marker (`0xAA`)
* **Packet ID** — Identifier for the packet
* **Component** — Numeric component ID
* **Float Value** — Value associated with the component
* **CRC** — CRC-16 checksum

Supported baud rates include:

```text
9600
19200
38400
57600
115200
```

> Rates can be added or removed by altering the code.


The serial communication layer is implemented in Go and is responsible for constructing and transmitting packets.

---

## Complete Example

```text
// Configure values

set A 121
set speed 100
set delay 500

// Display values

print $A
print $speed

// Move components

move JOINT_1 $speed
move JOINT_2 120

// Wait between operations

wait $delay

// Use an expression

set speed #[50*2]

move JOINT_1 $speed

wait #[250+250]

print $speed
```

This demonstrates the basic CDSL workflow:

```text
set      → define values
print    → inspect values
move     → control components
wait     → introduce delays
#[...]   → calculate values
$...     → reference variables
//       → add comments
```

## Validation and Errors

CDSL validates instructions and their arguments before execution.

For example:

```text
print
set X
set @C 12
move JOINT abc
wait -500
move BASE $UNKNOWN
```

These demonstrate errors such as:

* Missing arguments
* Invalid variable names
* Invalid numeric values
* Negative wait durations
* Undefined variables
* Invalid expressions
* Values outside configured component limits

Values supplied to `move` are checked against the configured limits of the specified component.
