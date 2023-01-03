<div align="center">
  <img src="https://benbusby.com/assets/images/y2k.svg">

  [![MPL License](https://img.shields.io/github/license/benbusby/y2k)](LICENSE)
  [![builds.sr.ht status](https://builds.sr.ht/~benbusby/y2k.svg)](https://builds.sr.ht/~benbusby/y2k?)
  [![Go Report Card](https://goreportcard.com/badge/github.com/benbusby/y2k)](https://goreportcard.com/report/github.com/benbusby/y2k)
</div>

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
    4. [Fibonacci I: Values < 2000](#fibonacci-i)
    5. [Fibonacci II: N-terms](#fibonacci-ii)
    6. [Fizz Buzz](#fizz-buzz)
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
  - Supported types: `int`, `string`
- Variable modification
  - Supported operations: `+=`, `-=`, `/=`, `*=`
  - Also supported:
    - Adding one variable's value to another
    - Assigning one variable's value to another
- Conditional logic
  - Supported types: `if`, `while`
  - Supported comparisons: `==`, `>`, `<`, and divisibility (`% N == 0`)
- Print statements
  - Supported types: `var`, `string`
- Debug mode
  - Outputs where/how each timestamp digit is being parsed

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

You can then pass the new output directory with `y2k` to verify that the program
still functions the same, but with completely empty 0-byte files.

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

After the timestamps have been concatenated into one long string, this string is
passed into the top level `interpreter.Parse()` function, which will read in a
command ID to determine which action to take. The interpreter will then parse the
fields that pertain to that command, followed by the value (if applicable), before
returning to parse the next command ID.

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
- `812415009.210000000` :: `1995-09-29 16:50:09.210000000`

This expands on the example given in the "How It Works" section (setting
variable "1" to the value 100) by also printing the variable out to the
console after setting it.

```elixir
8124 # Create new variable 1 with type int (2) and size 4
1500 # Insert 4 digits (1500) into variable 1

921 # Print variable 1
```

Output: `1500`

### Modify and Print Variable
[`examples/modify-and-print-var.y2k`](examples/modify-and-print-var.y2k)

Timestamp(s):
- `812415007.123500921` :: `1995-09-29 16:50:07.123500921`

This example takes an additional step after setting var "1" to 1500 by then
subtracting 500 from that variable, and then printing var "1".

```elixir
8124 # Create new variable 1 with type int (2) and size 4
1500 # Insert 4 digits (1500) into variable 1

7123 # On variable 1, call function "-=" (2) with a 3-digit argument
500  # Insert 3 digits (500) into function argument

921  # Print variable 1
```

Output: `1000`

### Hello World
[`examples/hello-world.y2k`](examples/hello-world.y2k)

Timestamp(s):
- `502090134.051212150` :: `1985-11-28 22:28:54.051212150`
- `X04915181.204630000` :: `1998-09-04 07:19:41.204630000`

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

### Fibonacci I
[`examples/fibonacci-lt-2000`](examples/fibonacci-lt-2000.y2k)

Timestamp(s):
- `812108221.161214200` :: `1995-09-26 03:37:01.161214200`
- `X09219227.151272511` :: `1998-10-24 02:53:47.151272511`

For this first Fibonacci Sequence program, we're printing all values that are
less than 2000. We're going to do something a little hacky in order to fit this
solution into only 2 file timestamps. First we'll create two variables that
will hold two values of the sequence at a time (starting with 0 and 1), then
add them to each other until the lower of the two values is above 2000. This
works since we know that an even number of terms is needed to reach `1597`, the
highest Fibonacci number that is less than `2000`.

**Note:** For a more robust Fibonacci Sequence implementation, see [Fibonacci
II](#fibonacci-ii).

```elixir
8121 # Create variable 1 with type int (2) and size 1
0    # Insert 1 digit (0) into variable 1
8221 # Create variable 2 with type int (2) and size 1
1    # Insert 1 digit (1) into variable 2

# Init while loop (while var 1 < 2000)
61214 # Create conditional using variable 1, with comparison "<" (2),
      # as a loop (1), and with a right hand value size of 4
2000  # Insert 4 digits (2000) into conditional's right hand value

# Begin while loop
    921 # Print variable 1
    922 # Print variable 2
    71512 # var 1 += var 2
    72511 # var 2 += var 1
```

Output:

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
```

### Fibonacci II
[`examples/fibonacci-n-terms.y2k`](examples/fibonacci-n-terms.y2k)

Timestamp(s):
- `812108221.183210693` :: `1995-09-26 03:37:01.183210693`
- `X11092173.611716127` :: `1998-11-14 18:09:33.611716127`
- `X25137921.100000000` :: `1999-04-26 08:45:21.100000000`

For this modification to the Fibonacci Sequence program, we're now using an
argument from the command line as the number of terms to print. The logic
from the previous solution can't be used, since it always prints 2 values
at a time, so we need to update it. In the new program, we take the command
line argument and create a new loop that decrements that value until it reaches 0.

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
    73611 # var 3 = var 1
    71612 # var 1 = var 2
    72513 # var 2 += var 3
    79211 # var 9 -= 1
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
- `502080901.040609262` :: `1985-11-28 19:55:01.040609262`
- `X60808010.402212626` :: `1997-04-11 19:20:10.402212626`
- `X05000187.319775188` :: `1995-07-05 21:09:47.319775188`
- `X12106121.310071111` :: `1995-09-26 03:02:01.310071111`
- `X61402159.274200061` :: `1997-04-18 16:22:39.274200061`
- `X40139294.200061401` :: `1996-08-15 14:01:34.200061401`
- `X59284200.009210000` :: `1997-03-25 03:03:20.009210000`

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
8791 # Create variable 7 with type "copy" (9) and length 1 (variable ID length)
9    # Use 1 digit variable ID (9) to copy values from var 9 to var 7
7751 # On variable 7, call function "append var" (5) with a 1 digit argument
8    # Use 1 digit variable ID (8) to append values from var 8 to var 7

# Create variable 1 for iterating from 0 to 100
8121 # Create variable 1 with type int (2) and size 1
0    # Insert 1 digit (0) into variable 1

# Begin the loop from 0 to 100
61213100 # while variable 1 < 100
    71111 # var 1 += 1
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

## FAQ

- **Why the pre-2000 timestamp limitation? Why the name Y2K?**

The language was originally designed to interpret timestamps of any length, but
both macOS and Go store Unix nanoseconds as an int64. The max value of an int64
has 19 digits (`9223372036854775807`) but it wouldn't be reliable to write
programs using all 19 digits, since ostensibly there could be programs that
exceed this value fairly easily (a program to print the letter 'c' would start
with `923...`, for example). As a result, all timestamps for Y2K programs have
18 digits, which results in a maximum timestamp that falls around the year
2000¹.

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
    72111 # Var 2 += 1
    81912 # Overwrite Var 1 with Var 2 values

# GOOD
# Loops as expected, Var 1's value is updated on each
# iteration with Var 2's value
61213100 # While Var 1 < 100
    72111 # Var 2 += 1
    71512 # Copy Var 2 value to Var 1
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
