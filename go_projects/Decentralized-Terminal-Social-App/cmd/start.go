package cmd

import (
	"github.com/Chris-Mwiti/ChatX/logger"
	"github.com/Chris-Mwiti/ChatX/nodes"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
 Use: "start",
 Short: "Use this command to start a node and check its info details such as id and address",
 Long: `
	This command is used to start a node and check its info details such as id and address.
	Example: chatX start [!no args]
	To be modified in the future
	`,
 Run: func(cmd *cobra.Command, args []string){
		_, log := logger.Init()
		log.Infoln("Start command called.")
		nodes.Start(log)
	},
}


func init() {
	rootCmd.AddCommand(startCmd)
}
