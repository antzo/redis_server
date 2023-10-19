package redis_server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// Server that handle client connection using Redis protocol (RESP2)
type Server struct {
	config   ServerConfig
	executor map[string]MessageHandler
	messages chan clientMessage
	errors   chan error
	sigs     chan os.Signal
}

type ServerConfig struct {
	Port           int
	ReadBufferSize int
}

type clientMessage struct {
	message Message
	con     net.Conn
}

// Start to listen for tcp connections and handle each of them concurrently.
// When the context is cancelled it stops and close the connection
func (s *Server) Start() (err error) {
	l, e := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if e != nil {
		return e
	}

	defer func() { err = l.Close() }()

	go func() {
		for {
			select {
			case e := <-s.errors:
				log.Print(e)
			case clientMsg := <-s.messages:
				cmd, err := GetCommand(clientMsg.message)
				if err != nil {
					// TODO: Send error message
					log.Println(err)
					continue
				}

				if handler, ok := s.executor[cmd.name]; ok {
					response := handler(cmd)
					clientMsg.con.Write(response.Data())
				} else {
					clientMsg.con.Write([]byte("+OK\r\n"))
				}
			}
		}
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			s.errors <- err
		}
		log.Println("client connected")

		go s.read(conn)
	}
}

// read creates a buffer of size ServerConfig.ReadBufferSize and waits for new data coming from the client.
// When new data is received, send a clientMessage to Server.messages channel
func (s *Server) read(c net.Conn) {
	defer func() {
		if err := c.Close(); err != nil {
			s.errors <- err
		}
	}()

	for {
		reader := bufio.NewReader(c)
		buf := make([]byte, s.config.ReadBufferSize)

		// Blocking call till we receive data from client
		if _, err := reader.Read(buf); err != nil {
			// client disconnected
			if err == io.EOF {
				return
			}

			s.errors <- err
			return
		}

		message, err := Deserialize(buf)
		if err != nil {
			s.errors <- err
			return
		}

		s.messages <- clientMessage{message: message, con: c}
	}
}

func NewServer(c ServerConfig) *Server {
	// TODO: Add options to ServerConfig
	c.ReadBufferSize = 128

	return &Server{
		config: c,
		executor: map[string]MessageHandler{
			"ping": Ping,
			"echo": Echo,
			"set":  Set,
			"get":  Get,
		},
		messages: make(chan clientMessage),
		errors:   make(chan error),
		sigs:     make(chan os.Signal, 1),
	}
}
