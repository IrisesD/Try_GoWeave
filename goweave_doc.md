## goweave docs

### features

The goweave now support 6 kinds of pointcuts:

* **call** 
```
These happen before, after or wrap around calling a method. The code is outside of the function.
```

* **execute**
```
These happen before or after executing a method. The code is put inside the method.
within:
```

* **within**
```
These happen for *every* statement within a function body declaration.
```

* **get**
```
These fire when a local/global variable has a get operation.
```

* **set**
```
These fire when a local/global variable has a set operation.
```

* **declaration**
```
This fires when a variable is declared.
```

### examples

#### call

CUT:
```Go
package main

import (
	"fmt"
)

func main() {
	a := 2
	fmt.Println(a)
}
```

instrumentation code:
```
aspect {
  pointcut: call(Println)
  imports (
    "fmt"
  )
  advice: {
	before: {
    	fmt.Println("before Println")
  	}
  }
}
```

ICUT:
```Go
package main

import (
"fmt"
)

func main() {
	a := 2
fmt.Println("before Println")
	fmt.Println(a)
}
```

#### execute

it now can just do instrumentation in user-defined functions.

CUT:
```Go
package main

import (
	"fmt"
)

func Print() {
	fmt.Println("Here.")
}

func main() {
	Print()
}
```

instrumentation code:
```
aspect {
  pointcut: execute(Print())
  imports (
    "fmt"
  )
  advice: {
	before: {
    	fmt.Println("before Print()")
  	}
  }
}
```

ICUT:
```Go
package main

import (
"fmt"
)

func Print() {
fmt.Println("before Print()")
	fmt.Println("Here.")
}

func main() {
	Print()
}
```

#### within

CUT:
```Go
package main

import (
	"fmt"
)

func Print() {
	fmt.Println("1")
	fmt.Println("2")
	fmt.Println("3")
}

func main() {
	Print()
}
```

instrumentation code:
```
aspect {
  pointcut: within(Print())
  imports (
    "fmt"
  )
  advice: {
	before: {
    	fmt.Println("before Print()")
  	}
  }
}
```

ICUT:
```Go
package main

import (
"fmt"
)

func Print() {
fmt.Println("before Print()")
	fmt.Println("1")
fmt.Println("before Print()")
	fmt.Println("2")
fmt.Println("before Print()")
	fmt.Println("3")
}

func main() {
	Print()
}
```

#### get

CUT:
```Go
package main

import (
	"fmt"
)

func main() {
	a := 2
	fmt.Println(a)
}
```

instrumentation code:
```
aspect {
  pointcut: get(a)
  imports (
    "fmt"
  )
  advice: {
	before: {
    	fmt.Println("before get a")
  	}
  }
}
```

ICUT:
```Go
package main

import (
"fmt"
)

func main() {
	a := 2
fmt.Println("before get a")
	fmt.Println(a)
}
```

#### set

CUT:
```Go
package main

import (
	"fmt"
)

func main() {
	var a int
    a = 2
	fmt.Println(a)
}
```

instrumentation code:
```
aspect {
  pointcut: set(a)
  imports (
    "fmt"
  )
  advice: {
	before: {
    	fmt.Println("before set a")
  	}
  }
}
```

ICUT:
```Go
package main

import (
"fmt"
)

func main() {
	var a int
fmt.Println("before set a")
	a = 2
	fmt.Println(a)
}
```

#### declaration

CUT:
```Go
package main

import "fmt"

func main() {
    ch := make(chan int, 2)
    ch <- 1
    ch <- 2
    fmt.Println(<-ch)
    fmt.Println(<-ch)
}
```

instrumentation code:
```
aspect {
  pointcut: declaration(a)
  imports (
    "fmt"
  )
  advice: {
	before: {
    	fmt.Println("before declaration of a")
  	}
  }
}
```

ICUT:
```Go
package main

import "fmt"

func main() {
fmt.Println("before declaration of ch")
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
```

### The basic logics of goweave

the search of statements is based on go/ast

goweave will use go/ast to parse the file into many kinds of statements.

For example, if we want to do instrumentation on CUT like below:
CUT:
```Go
package main

import (
	"fmt"
)

func main() {
	a := 2
	fmt.Println(a)
}
```

instrumentation code:
```
aspect {
  pointcut: get(a)
  imports (
    "fmt"
  )
  advice: {
	before: {
    	fmt.Println("before get a")
  	}
  }
}
```

goweave will trigger the below codes:
```Go
for _, decl := range file.Decls {
	fn, ok := decl.(*ast.FuncDecl)
	if !ok {
		continue
	}

	for x := 0; x < len(fn.Body.List); x++ {
		as, ok2 := fn.Body.List[x].(*ast.ExprStmt)
		if !ok2 {
			continue
		}

		blah, ok3 := as.X.(*ast.CallExpr)
		if !ok3 {
			continue
		}

		// can either be a unary || a ident (so far)
		fn2, ok4 := blah.Args[0].(*ast.UnaryExpr)
		if !ok4 {
			// look for ident

			fn3, ok5 := blah.Args[0].(*ast.Ident)
			if !ok5 {
				continue
			}

			if pk != fn3.Name {
				continue
			}

		} else {
			// look for channel

			blah2, ok4 := fn2.X.(*ast.Ident)
			if !ok4 {
				continue
			}

			if pk != blah2.Name {
				continue
			}
		}

		begin := fset.Position(as.Pos()).Line - 1
		after := fset.Position(as.End()).Line + 1

		before_advice := aspect.advize.before
		after_advice := aspect.advize.after

		if before_advice != "" {
			rout = w.writeAtLine(fname, begin+linecnt, before_advice)
			linecnt += strings.Count(before_advice, "\n") + 1
		}

		if after_advice != "" {
			rout = w.writeAtLine(fname, after+linecnt-1, after_advice)

			linecnt += strings.Count(after_advice, "\n") + 1
		}

	}
}
```
First it uses `fn, ok := decl.(*ast.FuncDecl)` to get the body of function declaration.
Then in the loop of `for x := 0; x < len(fn.Body.List); x++`, goweave searches for all kinds of statements that will cause a GET action, comparing them to the identification name defined in the specification file(.weave).
If the name matches, goweave uses`begin := fset.Position(as.Pos()).Line - 1` and `after := fset.Position(as.End()).Line + 1` to get where to instrument.
Then it will instrument the advice to that position.
This is how goweave works.

### Usage
first clone the goweave repo:
```shell
git clone https://github.com/deferpanic/goweave.git
```
then revise the main.go to
```Go
package main

import (
	"log"
	"./weave"
)

const (
	version = "v0.1"
)

// main is the main point of entry for running goweave
func main() {
	log.Println("goweave " + version)

	w := weave.NewWeave()
	w.Run()

}
```
use `go build -o gow` to build goweave and then move the executable `gow` to the project we want to instrument.

then we can use `./gow` instead of `go build` in the project.

**notes**:
the goweave will delete the ICUT after build the project, if we want to see what the code looks like after instrumentation, we can enter weave.go, there's a function named `Run`:
```Go
func (w *Weave) Run() {
	w.prep()
	w.loadAspects()

	// old regex parsing
	// only used for go routines currently
	// soon to be DEPRECATED
	w.transform()

	// applys around advice && evals execution joinpoints
	filepath.Walk(w.buildLocation, w.VisitFile)

	w.build()
}
```
if we comment the `w.build()`, then the project will not be built and the ICUT code can be maintained.