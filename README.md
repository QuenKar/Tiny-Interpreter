
# Little Interpreter

This program is for learning how to write a interpreter in Go.

### Process
source code -> Token -> AST(Abstract Syntax Tree) -> Evaluation.

### Interpreter Structure
#### Lexer
make source code become Tokens.

#### Parser
use recursive descent parser, which can make Tokens from Lexer become AST.

#### Evaluator
use tree walking interpreter.
A tree walking interpreter that recursively evaluates an AST is probably the slowest of all approaches, but easy to build, extend.

### Implement Functions
- variable bindings
- integers and booleans
- arithmetic expressions
- built-in functions(len, echo...)
- first-class and higher-order functions
- closures
- a string data structure

### Some Examples
```
echo("hello world!");

let age = 1;
let name = "quenkar";
let result = 10 * (20 / 2);

let sayhello = fn(x){
	echo("hello" + x);
};

let add = fn(a,b){return a+b;};

let newAdder = fn(x) {
	fn(y) { x + y };
};
let addTwo = newAdder(2);

```

#### TODO
- Array
- Hash Table