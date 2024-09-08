package main

import (
	"fmt"
	"net"
)

type Message struct {
	from    string
	payload []byte
}
type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 10),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	s.ln = ln
	go s.acceptLoop()
	fmt.Print("tcp Server has Started")

	defer s.ln.Close()

	<-s.quitch
	close(s.msgch)
	return nil
}

func (s *Server) acceptLoop() {

	for {
		conn, err := s.ln.Accept()

		if err != nil {
			fmt.Println("Accepted an error", err)
			continue
		}
		fmt.Println("Accepting connection from:", conn.RemoteAddr())
		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read Error", err)
			continue
		}

		s.msgch <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buf[:n],
		}
	}

}

func main() {
	server := NewServer(":3000")

	go func() {
		for msg := range server.msgch {
			fmt.Printf("Recived Message From Conn (%s): %s", msg.from, msg.payload)
		}
	}()

	err := server.Start()
	if err != nil {
		fmt.Println("Error Starting", err)
	}

}
