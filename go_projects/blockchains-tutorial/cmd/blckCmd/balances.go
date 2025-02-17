package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

//loading the latest DB state and printing it to the standard output
func balancesListCmd() *cobra.Command{
	var balancesListCmd = &cobra.Command{
		Use: "list",
		Short: "List all balances",
		Run: func(cmd *cobra.Command, args []string){
			//fetch the balances from the persisted state
		 
		},
	}
	return balancesListCmd
}

func incorrectUsageErr(error string){
	fmt.Printf("Incorrect usage: %s",error)
}

func balancesCmd() *cobra.Command {
	//creation of the command
	var balancesCmd = &cobra.Command{
		Use: "balances",
		Short: "Interact with balances (list...)",
		PreRunE: func(cmd *cobra.Command, args []string){
			return incorrectUsageErr() 
		},
		Run: func(cmd *cobra.Command, args []string){
			
		},
	}

	balancesCmd.AddCommand(balancesListCmd())
	return balancesCmd
}