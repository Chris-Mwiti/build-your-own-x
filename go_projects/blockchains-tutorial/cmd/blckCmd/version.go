package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

const Major = "0"
const Minor = "1"
const Fix = "0"

const Verbal = "TX Add & Balances List"

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "Describe version",
	Run: func(cmd *cobra.Command, args []string){
		fmt.Printf("Verision %s.%s.%s-beta %s", Major, Minor, Fix, Verbal)
	},
}