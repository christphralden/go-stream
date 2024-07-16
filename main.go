package main

import (
	"log"
	"github.com/christpheralden/go-stream/server/controllers"
)

func main(){
  srv := server.NewServer(
    server.WithProtocol(server.TCP),
    server.WithMaxConn(10),
  )
  
  srv.ShowOptions()

  ready := make(chan error)

  go func(){
    err := srv.Start()
    ready <- err
  }()


  if err := <-ready; err != nil {
    log.Println("Error in server: ", err)
    return
  }

  
  srv.Stop()
  log.Println("closed")
}
