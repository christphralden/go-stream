package server

import (
	"io"
	"log"
	"net"
	"sync"

	"github.com/christpheralden/go-stream/types"
)


const (
  TCP = "tcp"
  UDP = "udp"
)

type ServerOptionFunc func (*ServerOptions)

type ServerOptions struct{
  Id          string
  Protocol    string 
  ListenAddr  string
  MaxConn     int8
  Tls         bool
}

type Server struct{
  ServerOptions
  Ln          net.Listener 
  QuitCh      chan struct{}
  Wg          sync.WaitGroup
}

func defaultOptions() ServerOptions {
  return ServerOptions{
    Id: "default",
    Protocol: TCP,
    ListenAddr: "localhost:3306",
    MaxConn: 100,
    Tls: false,
  }
}

func WithTls(opts *ServerOptions){
  opts.Tls = true
}

func WithMaxConn(n int8) ServerOptionFunc {
  return func(opts *ServerOptions){
    opts.MaxConn = n
  }
}

func WithProtocol(protocol string) ServerOptionFunc {
  return func(opts *ServerOptions){
    opts.Protocol = protocol
  }
}


func NewServer(opts ...ServerOptionFunc) *Server  {
  o := defaultOptions()

  for _, fn := range opts {
    fn(&o)
  }
  
  return &Server{
    ServerOptions: o,
    QuitCh: make(chan struct{}),
  }

}

func (s *Server) ShowOptions() {
  if s == nil{
    log.Println("Server was not created")
    return
  }

  log.Printf("%+v\n", s)
}

func (s *Server) ShowConnectionStatus(){
  if s == nil {
    log.Println("Server was not created")
    return
  }

  if s.Ln == nil {
    log.Println("Server has not been started")
    return
  }

  host, port, err := net.SplitHostPort(s.Ln.Addr().String())

  if err != nil {
    log.Println("Server error listening")
    return
  }

  log.Printf("Listening to %v on port %v\n", host, port)
}



func (s *Server) Start() error {
  ln, err := net.Listen(s.Protocol, s.ListenAddr)

  if err != nil {
    return err
  }

  defer ln.Close()
  s.Ln = ln

  s.ShowConnectionStatus()
  
  go s.AcceptLoop()

  <-s.QuitCh
  
  return nil
}

func (s *Server) Stop() {
	close(s.QuitCh)
	if s.Ln != nil {
		s.Ln.Close()
	}
	s.Wg.Wait()
	log.Println("Server shutdown")
}


func (s *Server) AcceptLoop() {
  for{
    conn, err := s.Ln.Accept()

    if err != nil {
      select {
      case <-s.QuitCh:
        return
      default:
        log.Println("Error accepting connection: ", err)
      } 
      continue
    }

    s.Wg.Add(1)
    go s.readLoop(conn)
  }
}


func (s *Server) readLoop(conn net.Conn) {
  defer conn.Close()
  defer s.Wg.Done()

  for{
    payload, err := types.Decode(conn)

    if err != nil {
      if err == io.EOF{
        log.Println("Client has disconnected")
      }else{
        log.Println("Something went wrong: ", err)
      }

      break
    }

    log.Println("Server recieved: ", string(payload.Bytes()))

    reply := types.Binary("Read")

    _, err = reply.WriteTo(conn)
    
    if err != nil {
      log.Println("Error writing response: ", err)
      break
    }
  }
}
