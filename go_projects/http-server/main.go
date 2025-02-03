package main

import (
//	"errors"
	"fmt"
	"io"
	"net/http"
	//"os"
)

const ( keyServerAddr= "serverAddr" )

//creation of a get root handler
func getRoot(res http.ResponseWriter, req *http.Request){
	//get the context for the current reqeust made
	ctx := req.Context()
	fmt.Printf("%s :got /request\n", ctx.Value(keyServerAddr))
	io.WriteString(res, "Welcome to the home page of my website")
}

func getHello(res http.ResponseWriter, req *http.Request){
	ctx  := req.Context()
	fmt.Printf("%s: got /hello request\n", ctx.Value(keyServerAddr))
	io.WriteString(res, "Hello, HTTP!\n")
}

func main(){
	//handler func for the "/" route
	// http.HandleFunc("/", getRoot)

	// //handler func for the "/hello" route
	// http.HandleFunc("/hello", getHello)

	// //setup and serve ther server
	// //default using default server multiplexer
	// //blocks the call untill the http.ListenAndServer finishes to call
	// //no handler makes it a single multiplexer server
	// err := http.ListenAndServe("127.0.0.1:3333", nil)

	// //error checking
	// // 1. Checks whether the error is of type server closed
	// if errors.Is(err, http.ErrServerClosed) {
	// 	fmt.Printf("server closed \n")
	// } else if err != nil {
	// 	fmt.Printf("error starting server: %s\n", err)
	// 	os.Exit(1)
	// }


	//2. Spinning up the mutex server
	MutexServer()
}
