/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Chris-Mwiti/ChatX/logger"
	"github.com/Chris-Mwiti/ChatX/nodes"
	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
//@todo: Add a pre run command that will start the node independently
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Used to ping another peer connection and check if connection can be established",
	Long: `Ping is a subcommand used to send a ping protocal connection to another peer,
	checks if the connection is alive and if not establishes one.

	Example: (chatX ping [...addresses])
	`,
	Run: func(cmd *cobra.Command, args []string) {
		_, log := logger.Init()
		log.Infoln("Ping command called.")

		if len(args) <= 0 {
			log.Fatal("Ping command must have atleast one argument(address)")
		}
	
		err := nodes.Ping(log, args)
		if err != nil {
			log.Fatal("Error while pinging connection")
		}
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
