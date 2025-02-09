package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)


func RequestServer(){

	fmt.Println("Request server is up and running")
	const serverPort = 3338
	//creation of a context
	actualCtx := context.Background()
	//creation of new mux handler
	mux := http.NewServeMux()
	
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		reqCtx := r.Context()
		fmt.Printf("server: %s \n : port = %s", r.Method, reqCtx.Value("keyServerAddr"))
		fmt.Fprintf(w, `{"message": "hello!"}`)
	})

	server := http.Server{
		Addr: ":3338",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context{
			reqCtx := context.WithValue(actualCtx, "keyServerAddr",l.Addr().String())
			return reqCtx
		},
	}



	//creation of a go routine to server a server
	go func(){
		err := server.ListenAndServe()

		if err != nil {
			if errors.Is(err, http.ErrServerClosed){
				fmt.Printf("error running http server: %s\n", err)
			}
		}

	}()

	time.Sleep(100 * time.Millisecond)

	//creation of client instance
	requestURL := fmt.Sprintf("http://localhost:%d", serverPort)
	//making an actual request to the url
	//res, err := http.Get(requestURL)

	//making  a request using http.NewRequest to have more control over the request
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		//exit the system
		os.Exit(1)
	}

	//creation of a respond client
	//the http.DefaultClient.Do is used to call a request from a predefined request 
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)


	//read all the response body
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)



	//POST REQUEST
	
}