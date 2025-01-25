package main

//interface -> context.Context
//QueryContext -> sql.DB

import (
	"fmt"
	"context"
)


func doSomething(ctx context.Context) {
	fmt.Printf("doSomething: myKey's value is %s\n", ctx.Value("myKey"))

	anotherCtx := context.WithValue(ctx, "myKey", "anotherValue")
	doAnother(anotherCtx)

	fmt.Printf("doSomething: myKey's value is %s\n", ctx.Value("myKey"))
}

func doAnother(ctx context.Context) {
	fmt.Printf("doAnother: myKey's value is %s\n", ctx.Value("myKey"))
}

//Advantages of context:
//1. ability to access data stored inside a context

func main() {
	//one of the two ways of creating a context
	//used as a placeholder when you're not sure which context to use
	ctx := context.TODO()

	//second way of creating contexts
	//used to start a known context
	ctx = context.Background()
	ctx = context.WithValue(ctx, "myKey", "myValue")
	doSomething(ctx)

	//example of ending a context using withCancel
	ending()
}
