package makingrequest

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

const serverPort = 3338

func RequestServer(){
		//creation of a context
		actualCtx := context.Background()

		baseCtx, cancelCtx := context.WithCancel(actualCtx)

	//creation of a go routine to server a server
	go func(){
		//creation of new mux handler
		mux := http.NewServeMux()
		
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
			reqCtx := r.Context()
			fmt.Printf("server: %s \n : port = %s", r.Method, reqCtx.Value("keyServerAddr"))
		})

		server := http.Server{
			Addr: fmt.Sprintf(":%d", serverPort),
			Handler: mux,
			BaseContext: func(l net.Listener) context.Context{
				reqCtx := context.WithValue(baseCtx, "keyServerAddr",l.Addr().String())
				return reqCtx
			},
		}

		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed){
				fmt.Printf("error running http server: %s\n", err)
			}
		}

		cancelCtx()
	}()

	time.Sleep(100 * time.Millisecond)

	//creation of client instance
	requestURL := fmt.Sprintf("http://localhost:%d", serverPort)
	//making an actual request to the url
	res, err := http.Get(requestURL)

	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		//exit the system
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	<- baseCtx.Done()
}