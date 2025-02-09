package main

import (
	// "errors"
	"fmt"
	"io"
	"net/http"
)

const ( keyServerAddr= "serverAddr" )

//creation of a get root handler
func getRoot(res http.ResponseWriter, req *http.Request){
	//get the context for the current reqeust made
	ctx := req.Context()

	//check whether the url has a query
	hasFirst := req.URL.Query().Has("first")
	first := req.URL.Query().Get("first")
	fmt.Printf("The following is the value of first:%s", first)
	hasSecond := req.URL.Query().Has("second")
	second := req.URL.Query().Get("second")


	//extraction of body query
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
	}


	fmt.Printf("%s : got / request. first(%t)=%s, second(%t)=%s\n, body:\n%s\n", ctx.Value(keyServerAddr), hasFirst, first, hasSecond, second, body)
	io.WriteString(res, "Welcome to the home page of my website")
}

func getHello(res http.ResponseWriter, req *http.Request){
	ctx  := req.Context()
	fmt.Printf("%s: got /hello request\n", ctx.Value(keyServerAddr))

	//extraction of body values
	myName := req.PostFormValue("myName")
	if myName == ""{
		res.Header().Set("x-missing-field", "myName")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	io.WriteString(res, fmt.Sprintf("Hello. %s!\n", myName))
}

func main(){
	//handler func .for the "/" route
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

	//3. Spinning up a request server
	RequestServer()	


	//2. Spinning up the mutex server
	MutexServer()

}
