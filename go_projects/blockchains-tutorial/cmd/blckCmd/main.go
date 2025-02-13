package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)


func main(){
	//instatiate the cmd struct
	var cmd = &cobra.Command{
		Use: "blckCmd",
		Short: "The blockchain cli",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running a blockchain command")
		},
	}

	//add the version commnad
	cmd.AddCommand(versionCmd)

	err := cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}