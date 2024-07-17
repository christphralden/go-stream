package main

import (
	"log"

	"github.com/christpheralden/go-stream/internal/server/core"
)

func main() {


	srv := core.NewServer(
		core.WithProtocol(core.TCP),
		core.WithMaxConn(10),
	)

  defer srv.Stop()

	srv.ShowOptions()

	ready := make(chan error)
	go func() {
		err := srv.Start()
		ready <- err
	}()

	if err := <-ready; err != nil {
		log.Println("Error in server: ", err)
		return
	}

	srv.ShowConnectionStatus()
}
