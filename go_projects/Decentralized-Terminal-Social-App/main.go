package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Chris-Mwiti/ChatX/logger"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/multiformats/go-multiaddr"
)

func main (){
	//logger initialization
	_, sugarLogger := logger.Init()

	//now lets try to reconfigure the node ip address
	//@todo: Research more on how to customize a node settings

	//here we will adjust the ping service of the node
	node, err := libp2p.New(
		libp2p.Ping(false),
	)
	if err != nil {
		panic(err)
	}
	hostInfo := host.InfoFromHost(node)

	sugarLogger.Infow("Host info: ", 
		"Host Id", hostInfo.ID,
		"Host Addr", hostInfo.Addrs,
	)

	//configuration of our own ping service
	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID,pingService.PingHandler)

	//here lets generate the p2p address format for conn
	peerInfo := peerstore.AddrInfo{
		ID: node.ID(),
		Addrs: node.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	sugarLogger.Infoln("node address formated: ", addrs[0])

	

	//here we will try to add a commandline arg to capture the node addr to be connected to
	if len(os.Args) > 1 {
		addr, err := multiaddr.NewMultiaddr(os.Args[1])
		if err != nil {
			sugarLogger.Panic(err)
		}

		//create a peer AddrInfo from the addr string captured
		peer, err := peerstore.AddrInfoFromP2pAddr(addr)
		if err != nil {
			sugarLogger.Panic(err)
		}

		//@todo: Add context timeout to provide timeoutes for conn establishement
		if err := node.Connect(context.Background(), *peer); err != nil {
			sugarLogger.Panic(err)
		}
		sugarLogger.Infoln("pinging 5 messages to: ", addr)
		//@todo: context timeout to provide timeouts for conn establishment
		ch := pingService.Ping(context.Background(), peer.ID)
		for i := 0; i < 5; i++ {
			res := <-ch
			sugarLogger.Infoln("Got response!, round trip: ", res.RTT)
		}

	} else {
		//lets wait for a shutdown signal to close the connection
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		sugarLogger.Infoln("Received signal, shutting down...")
	}

	

	//shut down the node when an error has occured
	if err := node.Close(); err != nil {
		panic(err)
	}
}
