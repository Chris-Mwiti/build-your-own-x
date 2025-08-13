package main

import (
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/Chris-Mwiti/ChatX/logger"
)

func main (){
	//logger initialization
	_, sugarLogger := logger.Init()
	node, err := libp2p.New()
	if err != nil {
		panic(err)
	}
	hostInfo := host.InfoFromHost(node)

	sugarLogger.Infow("Host info: ", 
		"Host Id", hostInfo.ID,
		"Host Addr", hostInfo.Addrs,
	)

	//shut down the node when an error has occured
	if err := node.Close(); err != nil {
		panic(err)
	}
}
