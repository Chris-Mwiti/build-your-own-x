package nodes

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
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
	//here lets generate the node info
	hostInfo := peerstore.AddrInfo{
		ID: node.ID(),
		Addrs: node.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(&hostInfo)
	if err != nil {
		logger.Fatalf("Error while creating host addr from info: %s", err.Error())
	}

	//convert addresses to string format with newline separators
	//@todo: Add the capability of storing string json formatted syntax to capture and data correctly
	addrStrings := make([]string, len(addrs))
	for i, addr := range addrs {
		addrStrings[i] = addr.String()
	}

	addrData := []byte(strings.Join(addrStrings, "\n") + "\n")
	outputPath := "output/addresses.txt"
	f, err := os.Create(outputPath)
	if err != nil {
		logger.Fatalf("Error while creating address book: %v", err)
	}
	defer f.Close()

	//write the addresses to file
	_, err = f.Write(addrData)
	if err != nil {
		logger.Errorf("Error while writing to file: %v", err)
	}

	//ensure data is flushed to disk
	err = f.Sync()
	if err != nil {
		logger.Fatalf("Error syncing file: %v", err)
	}
	logger.Infow("Host Info: ", "Address", addrStrings, "ID", hostInfo.ID)

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



func Ping(logger *zap.SugaredLogger, peerAddrs []string)(error){
	node := Start(logger)
	defer Close(logger, node)

	//configure our own ping protocol
	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)

	for _, peerAddr := range peerAddrs {
		//here lets format the input arg to multiaddr 
		addr, err := multiaddr.NewMultiaddr(peerAddr)
		if err != nil {
			logger.Errorf("Error while formating addr format:%s : %s",addr, err.Error())
			return err
		}
		//create a peer info from the addr
		peer, err := peerstore.AddrInfoFromP2pAddr(addr)
		if err != nil {
			logger.Errorf("Error while creating peer info:%s : %s", peer.ID, err.Error())
			return err
		}

		//establish a connection from the host to peer
		ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Minute)
		defer cancel()
		if err := node.Connect(ctx, *peer); err != nil {
			logger.Errorf("Error while connecting peer %s: ", err.Error())
			cancel()
		}

		//ping to check if there's a connection
		ch := pingService.Ping(ctx, peer.ID)
		for i := 0; i < 4; i++ {
			res := <-ch
			logger.Infoln("ping response. RTT: ", res.RTT)
		}
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

//a handle function to handle read and write streams
func handleStream(logger *zap.SugaredLogger, stream network.Stream) {
	logger.Infoln("New stream created")

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	//@todo implement the read and write data to stream
	go readData(rw)
	go writeData(rw)
}

func readData(rw *bufio.ReadWriter) (error) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			return errors.New(fmt.Sprintf("[ReadData]: error while reading data: %v", err.Error()))
		}

		//here if we have an empty string we just loop out
		if str == ""{
			return nil
		}
		 
		if str != "\n"{
			fmt.Printf("\x1b[32m%s\x1b[0m>", str)
		} 
	}
}

func writeData(rw *bufio.ReadWriter) (error) {
	reader := bufio.NewReader(os.Stdin)	
	for {
		fmt.Print("> ")
		sendData, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error while reading stream data")
			return errors.New(fmt.Sprintf("[WriteData]: error while reading data: %v", err.Error()))
		}
		_, err = rw.WriteString(sendData)
		if err != nil {
			log.Println("Error while writing stream data")
			return errors.New(fmt.Sprintf("[WriteData]: error while writing data: %v", err.Error()))
		}
	}
}




