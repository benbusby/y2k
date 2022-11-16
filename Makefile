CC=chicken-csc

.PHONY: exe
exe:
	chicken-csc -o y2k interpreter.scm

static:
	chicken-csc -static -o y2k interpreter.scm

clean:
	rm -f y2k
