package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Chris-Mwiti/ChatX/logger"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
)

func main (){
	//logger initialization
	_, sugarLogger := logger.Init()

	//now lets try to reconfigure the node ip address
	//@todo: Research more on how to customize a node settings
	node, err := libp2p.New()
	if err != nil {
		panic(err)
	}
	hostInfo := host.InfoFromHost(node)

	sugarLogger.Infow("Host info: ", 
		"Host Id", hostInfo.ID,
		"Host Addr", hostInfo.Addrs,
	)

	//lets wait for a shutdown signal to close the connection
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	sugarLogger.Infoln("Received signal, shutting down...")

	//shut down the node when an error has occured
	if err := node.Close(); err != nil {
		panic(err)
	}
}
