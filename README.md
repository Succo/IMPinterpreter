# IMP is a test language described by simple semantic 

- arithmetical expression
```IMP
a ::=n
    |a '+' a
    |a '*' a
    |a '-' a
```
- boolean expression
```IMP
a ::=true
    |false
    |a = a
    |a <= a
    |!b
    |b||b
    |b&&b
```
- instruction 
```IMP
c ::=skip
    |X ::= a
    |c ; c
    |if b then c else c
    |while b do c
```

This is a go program to interpret and execute program written in IMP.
It only takes one argument, the name of a source file.
