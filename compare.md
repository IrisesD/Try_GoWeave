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

* GoWeave's GET/SET pointcuts seem to just support the local/global variables/functions, ignoring the class-variables. In this way, we should implement the support for class-variables/methods.
For example:
in the belowing CUT:
```Go
package main

import (
	"fmt"
)

type A struct {
	num int
}

func (a A) f() {
	fmt.Println(a.num)
}

func main() {
	a := A{num: 10}
	a.f()
}
```
use
```
aspect {
  pointcut: get(a.num) or get(A.num) or get(num)
  imports (
    "fmt"
  )
  advice: {
	before: {
    	fmt.Println("before get num")
  	}
  }
}
```
then it will raise an error since this hasn't been implemented.

* GoWeave now doesn't support struct method name like
```Go
call(a.f())
```
instead, it can use
```Go
call(f())
```
However, this may not be what we want since it may instrument extra codes if there's two structs have the same function name called `f()`.

Also it doesn't support partial match method name like
```Go
call(b.*)
```

* GoWeave does not support annotations like AspectJ. But as the author of GoWeave said:
```
NO CODE MODIFICATIONS - my main use cases involve not modifying code so that is why we initially did not support annotations - I'm not opposed to adding these but that's not my intended goal
```
They are not intended to support annotations since it will need to modify the code.

* GoWeave is much slower than AspectJ.
Because of the fact that GoWeave does instrumentation on source code level, so it needs to do a lot of I/O (read and write file), so the speed of instrumentation is much slower than AspectJ, which could do byte-code level instrumentation.

However, I think this overhead cannot be avoid since it's source-code instrumentation's own disadvantage, we can optimize it but the overhead will still be here.