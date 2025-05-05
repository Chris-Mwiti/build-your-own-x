package websockets

import (
	"bufio"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
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

type Conn struct {
	config *Config
	request *http.Request

	buf *bufio.ReadWriter
	rwc *io.ReadWriteCloser

	rio sync.Mutex
	frameReaderFactory

	wio sync.Mutex
	frameWriterFactory

	frameHandler
	PayloadType byte
	defaultCloseStatus int

	//limits the size of frame payload received over conn
	MaxPayloadBytes int
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

//framereader is an interface to read a websocket frame
type frameReader interface {
	//Reader is to read payload of the frame
	io.Reader

	//PayloadType returns payload type
	PayloadType()byte

	//HeaderReader returns a reader to read header of the frame
	HeaderReader() io.Reader

	//TrailerReader returns a reader to read trailer of the frame.
	//If it return nil, there is no trailer in the frame
	TrailerReader() io.Reader


	//returns the total length of the frame
	Len() int
}

type frameReaderFactory interface {
	NewFrameReader() (r frameReader, err error)
}

type frameWriter interface {
	io.WriteCloser
}

type frameWriterFactory interface {
	NewFrameWriter(payloadType byte) (W frameWriter, err error)
}

type frameHandler interface {
	HandleFrame(frame frameReader) (r frameReader, err error)
	WriteClose(status int)(err error)
}

//DialError -> an error that occurs while dialling a websocket server
type DialError struct {
	*Config
}