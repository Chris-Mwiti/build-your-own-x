package main

import (
	"log"

	"github.com/Chris-Mwiti/build-your-own-x/go_projects/websockets/cmd"
	"github.com/joho/godotenv"
)

//initialize to load the env variables
func init(){
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("unable to retrieve the .env file content: %v", err)
	}
}

func main() {
	cmd.Execute()
}
