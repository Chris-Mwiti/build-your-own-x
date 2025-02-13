package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)


func main(){
	//instatiate the cmd struct
	var cmd = &cobra.Command{
		Use: "blockchain",
		Short: "The blockchain cli",
		Run: func(cmd *cobra.Command, args []string) {

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