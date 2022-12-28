<div align="center">
  <img src="https://benbusby.com/assets/images/y2k.png">

  [![MPL License](https://img.shields.io/github/license/benbusby/y2k)](LICENSE)
</div>

___

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

## Install

### Windows, macOS, Linux
Download from [the latest release](https://github.com/benbusby/y2k/releases)

### Go
`go install github.com/benbusby/y2k@latest`

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
y2k [args] <directory>

Args:
  -d int
    	Set # of digits to parse at a time (default 1)
  -debug
    	Enable to view interpreter steps in console
  ?
        Additional arguments are passed to programs as variables
        (see Fibonacci Sequence II for an example)
```

____

The directory provided should have a list of `*.y2k` files that have already
had their timestamps modified as needed. The `init.sh` script in each example
folder should give a good idea of how each program is created. The basic
approach is to create files using `touch` (or if you're on macOS, `gtouch` from
the `coreutils` package) to set their date to the desired timestamp.

From the "[modify and print variable](#modify-and-print-variable)" example:
```
touch --date=@812310071.235009210 1.y2k
```

For programs that require multiple timestamps, just increment the number used
in the `.y2k` filename, and continue your program on the following timestamp,
ignoring the first digit (the reason for this is explained in the [How It
Works](#how-it-works) section).

From the "[Fibonacci Sequence < 2000](#fibonacci-i)" example:
```
touch --date=@812108221.161214200 1.y2k
touch --date=@909219227.151272511 2.y2k
```

In this second example, the first digit of the `2.y2k` timestamp will be ignored,
which allows the number `2000` to span across the two file timestamps.

See [Examples](#examples) for more detail on how each example works.

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
is passed into the top level `interpreter.Parse()` function for initial
parsing. `Parse` will check the first digit of the timestamp to determine
which action to take:

<table>
  <tr>
    <th>Command ID</th>
    <th>Action</th>
    <th>Struct</th>
    <th>Parser</th>
  </tr>
  <tr>
    <td><code>9</code></td>
    <td>Print var or string</td>
    <td><code>Y2KPrint</code></td>
    <td><code>src/print.go</code> -> <code>ParsePrint</code></td>
  </tr>
  <tr>
    <td><code>8</code></td>
    <td>Create variable</td>
    <td><code>Y2KVar</code></td>
    <td><code>src/variable.go</code> -> <code>ParseVariable</code></td>
  </tr>
  <tr>
    <td><code>7</code></td>
    <td>Modify variable</td>
    <td><code>Y2KModify</code></td>
    <td><code>src/modify.go</code> -> <code>ParseModify</code></td>
  </tr>
  <tr>
    <td><code>6</code></td>
    <td>Evaluate condition</td>
    <td><code>Y2KCond</code></td>
    <td><code>src/condition.go</code> -> <code>ParseCondition</code></td>
  </tr>
  <tr>
    <td><code>5</code></td>
    <td>Change interpreter state</td>
    <td><code>Y2K</code></td>
    <td><code>src/interpreter.go</code> -> <code>ParseMeta (Parse)</code></td>
  </tr>
  <tr>
    <td><code>4</code></td>
    <td>Continue</td>
    <td>N/A</td>
    <td>N/A</td>
  </tr>
</table>

Each action is associated with its own struct, which holds values that are
pertinent to the action it needs to perform, and its own parsing function.
The next N-digits after the command digit are used to populate the struct's
public fields before using that struct to perform an action:

<table>
  <tr>
    <th>Struct (Command ID)</th>
    <th># Public Fields</th>
    <th>Field Descriptions</th>
    <th>Example</th>
  </tr>
  <tr>
    <td><code>Y2KPrint (9)</code></td>
    <td>1</td>
    <td>
      <ol>
        <li>Type</li>
        <ul>
          <li>1 --> Variable</li>
          <li>2 --> String</li>
        </ul>
      </ol>
    </td>
    <td>Print var 3: <code>[9]13</code></td>
  </tr>
  <tr>
    <td><code>Y2KVar (8)</code></td>
    <td>3</td>
    <td>
      <ol>
        <li>ID (numeric "name" of var)</li>
        <li>Type</li>
        <ul>
          <li>1 --> String</li>
          <li>2 --> Int</li>
          <li>3 --> Copy Var</li>
        </ul>
        <li>Size (# of digits/chars, or the var ID for Type 3)</li>
      </ol>
    </td>
    <td>Set var 3 to 5000: <code>[8]3245000</code></td>
  </tr>
  <tr>
    <td><code>Y2KMod (7)</code></td>
    <td>3</td>
    <td>
      <ol>
        <li>Var ID (ID of the variable to modify)</li>
        <li>Function</li>
        <ul>
          <li>1 --> <code>+=</code></li>
          <li>2 --> <code>-=</code></li>
          <li>3 --> <code>*=</code></li>
          <li>4 --> <code>/=</code></li>
          <li>5 --> <code>+= other var value</code></li>
          <li>6 --> <code>Copy other var value</code></li>
        </ul>
        <li>Size (# of digits/chars to use for modifying)</li>
      </ol>
    </td>
    <td>Var 3 /= 200: <code>[7]343200</code></td>
  </tr>
  <tr>
    <td><code>Y2KCond (6)</code></td>
    <td>4</td>
    <td>
      <ol>
        <li>Var ID (ID of the variable to use in the "left hand" side of the condition)</li>
        <li>Function</li>
        <ul>
          <li>1 --> <code>==</code></li>
          <li>2 --> <code><</code></li>
          <li>3 --> <code>></code></li>
          <li>4 --> <code>Is evenly divisible by</code></li>
        </ul>
        <li>Loop</li>
        <ul>
          <li>0 --> Treat as an <code>if</code></li>
          <li>1 --> Treat as a <code>while</code></li>
        </ul>
        <li>Size (# of digits/chars to use in the "right hand" side of the condition)</li>
      </ol>
    </td>
    <td>If var 3 is 25, print 'a': <code>[6]310225[9]21</code></td>
  </tr>
  <tr>
    <td><code>Y2K (5)</code> (same as top-level interpreter struct)</td>
    <td>2</td>
    <td>
      <ol>
        <li>Debug</li>
        <ul>
          <li>0 --> Debug mode off</li>
          <li>1 --> Debug mode on</li>
        </ul>
        <li>Digits (# of digits to parse per pass)</li>
      </ol>
    </td>
    <td>Change interpreter digit parsing size to 2: <code>[5]02</code></td>
  </tr>
</table>

For example, to create a variable, your timestamp would need to start with
`8` to initiate variable creation, and then the following 3 digits would be
used to set the variable's `ID/name`, `Type`, and `Size` attributes.

Once the struct's public fields have been set, it's passed over to its parser
function to actually perform the action. Continuing from the previous example
of creating a variable, this would mean A) populating the variable's value
until its size matches the `Size` attribute determined earlier, and B) storing
the variable in a key/value mapping of ID->Variable for future access.

So to set a variable "1" to the integer value 100, the timestamp would need
to have the following values:

```
  Begin creating variable
 /  Set variable ID to "1"
|  /  Set variable Type to int (2)
| |  /  Set variable Size to 3 digits
| | |  /  Insert 1
| | | |  /  Insert 0
| | | | |  /  Insert 0
| | | | | | /
8 1 2 3 1 0 0
```

Now that the variable has been set, you can reference it in other parts of
your program using its "1" ID.

In the following section, I've outlined some small example programs that
should help with understanding the language's current functionality.

## Examples

Each example has its own directory in the repo with a script to create the
files with the appropriate timestamp. The documentation here is to explain
how each example works.

### Set and Print Variable
`examples/set-print-var`

Timestamp(s):
- `812310092.100000000`

This expands on the example given in the "How It Works" section (setting
variable "1" to the value 100) by also printing the variable out to the
console after setting it.

```
  Begin creating variable
 /  Set variable ID to "1"
|  /  Set variable Type to int (2)
| |  /  Set variable Size to 3 digits
| | |  /  Insert 1
| | | |  /  Insert 0
| | | | |  /  Insert 0
| | | | | |  /  Begin print command
| | | | | | |  /  Set print type to var
| | | | | | | |  /  Print var "1"
| | | | | | | | |  /
8 1 2 3 1 0 0 9 2 1
```

Output: `100`

### Modify and Print Variable
`examples/modify-print-var`

Timestamp(s):
- `812310071.235009210`

This example takes an additional step after setting var "1" to 100 by then
subtracting 500 from that variable, and then printing var "1".

```
  Begin creating variable
 /  Set variable ID to "1"
|  /  Set variable Type to int (2)
| |  /  Set variable Size to 3 digits
| | |  /  Insert 1
| | | |  /  Insert 0
| | | | |  /  Insert 0
| | | | | |  /  Begin modify command
| | | | | | |  /  Target var "1" for modification
| | | | | | | |  /  Set modifier function to -=
| | | | | | | | |  /  Set modifier size to 3 digits
| | | | | | | | | |  /  Insert 5
| | | | | | | | | | |  /  Insert 0
| | | | | | | | | | | |  /   Insert 0
| | | | | | | | | | | | |  /   Begin print command
| | | | | | | | | | | | | |  /  Set print type to var
| | | | | | | | | | | | | | |  /  Print var "1"
| | | | | | | | | | | | | | | |  /
8 1 2 3 1 0 0 7 1 2 3 5 0 0 9 2 1
```

Output: `-400`

### Hello World
`examples/hello-world`

Timestamp(s):
- `502090134.051212150`
- `X04915181.204630000`

In this example, we're printing the string "Hello World!". Since character
codes are easier to encapsulate with 2-digit codes, we need to switch the
interpreter to 2-digit parsing mode at the very beginning.

As noted at the end of the explanation, print strings are terminated using
two space characters ("0") * N-digit parsing size. So for 2-digit parsing,
we'll need "00 00" to tell the interpreter to stop printing the string.

I've broken up the timestamps below into separate sections to make it a
little easier to read.

```
  Begin changing interpreter state
 /  Set debug mode to "off" (default)
|  /  Switch to 2-digit parsing size
| |  /
5 0 2

  Begin print command
 /   Set print type to string
|   /   Print "H"
|  |   /   Print "e"
|  |  |   /   Print "l"
|  |  |  |   /   Print "l"
|  |  |  |  |   /   Print "o"
|  |  |  |  |  |   /   Print " "
|  |  |  |  |  |  |   /
09 01 34 05 12 12 15 00

  Print "W"
 /   Print "o"
|   /   Print "r"
|  |   /   Print "l"
|  |  |   /   Print "d"
|  |  |  |   /   Print "!"
|  |  |  |  |   /
49 15 18 12 04 63

End print string
 / \
|   |
00 00
```

Output: `Hello World!`

### Fibonacci I
`examples/fibonacci-lt-2000`

Timestamp(s):
- `812108221.161214200`
- `X09219227.151272511`

For this first Fibonacci Sequence program, we're printing all values that are
less than 2000. We're also going to do something a bit hacky in order to
fit this solution into only 2 file timestamps. We're going to create two
variables that will hold two values of the sequence at a time (starting with 0
and 1), and add them to each other until the lower of the two values is above 2000.
This works since we know that an even number of terms is needed to reach `1597`,
the highest Fibonacci number before `2000`.

Like the "Hello World!" example, the timestamps below have been broken up to
make them easier to read. I've also grouped some commands into chunks of digits
if they've already been covered in previous examples.

**Note:** For a more robust Fibonacci Sequence implementation, see [Fibonacci
II](#fibonacci-ii).

```
81210 : Set var "1" to 0
82211 : Set var "2" to 1

  Begin conditional command
 /  Use var "1" in left hand of conditional
|  /  Set comparison to "<"
| |  /  Enable loop ("while" mode)
| | |  /  Set right hand value size to 4
| | | |  /  Insert 2
| | | | |  /  Insert 0
| | | | | |  /  Insert 0
| | | | | | |  /  Insert 0
| | | | | | | |  /
6 1 2 1 4 1 0 0 0

921 : Print var "1"
922 : Print var "2"

71512 : Var 1 += Var 2
72511 : Var 2 += Var 1
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
`examples/fibonacci-n-terms`

Timestamp(s):
- `812108221.183210693`
- `X11092173.611716127`
- `X25137921.100000000`

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

As before, previously explained commands will be grouped into logical chunks instead
of being explained digit-by-digit.

```
81210 : Set var "1" to 0
82211 : Set var "2" to 1
83210 : Set var "3" to 0

693110 : While Var 9 (cli arg) > 0
   921   : Print Var 1
   73611 : Var 3 = Var 1
   71612 : Var 1 = Var 2
   72513 : Var 2 += Var 3
   79211 : Var 9 -= 1
```

Output 1 (`y2k examples/fibonacci-n-terms 15`):

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

Output 2 (`y2k examples/fibonacci-n-terms 20`):

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
`examples/fizz-buzz`

Timestamp(s):
- `502080901.040609262`
- `X60808010.402212626`
- `X05000187.319775188`
- `X12106121.310071111`
- `X61402159.274200061`
- `X40139294.200061401`
- `X59284200.009210000`

The Fizz Buzz program highlights a few features that haven't been covered yet,
namely terminating and "continue"-ing conditionals. We also have to tell the
interpreter to switch between 1- and 2-digit parsing in order to create our
words "fizz" and "buzz" while maintaining the efficiency of 1-digit parsing.

The value `2000`, when used with a non-looped conditional, tells the
interpreter where the "body" of the statement needs to end. This is an
arbitrary value (although fitting given the name of the language) that is used
multiple times in this program to tell the interpreter where an "if" statement
ends. There's also the new parse-command `4` (aka `CONTINUE`), which returns an
empty string to the parent parser function instead of the remainder of the
timestamp. Since this is being used inside a "while" loop, this returns the
interpreter back to the beginning of the loop to reevaluate instead of
continuing to the next part of the timestamp.

```
502 : Change interpreter to 2-digit parsing mode

08 09 01 04 06 09 26 26 : Set var 9 to "fizz"
08 08 01 04 02 21 26 26 : Set var 8 to "buzz"

05 00 01 : Change interpreter back to 1-digit parsing mode

87319 : Copy var 9 to var 7 (var 7 = "fizz")
77518 : Append var 8 to var 7 (var 7 += "buzz")

81210 : Set var 1 to 0

61213100 : While var 1 < 100
    71111 : Var 1 += 1
    6140215 : If var 1 % 15 == 0
        927 : Print var 7 ("fizzbuzz")
        4 : Continue
    2000 : end-if
    614013 : If var 1 % 3 == 0
        929 : Print var 9 ("fizz")
        4 : Continue
    2000 : end-if
    614015 : If var 1 % 5 == 0
        928 : Print var 8 ("buzz")
        4 : Continue
    2000 : end-if
    921 : Print var 1
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

- **Why the pre-2000 timestamp limitation?**

The language was originally designed to just interpret timestamps of any
length, but both macOS and Go store Unix nanoseconds as an int64, which means
most timestamps can only reliably contain 18 digits without causing problems in
the interpreter. As a result, the maximum timestamp value that the interpreter
can reliably read falls around the year 2000¹.

This limitation -- along with the interpreter's original 2-digit-only parsing
design -- were both relevant to [the "Y2K
Problem"](https://en.wikipedia.org/wiki/Year_2000_problem), hence the name.

- **What does 0-byte actually mean? How can a program be 0 bytes?**

Since file content is not actually read by the interpreter, each `.y2k` file
can be completely empty (0 bytes) without affecting how each program is
interpreted. And since every file has to have a timestamp associated with it
anyway, there aren't any extra bytes needed to achieve this functionality.
Technically though, there's no such thing as a 0 byte file -- the metadata for
that file does have to be stored somewhere. But for code golfing purposes, it
should be counted as 0 bytes.

- **Why are there two ways to copy a variable's value to a new variable?**

The method through the `SET` command (`8`) inserts a new reference to a
variable using the specified ID, whereas the method through the `MODIFY`
command (`7`) updates the existing reference in the variable table. The former
can be useful for instantiating a new variable from an existing one, but can
cause problems if you're within the scope of a condition that has referenced
that variable.

For example:

```
# Loops infinitely, since the loop is initialized
# with a previous reference to Var 1
61213100 : While Var 1 < 100
    74111 : Var 4 += 1
    81314 : Set Var 1 to Var 4

# Loops as expected
61213100 : While var 1 < 100
    74111 : Var 4 += 1
    71514 : Copy var 4 value to var 1
```

_____

¹ Technically Sept. 2001, but close enough...
