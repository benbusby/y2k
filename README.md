[![Y2K Logo](https://benbusby.com/assets/images/y2k-logo.jpeg)](https://benbusby.com/assets/images/y2k-logo.jpeg)

[![MPL License](https://img.shields.io/github/license/benbusby/y2k)](LICENSE)
![GitHub release](https://img.shields.io/github/v/release/benbusby/y2k)
[![builds.sr.ht status](https://builds.sr.ht/~benbusby/y2k.svg)](https://builds.sr.ht/~benbusby/y2k?)
[![Go Report Card](https://goreportcard.com/badge/github.com/benbusby/y2k)](https://goreportcard.com/report/github.com/benbusby/y2k)

___

<table>
    <tr>
        <td><a href="https://sr.ht/~benbusby/y2k">SourceHut</a></td>
        <td><a href="https://github.com/benbusby/y2k">GitHub</a></td>
    </tr>
</table>

Contents
1. [Install](#install)
2. [Features](#features)
3. [Usage](#usage)
4. [How It Works](#how-it-works)
5. [Examples](#examples)
    1. [Set and Print Variable](#set-and-print-variable)
    2. [Modify Variable](#modify-and-print-variable)
    3. [Print "Hello World!"](#hello-world)
    4. [Area of a Circle](#area-of-a-circle)
    5. [Fibonacci Sequence (N-terms)](#fibonacci-sequence)
    6. [Fizz Buzz](#fizz-buzz)
    7. [Count Up Forever (Golf Hack)](#count-up-forever)
6. [FAQ](#faq)
    1. [Why the pre-2000 timestamp limitation? Why the name Y2K?](#faq)
    2. [What does 0-byte actually mean? How can a program be 0 bytes?](#faq)
    3. [Why are there two ways to copy a variable's value to a new variable?](#faq)
    4. [How would I show proof of my solution in a code golf submission?](#faq)
    5. [Why doesn't Y2K have X feature?](#faq)
7. [Contributing](#contributing)

## Install

### Binary (Windows, macOS, Linux)
Download from [the latest release](https://github.com/benbusby/y2k/releases)

### Go
`go install github.com/benbusby/y2k@latest`

### From Source

1. Install Go: https://go.dev/doc/install
2. Clone and build project:
```
git clone https://github.com/benbusby/y2k.git
cd y2k
go build
```

## Features

- Variable creation
  - Supported types: `int`, `float`, `string`
- Variable modification
  - Supported operations: `+=`, `-=`, `/=`, `*=`, `**= (exponentiation)`, `= (overwrite)`
  - Accepts primitive types (`int`, `float`, `string`) or variable IDs as arguments
- Conditional logic
  - Supported types: `if`, `while`
  - Supported comparisons: `==`, `>`, `<`, and divisibility (`% N == 0`)
- Print statements
  - Supported types: `var`, `string`
- Debug mode
  - Outputs where/how each timestamp digit is being parsed
- "Raw" file reading/writing
  - Allows writing Y2K programs as file content (see [Examples](#examples)) and
    exporting to a set of new 0-byte files with their timestamps modified,
    rather than manually editing individual file timestamps.

## Usage

```
y2k [args] <input>

Args:
  -d int
        Set # of digits to parse at a time (default 1)
  -debug
        Enable to view interpreter steps in console
  -export
        Export a Y2K raw file to a set of timestamp-only files
  -outdir string
        Set the output directory for timestamp-only files when exporting a raw Y2K file.
        This directory will be created if it does not exist. (default "./y2k-out")
```

____

***Note:** See [CHEATSHEET.md](CHEATSHEET.md) for help with writing Y2K commands.*

The simple way to write Y2K programs is to write all commands to a file as
regular file content first.

For example, from [the "Set and Print Variable" program](#set-and-print-variable):

```elixir
# set-and-print-var.y2k
8124 # Create new variable 1 with type int (2) and size 4
1500 # Insert 4 digits (1500) into variable 1

921  # Print variable 1
```

```shell
$ y2k set-and-print-var.y2k
1500
```

You can then export this file to a set of empty 0-byte files (or in this
example, just one file) with their timestamps modified to achieve the same
functionality as the raw file:

```shell
$ y2k -export set-and-print-var.y2k
Writing ./y2k-out/0.y2k -- 812415009210000000 (1995-09-29 16:50:09.21 -0600 MDT)

$ ls ./y2k-out/*.y2k -lo --time-style="+%s%9N"
-rw-r--r-- 1 benbusby 0 812415009210000000 ./y2k-out/0.y2k
```

Then you could pass the new output directory as input to `y2k`, and verify that
the program still functions the same, but with completely empty 0-byte files.

```shell
$ y2k ./y2k-out
1500
```

See [Examples](#examples) for more detailed breakdowns of current example programs.

## How It Works

To preface, Y2K is obviously a fairly unconventional language. Since everything
is interpreted using numbers, it can perhaps be a bit confusing at first to get
a feel for how to write programs. If you find any of the below documentation
confusing, please let me know!

Y2K works by reading all files in a specified directory (sorted numerically)
and extracting each of their unix nanosecond timestamps. It then concatenates
each timestamp, stripping the first digit off of each timestamp except for the
first one. This is done to eliminate the potential issue of a command spanning
across multiple file timestamps where a 0 might be required at the beginning of
the timestamp. For example, if the number 1000 was being written to a variable
and the 0s needed to be at the beginning of the next file timestamp, this would
only be possible if the timestamp was prefixed with a non-zero digit (otherwise
leading 0s are ignored).

After the timestamps have been concatenated into one long string, this string
is passed into the top level `interpreter.Parse()` function, which will
interpret the first digit as a command ID in order to determine which action to
take. Command IDs are mapped to fields that are unique to that particular
command, and the interpreter will use the next N-digits to parse out values for
each of those fields. Some commands, such as setting and modifying variables,
have a "Size" field which tells the interpreter how many digits following the
command fields will be used to store/use a specific value. For instance, if you
wanted to store the number 100 in a variable, you would use the "Create
Variable" command ID, and the "Size" field for that command would be 3. The
following 3 digits of the timestamp would be "100", and the interpreter would
then read and store that 3-digit value in the variable.

Once the interpreter finishes reading the command ID, the command fields, and
any subsequent N-digit values (if applicable), it returns to the beginning to
parse the next command ID.

[CHEATSHEET.md](CHEATSHEET.md) contains a simplified breakdown of command IDs,
command fields, and when values are needed for the different commands. Please also
refer to the following [Examples](#examples) section for simple programs that help to
inform how the Y2K interpreter works.

## Examples

Each example below is taken directly from the [`examples`](examples) folder,
but with added explanation for how/why they work.

All examples can be exported to 0-byte solutions using the `-export` flag if
desired.

### Set and Print Variable
[`examples/set-and-print-var.y2k`](examples/set-and-print-var.y2k)

Timestamp(s):
- `812415009210000000 (1995-09-29 16:50:09.210000000)`

This expands on the example given in the "How It Works" section (setting
variable "1" to the value 100) by also printing the variable out to the
console after setting it.

```elixir
8124 # Create new variable 1 with type int (2) and size 4
1500 # Insert 4 digits (1500) into variable 1

921  # Print variable 1
```

Output: `1500`

### Modify and Print Variable
[`examples/modify-and-print-var.y2k`](examples/modify-and-print-var.y2k)

Timestamp(s):
- `812310071203500921 (1995-09-28 11:41:11.203500921)`

This example is very similar to the previous example, only this time we're
going to modify the variable after setting it. In this case, we set variable 1
to the int value 100, then subtract 500 from that variable.

```elixir
8123  # Create new variable 1 with type int (2) and size 3
100   # Insert 3 digits (100) into variable 1

71203 # On variable 1, call function "-=" (2) with a primitive (0) 3-digit argument
500   # Insert 3 digits (500) into function argument

921   # Print variable 1
```

Output: `-400`

### Hello World
[`examples/hello-world.y2k`](examples/hello-world.y2k)

Timestamp(s):
- `502090134051212150 (1985-11-28 22:28:54.05121215)`
- `804915181204630000 (1995-07-04 21:33:01.20463000)`

In this example, we're printing the string "Hello World!". Since character
codes are easier to encapsulate with 2-digit codes, we need to switch the
interpreter to 2-digit parsing mode at the very beginning.

As seen at the end of the explanation below, print strings are terminated using
two space ("0") characters * N-digit parsing size. So for 2-digit parsing,
we'll need "00 00" to tell the interpreter to stop printing the string.

```elixir
502 # Switch interpreter to 2-digit parsing size

09 01 # Begin print string command

34 05 12 12 15 00 # Print "Hello "
49 15 18 12 04 63 # Print "World!"

00 00 # End print command
```

Output: `Hello World!`

### Area of a Circle
[`examples/area-of-circle.y2k`](examples/area-of-circle.y2k)

Timestamp(s):
- `813913141592679501 (1995-10-17 00:59:01.592679501)`
- `827131199210000000 (1996-03-17 23:39:59.210000000)`

In this example, we're introducing a couple of new concepts. One is the ability
to include variables from the command line, and the other is modifying one
variable using another variable's value.

To include variable's from the command line, we simply pass the value after the
input. For example, `y2k my-program.y2k 10` would include a variable with the
value `10` that we can access from the beginning of the program. Since most Y2K
programs create variables using sequential IDs (i.e. 0 -> 1 -> 2, etc),
variables added from the command line are added to the back of the variable
map, with descending IDs from there. So if you're running Y2K in the default
1-digit parser mode, command line arguments are added as variables with IDs
starting at 9, then 8, and so on. As an example: `y2k my-program.y2k foo bar`
would have variable 9 set to "foo" and variable 8 set to "bar".

The other new concept is modifying a variable with the value from another
variable. In previous examples, we've used primitive types for arguments, but
in this case we need to multiply our "Pi" variable (1) by our squared radius.
To do this, we set the third field to "1" to tell the interpreter that the
value we're passing in is a variable ID, not a primitive type.

```elixir
8139      # Set variable 1 to type float (3) and size 9

131415926 # Insert 9 digits (131415926) into variable 1, using the first
          # digit (1) as the decimal placement (3.1415926)

79501     # Modify variable 9 (CLI arg) using the "**=" function (5),
          # with a non-variable (0) argument size of 1
2         # Use the number 2 as the function argument (var9 **= 2)

71311     # Modify variable 1 using the "*=" function (3), with a
          # variable argument (1) with a variable ID size of 1
9         # Use the variable ID 9 in the function argument (var1 *= var9)

921       # Print variable 1
```

Output (`y2k examples/area-of-circle.y2k 10`):

```
314.15926
```

Output (`y2k examples/area-of-circle.y2k 25`):

```
1963.495375
```

### Fibonacci Sequence
[`examples/fibonacci-n-terms.y2k`](examples/fibonacci-n-terms.y2k)

Timestamp(s):
- `812108221183210693 (1995-09-26 03:37:01.183210693)`
- `811092173911171911 (1995-09-14 09:22:53.911171911)`
- `827211137920110000 (1996-03-18 21:52:17.920110000)`

For this modification to the Fibonacci Sequence program, we're now using an
argument from the command line as the number of terms to print. In this new
program, we'll take the command line argument and create a new loop that
decrements that value until it reaches 0.

We need 3 variables for this program, not including the variable added from the
command line: a variable for the "current" value, a "placeholder" variable to
track the "current" value before it gets updated, and a "next" variable to track
the "next" value in the sequence. On each loop iteration, we 1) print "current", 2)
set "placeholder" to "current", 3) set "current" to "next", 4) add "placeholder"
to "next", and 5) decrement counter.

```elixir
8121 # Create variable 1 with type int (2) and size 1
0    # Insert 1 digit (0) into variable 1
8221 # Create variable 2 with type int (2) and size 1
1    # Insert 1 digit (1) into variable 2
8321 # Create variable 3 with type int (2) and size 1
0    # Insert 1 digit (0) into variable 3

# Init while loop (while var 9 > 0)
69311 # Create conditional using variable 9, with comparison ">" (3),
      # as a loop (1), and with a right hand value size of 1
0     # Insert 1 digit (0) into conditional's right hand value

# Begin while loop
    921 # Print var 9
    739111 # var 3 = var 1
    719112 # var 1 = var 2
    721113 # var 2 += var 3
    792011 # var 9 -= 1
```

Output 1 (`y2k examples/fibonacci-n-terms.y2k 15`):

```
0
1
1
2
3
5
8
13
21
34
55
89
144
233
377
```

Output 2 (`y2k examples/fibonacci-n-terms.y2k 20`):

```
0
1
1
2
3
5
8
13
21
34
55
89
144
233
377
610
987
1597
2584
4181
```

### Fizz Buzz
[`examples/fizz-buzz.y2k`](examples/fizz-buzz.y2k)

Timestamp(s):
- `502080901040609262 (1985-11-28 19:55:01.040609262)`
- `860808010402212626 (1997-04-11 19:20:10.402212626)`
- `805000187919771118 (1995-07-05 21:09:47.919771118)`
- `881210612131007110 (1997-12-03 21:43:32.131007110)`
- `811614021592742000 (1995-09-20 10:20:21.592742000)`
- `861401392942000614 (1997-04-18 16:09:52.942000614)`
- `801592842000921000 (1995-05-27 10:40:42.000921000)`

The Fizz Buzz program highlights a few features that haven't been covered yet,
namely terminating and "continue"-ing conditionals. We also have to tell the
interpreter to switch between 1- and 2-digit parsing in order to create our
words "fizz" and "buzz" while maintaining the efficiency of 1-digit parsing.

The value `2000`, when used within a non-looped conditional, tells the
interpreter where the "body" of the statement needs to end. This is an
arbitrarily chosen value (but fits with the name of the language) that is used
multiple times in this program to tell the interpreter where an "if" statement
ends. There's also the new command ID `4` (aka `CONTINUE`), which returns an
empty string to the parent parser function instead of the remainder of the
timestamp. Since this is being used inside a "while" loop, this returns the
interpreter back to the beginning of the loop to reevaluate instead of
continuing to the next part of the timestamp.


```elixir
502 # Change interpreter to 2-digit parsing mode

# Set variables 9 and 8 to "fizz" and "buzz" respectively
08 09 01 04 # Create variable 9 with type string (1) and length 4
06 09 26 26 # Insert 4 chars ("fizz") into variable 9
08 08 01 04 # Create variable 8 with type string (1) and length 4
02 21 26 26 # Insert 4 chars ("buzz") into variable 8

05 00 01 # Change interpreter back to 1-digit parsing mode

# Set variable 7 to "fizzbuzz"
8791  # Create variable 7 with type "copy" (9) and length 1 (variable ID length)
9     # Use 1 digit variable ID (9) to copy values from var 9 to var 7
77111 # On variable 7, call function "+=" (5) using a variable (1) with a 1 digit ID
8     # Use 1 digit variable ID (8) to append values from var 8 to var 7

# Create variable 1 for iterating from 0 to 100
8121 # Create variable 1 with type int (2) and size 1
0    # Insert 1 digit (0) into variable 1

# Begin the loop from 0 to 100
61213100 # while variable 1 < 100
    711011 # var 1 += 1
    6140215 # if var 1 % 15 == 0
        927 # print var 7 ("fizzbuzz")
          4 # continue
    2000 # end-if
    614013 # if var 1 % 3 == 0
       929 # print var 9 ("fizz")
         4 # continue
    2000 # end-if
    614015 # if var 1 % 5 == 0
        928 # print var 8 ("buzz")
          4 # continue
    2000 # end-if
    921 # print var 1
```

Output:
```
1
2
fizz
4
buzz
fizz
7
8
fizz
buzz
11
fizz
13
14
fizzbuzz
16
17
fizz
19
buzz
fizz
22
23
fizz
buzz
26
fizz
28
29
fizzbuzz
...<continued>...
```

### Count Up Forever
[`examples/count-up-forever.y2k`](examples/count-up-forever.y2k)

Timestamp(s):
- `611110721011922000 (1989-05-13 18:58:41.011922000)`

*Originally from [this problem on
codegolf.stackexchange.com](https://codegolf.stackexchange.com/questions/63834/count-up-forever/)*

This program highlights a "hacky" feature that is included in Y2K, which is the ability
to create new variables by referencing their IDs before they've been created. In this
example, we create variable 1 through its reference in the while loop, and variable 2
the first time that we try to modify it. When you do this, an "empty" variable is created
without a specific type and a numeric value of 0.

Creating variables this way isn't necessarily recommended, since it makes
programs more difficult to read and can only be used for creating variables
with a value of 0, but it can be a useful way to condense a solution into an
even smaller footprint. In this case, we can fit the solution to the problem in
a single file timestamp (and in raw format is only 15 bytes after comments and
newlines are removed).

```
611110 # while var 1 == 0
    721011 # var 2 += 1
    922    # print var 2
```

Output:
```
1
2
3
4
5
...<continued until killed>...
```

## FAQ

- **Why the pre-2000 timestamp limitation? Why the name Y2K?**

The language was originally designed to interpret timestamps of any length, but
both macOS and Go store Unix nanoseconds as an int64. The max value of an int64
has 19 digits (`9223372036854775807`) but it wouldn't be reliable to write
programs using all 19 digits, since there can be programs that exceed this
value fairly easily (a program to print the letter 'c' would start with
`923...`, for example). As a result, all timestamps for Y2K programs have 18
digits, which results in a maximum timestamp that falls around the year 2000¹.

The interpreter was also originally designed to only ever read 2 digits at a time.
These combined limitations reminded me of [the "Y2K
Problem"](https://en.wikipedia.org/wiki/Year_2000_problem), hence the name.

- **What does 0-byte actually mean? How can a program be 0 bytes?**

Since the interpreter only reads file *timestamps* and not file *content*, each
`.y2k` file can be completely empty (0 bytes) without affecting how each
program is interpreted. And since every file has to have a timestamp associated
with it anyway, there aren't any extra bytes needed to achieve this
functionality. Technically though, there's no such thing as a 0 byte file --
the metadata for that file does have to be stored somewhere. But for code
golfing purposes, I believe it would be counted as 0 bytes.

- **Why are there two ways to copy a variable's value to a new variable?**

The method through the `SET` command (`8`) inserts a new reference to a
variable using the specified ID, whereas the method through the `MODIFY`
command (`7`) updates the existing reference in the variable table. The former
can be useful for instantiating a new variable from an existing one, but can
cause problems if you're within the scope of a condition that has referenced
that variable.

For example:

```elixir
81210 # int var 1 = 0
82210 # int var 2 = 0

# BAD
# Loops infinitely, since the reference to Var 1 that
# was used to create the loop is overwritten, and the
# value of the original reference is never updated
61213100 # While Var 1 < 100
    721101 # Var 2 += 1
    81912  # Overwrite Var 1 with Var 2 values

# GOOD
# Loops as expected, Var 1's value is updated on each
# iteration with Var 2's value
61213100 # While Var 1 < 100
    721101 # Var 2 += 1
    719112 # Copy Var 2 value to Var 1
```

- **How would I show proof of my solution in a code golf submission?**

I'm not sure the best way to do this yet. Assuming you wrote your solution
as a "raw" Y2K file, you can run `y2k -export my-program.y2k`, and then
run the following command:

```shell
$ ls ./y2k-out/*.y2k -lo --time-style="+%s%9N"
-rw-r--r-- 1 benbusby 0 502090134051212150 0.y2k
-rw-r--r-- 1 benbusby 0 104915181204630000 1.y2k
```

You could also include your raw Y2K file contents along with the 0-byte
proof, to be extra thorough.

- **Why doesn't Y2K have X feature?**

The language is still in development. Feel free to open an issue, or refer to
the [Contributing](#contributing) section if you'd like to help out!

_____

¹ Technically Sept. 2001, but close enough...

## Contributing

I would appreciate any input/contributions from anyone. Y2K still needs a lot
of work, so feel free to submit a PR for a new feature, browse the issues tab
to see if there's anything that you're interested in working on, or add a new
example program.

The main thing that would help is trying to solve current or past code-golfing
problems from https://codegolf.stackexchange.com. If there's a limitation in
Y2K (there are definitely a ton) that prevents you from solving the problem,
open an issue or PR so that it can be addressed!
