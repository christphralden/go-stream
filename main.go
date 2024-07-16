package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/christpheralden/go-stream/client"
  "github.com/christpheralden/go-stream/server/controllers"
)

func simulateServer(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	srv := server.NewServer(
		server.WithProtocol(server.TCP),
		server.WithMaxConn(10),
	)
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

	time.Sleep(10 * time.Second)

	srv.Stop()
	log.Println("Server closed")
}

func simulateClient(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	cl := client.NewClient(
		client.WithProtocol(client.TCP),
	)

	if err := cl.Dial(); err != nil {
		log.Println("Error connecting to server:", err)
		return
	}

	fmt.Println("Input your messages:")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()

		if message == "quit" {
			fmt.Println("Disconnecting from server")
			cl.Stop()
			break
		}

		if err := cl.SendMessage(message); err != nil {
			fmt.Println("Error sending message:", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from input:", err)
	}
}

func main() {
	var wg sync.WaitGroup

	go simulateServer(&wg)

	time.Sleep(1 * time.Second)

	go simulateClient(&wg)

	wg.Wait()
}
