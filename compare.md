## Compare GoWeave with AspectJ

**it's hard to make a table to compare them since I don't know all the features they have, so I just list some features below**

### some features that AspectJ has but GoWeave don't have:

* GoWeave cannot identify all the statements because it didn't implement all. For example, in declaration statements they use:
```Go
if len(as.Rhs[0].(*ast.CallExpr).Args) == 2 {
	_, k := (as.Rhs[0].(*ast.CallExpr).Args[0]).(*ast.ChanType)
	if !k {
	    continue
	}

	r2, k2 := (as.Rhs[0].(*ast.CallExpr).Args[1]).(*ast.BasicLit)
	if !k2 {
		continue
	}
...
}
```
to get the declaration of variables, which means they just consider the definition of channel, since channeld's definition is like:
```Go
ch := make(chan int, 2)
```

So if we want it to be more complete, we must implement other kinds of form of statements like:
```Go
var a int
```
or 
```Go
a := 2
```

* GoWeave now can just **execute&within** support the instrumentation of user-defined function, if we want to do instrumentation on library function, it may not work. For example:
if we want to do instrumentation to the belowing CUT:
```Go
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello World!")
}
```
using:
```
aspect {
  pointcut: execute(Println)
  imports (
    "fmt"
  )
  advice: {
	before: {
    	fmt.Println("in Println")
  	}
  }
}
```
then the output will not change:
```shell
(base) ➜  main git:(master) ✗ ./main
Hello World!
```

However, I think this feature is not so important, since in real scenario, we are not usually intended to revise the library function, since it may raise a lot of mistakes. The GoWeave can support the **call** pointcut of library functions, of which I think is enough.

