package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

func balancesListCmd() *cobra.Command{
	var balancesListCmd = &cobra.Command{
		Use: "list",
		Short: "List all balances",
		Run: func(cmd *cobra.Command, args []string){
			//fetch the balances from the persisted state
		 
		},
	}
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
			fmt.Printf("Running the balances cmd")
		},
	}

	balancesCmd.AddCommand(balancesListCmd)
	return balancesCmd
}