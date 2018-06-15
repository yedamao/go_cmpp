package cmpptest

import (
	"log"
	"net"
	"time"
)

// Server is a mock SMGW server
type Server struct {
	listener *net.TCPListener

	done chan struct{}
}

func NewServer(addr string) (*Server, error) {

	laddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return nil, err
	}
	log.Println("Server Listen: ", listener.Addr())

	return &Server{
		listener: listener,
		done:     make(chan struct{}),
	}, nil
}

func (s *Server) Run() {

	for {
		select {
		case <-s.done:
			log.Println("stopping listening on", s.listener.Addr())
			return
		default:
		}

		s.listener.SetDeadline(time.Now().Add(1e9))
		conn, err := s.listener.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Println(err)
		}
		log.Println("conn from: ", conn.RemoteAddr())

		newSession(conn)
	}
}

func (s *Server) Stop() {
	close(s.done)
}
