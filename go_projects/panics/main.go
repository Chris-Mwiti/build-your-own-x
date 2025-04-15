package main

import (
	"fmt"
)

//an example of a recover panic function
func recoverPanicFunc(){
	defer fmt.Println("defer position 1");
	defer func(){
		fmt.Println("defer position 2")

		if e := recover(); e != nil {
			fmt.Println("Recover with the following error: ", e)
		}
	}()

	//execute a panic func
	panic("Just panicking for the sake of example")
	
	//this is a func that will never be executed
	fmt.Println("Hello am a statement that will never be executed")
}

func main(){
	fmt.Println("Start panicing execution")
	recoverPanicFunc()
	fmt.Println("The main program regains control after the first panic")
}