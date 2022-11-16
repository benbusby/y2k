# Y2K

Contents
1. [Install](#install)
2. [Usage](#usage)
3. [Examples](#examples)
4. [FAQ](#faq)

## Install

### From Source
1. Install [Chicken Scheme](https://code.call-cc.org/)
2. Clone repo: `git clone https://git.sr.ht:~benbusby/y2k`
3. Build: `chicken-csc -o y2k interpreter.scm`

## Usage

```
y2k <directory>
```

**Note:** The directory provided should have a list of `*.y2k` files that have already
had their timestamps modified as needed. See [Examples](#examples) for more
detail.

## Examples

Each example has its own directory in the repo with a script to create the
files with the appropriate timestamp. The documentation here is to explain
how each example works.

### Hello World (`examples/hello_world`)

This example prints "hello world!".

The program starts with the `02` command, which begins a routine of converting
each following 2-digit code into a different character. Programs like this one
that need to span across multiple timestamps must have subsequent files begin
with a `01` command to tell the interpreter to pick up where it left off.
Otherwise, there would be instances where the file timestamp would need to have
leading 0s that would not actually be valid (i.e. a file with a timestamp of
"0001" would just be stored as "1", and the leading 0s would not be seen by the
interpreter).

- Timestamp 1 -- 208051212 -- 1976-08-04 18:00:12
```
  Begin Print
 /    Print "h"
|   /   Print "e"
|  |   /   Print "l"
|  |  |   /    Print "l"
|  |  |  |   /
|  |  |  |  |
02 08 05 12 12
```
- Timestamp 2 -- 115002315 - 1973-08-23 19:05:15
```
  Continue
 /   Print "o"
|   /   Print " "
|  |   /   Print "w"
|  |  |   /   Print "o"
|  |  |  |   /
|  |  |  |  |
01 15 00 23 15
```
- Timestamp 3 -- 118120427 -- 1973-09-28 21:13:47
```
  Continue
 /   Print "r"
|   /   Print "l"
|  |   /   Print "d"
|  |  |   /   Print "!"
|  |  |  |   /
|  |  |  |  |
01 18 12 04 27
```

Final result:
```
y2k examples/hello_world
hello world!
```

### Set and Print Integer Variable (`examples/set_print_var`)

This example sets a numeric variable "01" to the value `100` and prints it.

The program is slightly more complex, but shorter than the "hello world" example.
It begins with a `04` command to begin creating a new variable. After this, a
unique set of steps is invoked to tell the interpreter what kind of variable it
is, how big it is, etc. Refer to the [Usage](#usage) section for a table
explaining each value.

- Timestamp 1 -- 401020310 -- 1982-09-16 04:31:50

```
  Begin Create Variable
 /    Set variable name/ID to "01"
|   /   Set variable type to numeric
|  |   /   Set variable size/digits to 3 [Begin single digit parsing]
|  |  |   /    Store 1
|  |  |  |   /  Store 0
|  |  |  |  |  /
04 01 02 03 1 0
```

- Timestamp 2 -- 100301 -- 1970-01-01 20:51:41
```
  Continue
 /    Store 0 [size of variable is now 3, switch back to 2-digit parsing]
|   /   Begin print variable
|  |   /   Print variable "01"
|  |  |   /
|  |  |  |
|  |  |  |
01 0  03 01
```

Final result:
```
y2k examples/set_print_var
100
```

## FAQ

- **Why the pre-2000 timestamp limit?**

The language was originally designed to just interpret timestamps in general,
but when testing it on macOS, I found that timestamps are stored using a signed
64 bit integer, so they can't reliably store any 10-digit timestamps without
limitation. Reducing the maximum number of digits to 9 means that timestamps
can't go past 999999999, which corresponds to ~Sept 2001. The date limitation
was an inspiration to cap the max timestamp to Jan 1 2000 build the language
around that limitation.

- **What does 0-byte actually mean?**

Since file content is not actually read by the interpreter, each `.y2k` file
can be completely empty (0 bytes) without affecting how each program is
written. Technically though, there's no such thing as a 0 byte file -- the
metadata for that file has to be stored somewhere. But for code golfing
purposes, it would be counted as 0 bytes.
