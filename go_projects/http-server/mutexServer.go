package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
)

//creation of a mutex server
func MutexServer(){
	mux := http.NewServeMux()
	//addition of mux route handlers
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)

	//creation of multiple server using the http.Server interface method

	//creation of a context that will server in the background 
	actualCtx := context.Background()

	//create a new context with cancel functionality
	ctx, cancelCtx := context.WithCancel(actualCtx)

	//server one is a http
	serverOne := &http.Server{
		Addr: ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			//create a base context
			baseCtx := context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return baseCtx
		},
	}
	serverTwo := &http.Server{
		Addr: ":4444",
		Handler: mux,
		BaseContext: func(listener net.Listener) context.Context {
			//base context server addr
			baseCtx := context.WithValue(ctx, keyServerAddr, listener.Addr().String())
			return baseCtx
		},
	}

	//creation of a go routine to server the first server
	go func(){
		err := serverOne.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed){
			fmt.Printf("server has been closed\n")
		} else if err != nil{
			fmt.Printf("error closing server: %s\n", err)
			os.Exit(1)
		}
		cancelCtx()
	}()

	//creation of 2 go routine to serve the second server
	go func(){
		err := serverTwo.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed){
			fmt.Printf("the second server has been closed\n")
		} else if err != nil {
			fmt.Printf("erro closing the second server: %s\n", err)
		}
		cancelCtx()
	}()

	<-ctx.Done()
	
}
