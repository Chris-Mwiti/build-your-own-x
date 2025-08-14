/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Chris-Mwiti/ChatX/logger"
	"github.com/Chris-Mwiti/ChatX/nodes"
	"github.com/spf13/cobra"
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Used to start a p2p node and set it to actively wait connection activities",
	Long: `The command will spin a node host server and await incoming connections.
		Example (chatX listen [!no option/args])
	`,
	Run: func(cmd *cobra.Command, args []string) {
		_, log := logger.Init()
		log.Infoln("Listen command called.")
		nodes.Listen(log)
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
