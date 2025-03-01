package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)


func RequestServer(){

	fmt.Println("Request server is up and running")
	const serverPort = 3338
	//creation of a context
	actualCtx := context.Background()
	//creation of new mux handler
	mux := http.NewServeMux()

	//GET REQUEST MUX HANDLER
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		reqCtx := r.Context()
		fmt.Printf("server: %s \n : port = %s", r.Method, reqCtx.Value("keyServerAddr"))
		fmt.Fprintf(w, `{"message": "hello!"}`)
	})

	//POST REQUEST MUX HANDLER
	mux.HandleFunc("/post", func (w http.ResponseWriter, r *http.Request){
		//print the method type
		fmt.Printf("server: %s/\n", r.Method)
		fmt.Printf("server: query id: %s\n", r.URL.Query().Get("id"))
		fmt.Printf("server: content-type: %s\n", r.Header.Get("content-type"))

		//headers
		fmt.Printf("server: headers:\n")
		for headerName, headerValue := range r.Header {
			fmt.Printf("\t%s = %s\n", headerName, strings.Join(headerValue, ", "))
		}

		//read all the content type in the body request
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("server: could not read the body content type")
		}
		fmt.Printf("server: request body: %s\n", reqBody)

		fmt.Fprintf(w, `{"message": "hello!"}`)

		//timeout simulation
		time.Sleep(35 * time.Second)
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

	jsonBody := []byte(`{"client_message": "hello, server!"}`)
	bodyReader := bytes.NewReader(jsonBody)
	//creation of client instance
	requestURL := fmt.Sprintf("http://localhost:%d/post?id=1234", serverPort)
	//making an actual request to the url
	//res, err := http.Get(requestURL)

	//making  a request using http.NewRequest to have more control over the request
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		//exit the system
		os.Exit(1)
	}

	//customization of the request headers
	req.Header.Set("Content-Type", "application/json")


	//creation of customizable response client
	client := http.Client{
		//add a timeout to prevent response hanging
		Timeout: 30 * time.Second,
	}
	//creation of a respond client
	//the http.DefaultClient.Do is used to call a request from a predefined request 
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)


	//read all the response bodyhttps://github.com/gitau-BSC/Project_final.git
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)



	//POST REQUEST
}