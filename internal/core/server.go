package core

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/christpheralden/go-stream/pkg/types"
)

type Server interface{
  Start()
  Stop()
}

type BaseServer struct{
  ServerOptions
  Ln          net.Listener 
  QuitCh      chan struct{}
  Wg          sync.WaitGroup
}

func NewServer(opts ...ServerOptionFunc) *BaseServer  {
  o := defaultOptions()

  for _, fn := range opts {
    fn(&o)
  }
  
  return &BaseServer{
  	ServerOptions: o,
  	Ln:            nil,
  	QuitCh:        make(chan struct{}),
  	Wg:            sync.WaitGroup{},
  }

}

func (s *BaseServer) Start() error {
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

func (s *BaseServer) Stop() {
	close(s.QuitCh)
	if s.Ln != nil {
		s.Ln.Close()
	}
	s.Wg.Wait()
	log.Println("Server stopping")
}


func (s *BaseServer) AcceptLoop() {
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


func (s *BaseServer) readLoop(conn net.Conn) {
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

    log.Printf("Server recieved from %v: %v", getConnectionStatus(s.Ln), string(payload.Bytes()))

    reply := types.Binary("Read")

    _, err = reply.WriteTo(conn)
    
    if err != nil {
      log.Println("Error writing response: ", err)
      break
    }
  }
}

func (s *BaseServer) ShowOptions() {
  if s == nil{
    log.Println("Server was not created")
    return
  }

  log.Printf("%+v\n", s)
}

func (s *BaseServer) ShowConnectionStatus(){
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
    log.Println("Server error listening: ", err)
    return
  }

  log.Printf("Listening to %v on port %v\n", host, port)
}

func getConnectionStatus(ln net.Listener) string { //TODO : remember
  if ln == nil {
    return "Server has not been started"
  }

  host, port, err := net.SplitHostPort(ln.Addr().String())
  
  if err != nil {
    return fmt.Sprintf("Server error listening : %+v", err)
  }
      
  return fmt.Sprintf("%v:%v\n", host, port)

}
