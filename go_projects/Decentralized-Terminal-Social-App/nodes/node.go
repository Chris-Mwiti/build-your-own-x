package nodes

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/multiformats/go-multiaddr"
	"go.uber.org/zap"
)

func Start(logger *zap.SugaredLogger)(host.Host){
	logger.Infoln("Staring node...")
	node, err := libp2p.New()

	if err != nil {
		logger.Fatalf("Could not be able to start node: %v\n", err)
	}

	logger.Infoln("Node created!")

	return node
}

func Close(logger *zap.SugaredLogger, node host.Host){
	logger.Infoln("Shutting down node...")
	err := node.Close()

	if err != nil {
		logger.Fatalf("Could not be able to shutdown node: %v\n", err)
	}

	logger.Infoln("Node shutdown!")
}



func Ping(logger *zap.SugaredLogger, peerAddr string)(error){
	node := Start(logger)
	defer Close(logger, node)
	//here lets generate the node info
	hostInfo := peerstore.AddrInfo{
		ID: node.ID(),
		Addrs: node.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(&hostInfo)
	if err != nil {
		logger.Errorf("Error while creating host addr from info: %s", err.Error())
		return err
	}
	logger.Infow("Host Info: ", "Address", addrs, "ID", hostInfo.ID)
	
	//configure our own ping protocol
	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)


	//here lets format the input arg to multiaddr 
	addr, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		logger.Errorf("Error while formating addr format: %s", err.Error())
		return err
	}
	//create a peer info from the addr
	peer, err := peerstore.AddrInfoFromP2pAddr(addr)
	if err != nil {
		logger.Errorf("Error while creating peer info: %s", err.Error())
		return err
	}

	//establish a connection from the host to peer
	ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Minute)
	defer cancel()
	if err := node.Connect(ctx, *peer); err != nil {
		logger.Errorf("Error while pinging peer %s: ", err.Error())
		cancel()
	}

	//ping to check if there's a connection
	ch := pingService.Ping(ctx, peer.ID)
	for i := 0; i < 5; i++ {
		res := <-ch
		logger.Infoln("ping response. RTT: ", res.RTT)
	}

	return nil
}

func Listen(logger *zap.SugaredLogger){
	logger.Infoln("Node listening...")
	node := Start(logger)
	defer Close(logger, node)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	logger.Infoln("Received signal shutting down...")
}


