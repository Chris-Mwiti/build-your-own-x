package websockets

import (
	"bufio"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
)

//represents the configuration protocol for a websocket connection
type Config struct {
	//websocket server address
	Location *url.URL
	//websocket client origin
	Origin *url.URL
	//websocket subprotocols
	Protocol []string
	//websocket protocol version
	Version int
	//TLS config for secure websocket(wss)
	TlsConfig *tls.Config
	//additional header fiels to be sent in websocket opening handshake
	Header http.Header
	//Dialer used when opening websocket connection
	Dialer *net.Dialer

	handshakeData map[string]string

}

//serverHandshaker is an interface to handle websocket server side handshake
type serverHandshaker interface {
	// ReadHandshake reads handshake request message from client
	// Return http repsonse code and error if any
	ReadHandshake(buf *bufio.Reader, req *http.Request) (code int, err error)

	//AcceptHandshake accepts the client handshake request and sends
	//handshake response back to the client
	AcceptHandshake(buf *bufio.Writer)(err error)

	//NewServerConn creates a new Websocket connection
	NewServerConn(buf *bufio.ReadWriter, rwc io.ReadWriteCloser, request *http.Request) (conn *Conn)
}

//DialError -> an error that occurs while dialling a websocket server
type DialError struct {
	*Config
}